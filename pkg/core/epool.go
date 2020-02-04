package core

import (
	"github.com/gobwas/ws/wsutil"
	"github.com/golang/glog"
	"github.com/mkawserm/pisig/pkg/context"
	"github.com/mkawserm/pisig/pkg/event"
	"golang.org/x/sys/unix"
	"net"
	"reflect"
	"sync"
	"syscall"
)

//type EPoolProcessMessageHook func(conn net.Conn, msg []byte, opCode byte) error
type EPoolRemoveConnectionHook func(conn net.Conn) error

type EPool struct {
	mPisigContext   *context.PisigContext
	mFd             int
	mConnectionMap  map[int]net.Conn
	mEventQueueSize int
	mWaitingTime    int

	mRemoveConnectionHook EPoolRemoveConnectionHook

	mRWLock *sync.RWMutex
}

func NewEPool(ctx *context.PisigContext,
	removeConnectionHook EPoolRemoveConnectionHook,
) (*EPool, error) {
	fd, err := unix.EpollCreate1(0)

	if err != nil {
		return nil, err
	}

	return &EPool{
		mPisigContext:         ctx,
		mFd:                   fd,
		mRWLock:               &sync.RWMutex{},
		mConnectionMap:        make(map[int]net.Conn),
		mEventQueueSize:       ctx.PisigSettings.EventPoolQueueSize,
		mWaitingTime:          ctx.PisigSettings.EventPoolWaitingTime,
		mRemoveConnectionHook: removeConnectionHook,
	}, nil
}

//func NewDefaultEventPool() (*EPool, error) {
//	return NewEPool(100,100)
//}

func (e *EPool) Setup() {
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

func (e *EPool) GetConnection(connectionId int) (net.Conn, bool) {
	e.mRWLock.RLock()
	defer e.mRWLock.RUnlock()

	v, ok := e.mConnectionMap[connectionId]
	return v, ok
}

func (e *EPool) AddConnection(conn net.Conn) error {
	// Extract file descriptor associated with the connection
	if glog.V(3) {
		glog.Infof("Connection: %d\n", conn)
	}

	fd := WebsocketFileDescriptor(conn)

	if glog.V(3) {
		glog.Infof("Connect id: %d\n", fd)
	}

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

func (e *EPool) RemoveConnection(conn net.Conn) error {
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

func (e *EPool) Wait() ([]net.Conn, error) {
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

func (e *EPool) TotalActiveConnections() int {
	e.mRWLock.RLock()
	defer e.mRWLock.RUnlock()

	return len(e.mConnectionMap)
}

func (e *EPool) GetConnectionIdSlice() []int {
	e.mRWLock.RLock()
	defer e.mRWLock.RUnlock()
	var idList []int

	for i := range e.mConnectionMap {
		idList = append(idList, i)
	}

	return idList
}

func (e *EPool) GetConnectionSlice() []net.Conn {
	e.mRWLock.RLock()
	defer e.mRWLock.RUnlock()
	var connections []net.Conn

	for _, con := range e.mConnectionMap {
		connections = append(connections, con)
	}

	return connections
}

func (e *EPool) GetConnectionMap() map[int]net.Conn {
	e.mRWLock.RLock()
	defer e.mRWLock.RUnlock()

	return e.mConnectionMap
}

func (e *EPool) RunMainEventLoop() {
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

			//fmt.Println(msg)
			//fmt.Println(opCode)
			//fmt.Println(err)
			//
			//fmt.Println(e.mRemoveConnectionHook)
			//fmt.Println(e.mProcessMessageHook)

			if err != nil {
				if e.mRemoveConnectionHook != nil {
					if err := e.mRemoveConnectionHook(conn); err != nil {
						glog.Errorf("Failed to remove connection, error: %v\n", e)
					}
				}
				_ = conn.Close()
				continue
			}

			if err := e.ProcessWebSocketMessage(conn, msg, byte(opCode)); err != nil {
				glog.Errorf("Failed to process message, error: %v\n", e)
			}
		}
	}
}

func (e *EPool) ProcessWebSocketMessage(conn net.Conn, msg []byte, opCode byte) error {
	if glog.V(3) {
		glog.Infof("Process message\n")
	}

	topic := event.Topic{
		Name: "WebSocketEvent",
		Key:  []byte(""),
		Data: event.WebSocketEvent{
			Conn:    conn,
			OpCode:  opCode,
			Message: msg,
		},
	}

	e.mPisigContext.ProduceTopic(topic)

	//if glog.V(3) {
	//	glog.Infof("Message: %s\n", string(msg))
	//	glog.Infof("OpCode: %d\n", opCode)
	//	//err := wsutil.WriteServerText(conn,[]byte("Hello World"))
	//	//println(err)
	//}

	return nil
}

func WebsocketFileDescriptor(conn net.Conn) int {
	if glog.V(3) {
		glog.Infof("Inspecting websocket file descriptor")
	}

	//if glog.V(3) {
	//	glog.Infof("Connection: %d\n", conn)
	//}

	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	if glog.V(3) {
		//glog.Infof(pfdVal.FieldByName("Sysfd").String())
		glog.Infof("Websocket file descriptor found\n")
	}
	return int(pfdVal.FieldByName("Sysfd").Int())
}
