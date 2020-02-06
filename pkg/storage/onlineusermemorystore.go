package storage

import "sync"

type OnlineUserMemoryStore struct {
	mSocketIdToUniqueId     map[int]string
	mUniqueIdToSocketIdList map[string][]int

	mUniqueIdToData map[string]interface{}

	mUniqueIdToGroupId     map[string]string
	mGroupIdToUniqueIdList map[string][]string

	mRWMutex *sync.RWMutex
}

func (o *OnlineUserMemoryStore) AddUser(uniqueId string, groupId string, socketId int, data interface{}) bool {

	return false
}

func (o *OnlineUserMemoryStore) RemoveUser(socketId int) bool {

	return false
}

func (o *OnlineUserMemoryStore) GetUniqueIdFromSocketId(socketId int) string {

	return ""
}

func (o *OnlineUserMemoryStore) GetSocketIdListFromUniqueId(uniqueId string) []int {

	return []int{}
}

func (o *OnlineUserMemoryStore) GetDataFromUniqueId(uniqueId string) interface{} {

	return nil
}

func (o *OnlineUserMemoryStore) GetGroupIdFromUniqueId(uniqueId string) string {

	return ""
}

func (o *OnlineUserMemoryStore) GetUniqueIdListFromGroupId(groupId string) []string {

	return []string{}
}

func (o *OnlineUserMemoryStore) GetGroupIdList() []string {

	return []string{}
}

func (o *OnlineUserMemoryStore) GetUniqueIdList() []string {

	return []string{}
}

func (o *OnlineUserMemoryStore) GetTotalGroupId() int {

	return 0
}

func (o *OnlineUserMemoryStore) GetTotalUniqueId() int {

	return 0
}
