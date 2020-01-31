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

type EventPool struct {
	fd             int
	connectionMap  map[int]net.Conn
	eventQueueSize int
	waitingTime    int

	processMessageHook   func(conn net.Conn, msg []byte, opCode byte) error
	removeConnectionHook func(conn net.Conn) error

	lock *sync.RWMutex
}

func MakeEventPool() (*EventPool, error) {
	fd, err := unix.EpollCreate1(0)

	if err != nil {
		return nil, err
	}

	return &EventPool{
		fd:             fd,
		lock:           &sync.RWMutex{},
		connectionMap:  make(map[int]net.Conn),
		eventQueueSize: 100,
		waitingTime:    100,
	}, nil

}

func MakeCustomEventPool(eventQueueSize int, waitingTime int) (*EventPool, error) {
	fd, err := unix.EpollCreate1(0)

	if err != nil {
		return nil, err
	}

	return &EventPool{
		fd:             fd,
		lock:           &sync.RWMutex{},
		connectionMap:  make(map[int]net.Conn),
		eventQueueSize: eventQueueSize,
		waitingTime:    waitingTime,
	}, nil

}

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
	e.lock.RLock()
	defer e.lock.RUnlock()

	v, ok := e.connectionMap[connectionId]
	return v, ok
}

func (e *EventPool) AddConnection(conn net.Conn) error {
	// Extract file descriptor associated with the connection
	fd := WebsocketFileDescriptor(conn)

	err := unix.EpollCtl(e.fd,
		syscall.EPOLL_CTL_ADD,
		fd,
		&unix.EpollEvent{
			Events: unix.POLLIN | unix.POLLHUP,
			Fd:     int32(fd),
		})

	if err != nil {
		return err
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	e.connectionMap[fd] = conn
	if len(e.connectionMap)%100 == 0 {
		glog.V(1).Infof("Total number of connections: %v\n", len(e.connectionMap))
	}

	return nil
}

func (e *EventPool) RemoveConnection(conn net.Conn) error {
	fd := WebsocketFileDescriptor(conn)

	err := unix.EpollCtl(e.fd,
		syscall.EPOLL_CTL_DEL,
		fd, nil)

	if err != nil {
		return err
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	delete(e.connectionMap, fd)

	if len(e.connectionMap)%100 == 0 {
		glog.V(1).Infof("Total number of connections: %v\n", len(e.connectionMap))
	}

	return nil
}

func (e *EventPool) Wait() ([]net.Conn, error) {
	events := make([]unix.EpollEvent, e.eventQueueSize)
	n, err := unix.EpollWait(e.fd, events, e.waitingTime)

	if err != nil {
		return nil, err
	}

	e.lock.RLock()
	defer e.lock.RUnlock()

	var connections []net.Conn
	for i := 0; i < n; i++ {
		conn := e.connectionMap[int(events[i].Fd)]
		connections = append(connections, conn)
	}
	return connections, nil
}

func (e *EventPool) TotalActiveConnections() int {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return len(e.connectionMap)
}

func (e *EventPool) GetConnectionIdSlice() []int {
	e.lock.RLock()
	defer e.lock.RUnlock()
	var idList []int

	for i := range e.connectionMap {
		idList = append(idList, i)
	}

	return idList
}

func (e *EventPool) GetConnectionSlice() []net.Conn {
	e.lock.RLock()
	defer e.lock.RUnlock()
	var connections []net.Conn

	for _, con := range e.connectionMap {
		connections = append(connections, con)
	}

	return connections
}

func (e *EventPool) GetConnectionMap() map[int]net.Conn {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.connectionMap
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
				if e.removeConnectionHook != nil {
					if err := e.removeConnectionHook(conn); err != nil {
						glog.Errorf("Failed to remove connection, error: %v\n", e)
					}
				}
				_ = conn.Close()
				continue
			}

			if e.processMessageHook != nil {
				if err := e.processMessageHook(conn, msg, byte(opCode)); err != nil {
					glog.Errorf("Failed to process message, error: %v\n", e)
				}
			}
		}
	}
}

func WebsocketFileDescriptor(conn net.Conn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}
