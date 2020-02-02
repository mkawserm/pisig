package core

import (
	"github.com/gobwas/ws/wsutil"
	"github.com/golang/glog"
	"golang.org/x/sys/unix"
	"net"
	"reflect"
	"sync"
	"syscall"
)

type EventPoolProcessMessageHook func(conn net.Conn, msg []byte, opCode byte) error
type EventPoolRemoveConnectionHook func(conn net.Conn) error

type EventPool struct {
	mFd             int
	mConnectionMap  map[int]net.Conn
	mEventQueueSize int
	mWaitingTime    int

	mProcessMessageHook   EventPoolProcessMessageHook
	mRemoveConnectionHook EventPoolRemoveConnectionHook

	mRWLock *sync.RWMutex
}

func NewEventPool(
	eventQueueSize int,
	waitingTime int,
	processMessageHook EventPoolProcessMessageHook,
	removeConnectionHook EventPoolRemoveConnectionHook,
) (*EventPool, error) {
	fd, err := unix.EpollCreate1(0)

	if err != nil {
		return nil, err
	}

	return &EventPool{
		mFd:                   fd,
		mRWLock:               &sync.RWMutex{},
		mConnectionMap:        make(map[int]net.Conn),
		mEventQueueSize:       eventQueueSize,
		mWaitingTime:          waitingTime,
		mProcessMessageHook:   processMessageHook,
		mRemoveConnectionHook: removeConnectionHook,
	}, nil
}

//func NewDefaultEventPool() (*EventPool, error) {
//	return NewEventPool(100,100)
//}

func (e *EventPool) Setup() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		glog.Errorf("Get rlimit call failed \n")
		panic(err)
	}

	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		glog.Errorf("Set rlimit call failed \n")
		panic(err)
	}
}

func (e *EventPool) GetConnection(connectionId int) (net.Conn, bool) {
	e.mRWLock.RLock()
	defer e.mRWLock.RUnlock()

	v, ok := e.mConnectionMap[connectionId]
	return v, ok
}

func (e *EventPool) AddConnection(conn net.Conn) error {
	// Extract file descriptor associated with the connection
	fd := WebsocketFileDescriptor(conn)

	err := unix.EpollCtl(e.mFd,
		syscall.EPOLL_CTL_ADD,
		fd,
		&unix.EpollEvent{
			Events: unix.POLLIN | unix.POLLHUP,
			Fd:     int32(fd),
		})

	if err != nil {
		return err
	}

	e.mRWLock.Lock()
	defer e.mRWLock.Unlock()

	e.mConnectionMap[fd] = conn
	if len(e.mConnectionMap)%100 == 0 {
		if glog.V(3) {
			glog.Infof("Total number of connections: %v\n", len(e.mConnectionMap))
		}
	}

	return nil
}

func (e *EventPool) RemoveConnection(conn net.Conn) error {
	fd := WebsocketFileDescriptor(conn)

	err := unix.EpollCtl(e.mFd,
		syscall.EPOLL_CTL_DEL,
		fd, nil)

	if err != nil {
		return err
	}

	e.mRWLock.Lock()
	defer e.mRWLock.Unlock()

	delete(e.mConnectionMap, fd)

	if len(e.mConnectionMap)%100 == 0 {
		if glog.V(3) {
			glog.Infof("Total number of connections: %v\n", len(e.mConnectionMap))
		}
	}

	return nil
}

func (e *EventPool) Wait() ([]net.Conn, error) {
	events := make([]unix.EpollEvent, e.mEventQueueSize)
	n, err := unix.EpollWait(e.mFd, events, e.mWaitingTime)

	if err != nil {
		return nil, err
	}

	e.mRWLock.RLock()
	defer e.mRWLock.RUnlock()

	var connections []net.Conn
	for i := 0; i < n; i++ {
		conn := e.mConnectionMap[int(events[i].Fd)]
		connections = append(connections, conn)
	}
	return connections, nil
}

func (e *EventPool) TotalActiveConnections() int {
	e.mRWLock.RLock()
	defer e.mRWLock.RUnlock()

	return len(e.mConnectionMap)
}

func (e *EventPool) GetConnectionIdSlice() []int {
	e.mRWLock.RLock()
	defer e.mRWLock.RUnlock()
	var idList []int

	for i := range e.mConnectionMap {
		idList = append(idList, i)
	}

	return idList
}

func (e *EventPool) GetConnectionSlice() []net.Conn {
	e.mRWLock.RLock()
	defer e.mRWLock.RUnlock()
	var connections []net.Conn

	for _, con := range e.mConnectionMap {
		connections = append(connections, con)
	}

	return connections
}

func (e *EventPool) GetConnectionMap() map[int]net.Conn {
	e.mRWLock.RLock()
	defer e.mRWLock.RUnlock()

	return e.mConnectionMap
}

func (e *EventPool) RunMainEventLoop() {
	for {
		connections, err := e.Wait()
		if err != nil {
			glog.Warningf("Failed to wait on eventPool, error: %v\n", err)
			continue
		}

		for _, conn := range connections {

			if conn == nil {
				break
			}

			msg, opCode, err := wsutil.ReadClientData(conn)

			if err != nil {
				if e.mRemoveConnectionHook != nil {
					if err := e.mRemoveConnectionHook(conn); err != nil {
						glog.Errorf("Failed to remove connection, error: %v\n", e)
					}
				}
				_ = conn.Close()
				continue
			}

			if e.mProcessMessageHook != nil {
				if err := e.mProcessMessageHook(conn, msg, byte(opCode)); err != nil {
					glog.Errorf("Failed to process message, error: %v\n", e)
				}
			}
		}
	}
}

func WebsocketFileDescriptor(conn net.Conn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("mFd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}
