package core

import (
	"github.com/gobwas/ws/wsutil"
	"github.com/golang/glog"
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
	mPisigContext   *PisigContext
	mFd             int
	mConnectionMap  map[int]net.Conn
	mEventQueueSize int
	mWaitingTime    int

	mRemoveConnectionHook EPoolRemoveConnectionHook

	mRWLock *sync.RWMutex
}

func NewEPool(pisigContext *PisigContext,
	removeConnectionHook EPoolRemoveConnectionHook,
) (*EPool, error) {
	fd, err := unix.EpollCreate1(0)

	if err != nil {
		return nil, err
	}

	return &EPool{
		mPisigContext:         pisigContext,
		mFd:                   fd,
		mRWLock:               &sync.RWMutex{},
		mConnectionMap:        make(map[int]net.Conn),
		mEventQueueSize:       pisigContext.PisigSettings.EventPoolQueueSize,
		mWaitingTime:          pisigContext.PisigSettings.EventPoolWaitingTime,
		mRemoveConnectionHook: removeConnectionHook,
	}, nil
}

//func NewDefaultEventPool() (*EPool, error) {
//	return NewEPool(100,100)
//}

func (ePool *EPool) Setup() {
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

func (ePool *EPool) GetConnection(connectionId int) (net.Conn, bool) {
	ePool.mRWLock.RLock()
	defer ePool.mRWLock.RUnlock()

	v, ok := ePool.mConnectionMap[connectionId]
	return v, ok
}

func (ePool *EPool) AddConnection(conn net.Conn) error {
	// Extract file descriptor associated with the connection
	if glog.V(3) {
		glog.Infof("Connection: %d\n", conn)
	}

	fd := WebsocketFileDescriptor(conn)

	if glog.V(3) {
		glog.Infof("Connect id: %d\n", fd)
	}

	err := unix.EpollCtl(ePool.mFd,
		syscall.EPOLL_CTL_ADD,
		fd,
		&unix.EpollEvent{
			Events: unix.POLLIN | unix.POLLHUP,
			Fd:     int32(fd),
		})

	if err != nil {
		return err
	}

	ePool.mRWLock.Lock()
	defer ePool.mRWLock.Unlock()

	ePool.mConnectionMap[fd] = conn
	if len(ePool.mConnectionMap)%100 == 0 {
		if glog.V(3) {
			glog.Infof("Total number of connections: %v\n", len(ePool.mConnectionMap))
		}
	}

	return nil
}

func (ePool *EPool) RemoveConnection(conn net.Conn) error {
	fd := WebsocketFileDescriptor(conn)

	err := unix.EpollCtl(ePool.mFd,
		syscall.EPOLL_CTL_DEL,
		fd, nil)

	if err != nil {
		return err
	}

	ePool.mRWLock.Lock()
	defer ePool.mRWLock.Unlock()

	delete(ePool.mConnectionMap, fd)

	if len(ePool.mConnectionMap)%100 == 0 {
		if glog.V(3) {
			glog.Infof("Total number of connections: %v\n", len(ePool.mConnectionMap))
		}
	}

	return nil
}

func (ePool *EPool) Wait() ([]net.Conn, error) {
	events := make([]unix.EpollEvent, ePool.mEventQueueSize)
	n, err := unix.EpollWait(ePool.mFd, events, ePool.mWaitingTime)

	if err != nil {
		return nil, err
	}

	ePool.mRWLock.RLock()
	defer ePool.mRWLock.RUnlock()

	var connections []net.Conn
	for i := 0; i < n; i++ {
		conn := ePool.mConnectionMap[int(events[i].Fd)]
		connections = append(connections, conn)
	}
	return connections, nil
}

func (ePool *EPool) TotalActiveConnections() int {
	ePool.mRWLock.RLock()
	defer ePool.mRWLock.RUnlock()

	return len(ePool.mConnectionMap)
}

func (ePool *EPool) GetConnectionIdSlice() []int {
	ePool.mRWLock.RLock()
	defer ePool.mRWLock.RUnlock()
	var idList []int

	for i := range ePool.mConnectionMap {
		idList = append(idList, i)
	}

	return idList
}

func (ePool *EPool) GetConnectionSlice() []net.Conn {
	ePool.mRWLock.RLock()
	defer ePool.mRWLock.RUnlock()
	var connections []net.Conn

	for _, con := range ePool.mConnectionMap {
		connections = append(connections, con)
	}

	return connections
}

func (ePool *EPool) GetConnectionMap() map[int]net.Conn {
	ePool.mRWLock.RLock()
	defer ePool.mRWLock.RUnlock()

	return ePool.mConnectionMap
}

func (ePool *EPool) RunMainEventLoop() {
	for {
		connections, err := ePool.Wait()
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
			//fmt.Println(ePool.mRemoveConnectionHook)
			//fmt.Println(ePool.mProcessMessageHook)

			if err != nil {
				if ePool.mRemoveConnectionHook != nil {
					if err := ePool.mRemoveConnectionHook(conn); err != nil {
						glog.Errorf("Failed to remove connection, error: %v\n", ePool)
					}
				}
				_ = conn.Close()
				continue
			}

			if err := ePool.ProcessWebSocketMessage(conn, msg, byte(opCode)); err != nil {
				glog.Errorf("Failed to process message, error: %v\n", ePool)
			}
		}
	}
}

func (ePool *EPool) ProcessWebSocketMessage(conn net.Conn, msg []byte, opCode byte) error {
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

	ePool.mPisigContext.ProduceTopic(topic)

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
