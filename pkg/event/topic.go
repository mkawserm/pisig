package event

import (
	"reflect"
	"unsafe"
)

//type Protocol int
//
//const (
//	Unknown   Protocol = 0
//	TCP       Protocol = 1
//	UDP       Protocol = 2
//	gRPC      Protocol = 3
//	WebSocket Protocol = 4
//	HTTP      Protocol = 5
//)
//
//func (protocol Protocol) String() string {
//	names := [...]string{
//		"Unknown",
//		"TCP",
//		"UDP",
//		"gRPC",
//		"WebSocket",
//		"HTTP",
//	}
//
//	if protocol < TCP || protocol > HTTP {
//		return "Unknown"
//	}
//	return names[protocol]
//}

type Topic struct {
	Name string
	Key  string
	Data interface{}
}

func (t *Topic) DataType() string {
	return reflect.TypeOf(t.Data).String()
}

func (t *Topic) GetNameString() string {
	return t.Name
}

func (t *Topic) GetNameBytes() []byte {
	return UnsafeBytes(t.Name)
}

func (t *Topic) GetKeyString() string {
	return t.Key
}

func (t *Topic) GetKeyBytes() []byte {
	return UnsafeBytes(t.Key)
}

func (t *Topic) GetDataBytes() []byte {
	return t.Data.([]byte)
}

type TopicQueue chan Topic
type TopicQueuePool chan TopicQueue

func UnsafeString(bytes []byte) string {
	hdr := *(*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
	}))
}

func UnsafeBytes(str string) []byte {
	hdr := *(*reflect.StringHeader)(unsafe.Pointer(&str))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
		Cap:  hdr.Len,
	}))
}
