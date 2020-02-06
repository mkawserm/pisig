package storage

import (
	"sort"
	"sync"
)

type OnlineUserMemoryStore struct {
	mSocketIdToUniqueId     map[int]string
	mUniqueIdToSocketIdList map[string][]int

	mSocketIdToGroupId     map[int]string
	mGroupIdToSocketIdList map[string][]int

	mSocketIdToData map[int]interface{}
	mRWMutex        *sync.RWMutex
}

func (o *OnlineUserMemoryStore) AddUser(uniqueId string, groupId string, socketId int, data interface{}) bool {
	o.mRWMutex.Lock()
	defer o.mRWMutex.Unlock()

	o.mSocketIdToUniqueId[socketId] = uniqueId

	if _, ok := o.mUniqueIdToSocketIdList[uniqueId]; ok {
		o.mUniqueIdToSocketIdList[uniqueId] = append(o.mUniqueIdToSocketIdList[uniqueId], socketId)
	} else {
		o.mUniqueIdToSocketIdList[uniqueId] = []int{socketId}
	}

	if data != nil {
		o.mSocketIdToData[socketId] = data
	}

	if groupId != "" {
		o.mSocketIdToGroupId[socketId] = groupId
		if _, ok := o.mGroupIdToSocketIdList[groupId]; ok {
			o.mGroupIdToSocketIdList[groupId] = append(o.mGroupIdToSocketIdList[groupId], socketId)
		} else {
			o.mGroupIdToSocketIdList[groupId] = []int{socketId}
		}
	}

	return true
}

func (o *OnlineUserMemoryStore) RemoveUser(socketId int) bool {
	o.mRWMutex.Lock()
	defer o.mRWMutex.Unlock()

	var uniqueId string
	var groupId string
	var ok bool
	var uniqueIdToSocketIdList []int
	var groupIdToSocketIdList []int

	// cleanup data
	if _, ok := o.mSocketIdToData[socketId]; ok {
		delete(o.mSocketIdToData, socketId)
	}

	uniqueId, ok = o.mSocketIdToUniqueId[socketId]

	// cleanup unique id information
	if ok {
		delete(o.mSocketIdToUniqueId, socketId)
	} else {
		return false
	}

	// Remove current Socket id from Unique Id to Socket Id list map
	uniqueIdToSocketIdList, ok = o.mUniqueIdToSocketIdList[uniqueId]
	if ok {
		index := sort.Search(len(uniqueIdToSocketIdList), func(i int) bool {
			return uniqueIdToSocketIdList[i] == socketId
		})

		if index >= len(uniqueIdToSocketIdList) || uniqueIdToSocketIdList[index] != socketId {
			// invalid index nothing to do
		} else {
			uniqueIdToSocketIdList[index] = uniqueIdToSocketIdList[len(uniqueIdToSocketIdList)-1]
			uniqueIdToSocketIdList = uniqueIdToSocketIdList[0 : len(uniqueIdToSocketIdList)-1]
		}

		if len(uniqueIdToSocketIdList) == 0 {
			delete(o.mUniqueIdToSocketIdList, uniqueId)
		} else {
			o.mUniqueIdToSocketIdList[uniqueId] = uniqueIdToSocketIdList
		}
	}

	// cleanup group information
	groupId, ok = o.mSocketIdToGroupId[socketId]
	if ok {
		delete(o.mSocketIdToGroupId, socketId)
	}

	// Remove current Socket id from group to Socket Id list map
	groupIdToSocketIdList, ok = o.mGroupIdToSocketIdList[groupId]
	if ok {
		index := sort.Search(len(groupIdToSocketIdList), func(i int) bool {
			return groupIdToSocketIdList[i] == socketId
		})

		if index >= len(groupIdToSocketIdList) || groupIdToSocketIdList[index] != socketId {
			// invalid index nothing to do
		} else {
			groupIdToSocketIdList[index] = groupIdToSocketIdList[len(groupIdToSocketIdList)-1]
			groupIdToSocketIdList = groupIdToSocketIdList[0 : len(groupIdToSocketIdList)-1]
		}

		if len(groupIdToSocketIdList) == 0 {
			delete(o.mGroupIdToSocketIdList, groupId)
		} else {
			o.mGroupIdToSocketIdList[groupId] = groupIdToSocketIdList
		}
	}

	return true
}

func (o *OnlineUserMemoryStore) GetUniqueIdFromSocketId(socketId int) string {
	o.mRWMutex.Lock()
	defer o.mRWMutex.Unlock()

	val, ok := o.mSocketIdToUniqueId[socketId]
	if ok {
		return val
	}

	return ""
}

func (o *OnlineUserMemoryStore) GetSocketIdListFromUniqueId(uniqueId string) []int {
	o.mRWMutex.Lock()
	defer o.mRWMutex.Unlock()

	val, ok := o.mUniqueIdToSocketIdList[uniqueId]
	if ok {
		return val
	}

	return []int{}
}

func (o *OnlineUserMemoryStore) GetDataFromSocketId(socketId int) interface{} {
	o.mRWMutex.Lock()
	defer o.mRWMutex.Unlock()

	val, ok := o.mSocketIdToData[socketId]
	if ok {
		return val
	}

	return nil
}

func (o *OnlineUserMemoryStore) GetGroupIdFromSocketId(socketId int) string {
	o.mRWMutex.Lock()
	defer o.mRWMutex.Unlock()

	val, ok := o.mSocketIdToGroupId[socketId]
	if ok {
		return val
	}

	return ""
}

func (o *OnlineUserMemoryStore) GetSocketIdListFromGroupId(groupId string) []int {
	o.mRWMutex.Lock()
	defer o.mRWMutex.Unlock()

	val, ok := o.mGroupIdToSocketIdList[groupId]
	if ok {
		return val
	}

	return []int{}
}

func (o *OnlineUserMemoryStore) GetGroupIdList() []string {
	o.mRWMutex.Lock()
	defer o.mRWMutex.Unlock()

	return []string{}
}

func (o *OnlineUserMemoryStore) GetUniqueIdList() []string {
	o.mRWMutex.Lock()
	defer o.mRWMutex.Unlock()

	return []string{}
}

func (o *OnlineUserMemoryStore) GetTotalGroupId() int {
	o.mRWMutex.Lock()
	defer o.mRWMutex.Unlock()

	return len(o.mGroupIdToSocketIdList)
}

func (o *OnlineUserMemoryStore) GetTotalUniqueId() int {
	o.mRWMutex.Lock()
	defer o.mRWMutex.Unlock()

	return len(o.mUniqueIdToSocketIdList)
}

func NewOnlineUserMemoryStore() *OnlineUserMemoryStore {
	return &OnlineUserMemoryStore{
		mSocketIdToUniqueId:     make(map[int]string),
		mUniqueIdToSocketIdList: make(map[string][]int),
		mSocketIdToGroupId:      make(map[int]string),
		mGroupIdToSocketIdList:  make(map[string][]int),
		mSocketIdToData:         make(map[int]interface{}),
		mRWMutex:                &sync.RWMutex{},
	}
}
