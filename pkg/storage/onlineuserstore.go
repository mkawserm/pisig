package storage

type OnlineUserStore interface {
	AddUser(uniqueId string, groupId string, socketId int, data interface{}) bool
	RemoveUser(socketId int) bool

	GetUniqueIdFromSocketId(socketId int) string
	GetSocketIdListFromUniqueId(uniqueId string) []int

	GetDataFromUniqueId(uniqueId string) interface{}

	GetGroupIdFromUniqueId(uniqueId string) string
	GetUniqueIdListFromGroupId(groupId string) []string

	GetGroupIdList() []string
	GetUniqueIdList() []string

	GetTotalGroupId() int
	GetTotalUniqueId() int
}
