package storage

type OnlineUserStore interface {
	AddUser(uniqueId string, groupId string, socketId int, data interface{}) bool
	RemoveUser(socketId int) bool

	GetUniqueIdFromSocketId(socketId int) string
	GetSocketIdListFromUniqueId(uniqueId string) []int

	GetDataFromSocketId(socketId int) interface{}

	GetGroupIdFromSocketId(socketId int) string
	GetSocketIdListFromGroupId(groupId string) []int

	GetGroupIdList() []string
	GetUniqueIdList() []string

	GetTotalGroupId() int
	GetTotalUniqueId() int
}
