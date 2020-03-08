package core

import (
	"errors"
	"sort"
	"sync"
)

type InterfaceMap map[string]interface{}

type PisigServiceRegistry struct {
	mStore         map[string]InterfaceMap
	mTopicListener map[string][]interface{}

	mGroupList   []string
	mServiceList []string

	mDefaultServiceMap    map[string]string
	mAllAvailableServices map[string][]string

	mRWLock *sync.RWMutex
}

func (psr *PisigServiceRegistry) SetDefaultServiceName(groupName, serviceName string) {
	psr.mRWLock.Lock()
	defer psr.mRWLock.Unlock()

	psr.mDefaultServiceMap[groupName] = serviceName
}

func (psr *PisigServiceRegistry) GetDefaultServiceName(groupName string) (string, bool) {
	psr.mRWLock.RLock()
	defer psr.mRWLock.RUnlock()

	val, ok := psr.mDefaultServiceMap[groupName]

	return val, ok
}

func (psr *PisigServiceRegistry) GetDefaultServiceMap() map[string]string {
	psr.mRWLock.RLock()
	defer psr.mRWLock.RUnlock()

	return psr.mDefaultServiceMap
}

func (psr *PisigServiceRegistry) GetAllAvailableServices() map[string][]string {
	psr.mRWLock.RLock()
	defer psr.mRWLock.RUnlock()

	return psr.mAllAvailableServices
}

func (psr *PisigServiceRegistry) GetTopicListenerList(topicName string) []interface{} {
	psr.mRWLock.RLock()
	defer psr.mRWLock.RUnlock()

	val, ok := psr.mTopicListener[topicName]

	if ok {
		return val
	}

	return make([]interface{}, 0)
}

func (psr *PisigServiceRegistry) AddTopicListener(topicName string, pisigService PisigService) {
	psr.mRWLock.Lock()
	defer psr.mRWLock.Unlock()

	val, ok := psr.mTopicListener[topicName]

	if !ok {
		val = make([]interface{}, 0)
		val = append(val, pisigService)
		psr.mTopicListener[topicName] = val
		return
	}

	found := false
	for i := range val {
		if val[i] == pisigService {
			found = true
		}
	}

	if !found {
		val = append(val, pisigService)
		psr.mTopicListener[topicName] = val
	}
}

func (psr *PisigServiceRegistry) GetGroupList() []string {
	psr.mRWLock.RLock()
	defer psr.mRWLock.RUnlock()

	return psr.mGroupList
}

func (psr *PisigServiceRegistry) GetServiceList() []string {
	psr.mRWLock.RLock()
	defer psr.mRWLock.RUnlock()

	return psr.mServiceList
}

func (psr *PisigServiceRegistry) IsGroupExistsInList(groupName string) bool {
	psr.mRWLock.RLock()
	defer psr.mRWLock.RUnlock()

	i := sort.Search(len(psr.mGroupList), func(i int) bool { return psr.mGroupList[i] >= groupName })
	if i < len(psr.mGroupList) && psr.mGroupList[i] == groupName {
		return true
	}

	return false
}

func (psr *PisigServiceRegistry) IsServiceExistsInList(serviceName string) bool {
	psr.mRWLock.RLock()
	defer psr.mRWLock.RUnlock()

	i := sort.Search(len(psr.mServiceList), func(i int) bool { return psr.mServiceList[i] >= serviceName })
	if i < len(psr.mGroupList) && psr.mGroupList[i] == serviceName {
		return true
	}

	return false
}

func (psr *PisigServiceRegistry) AddService(pisigService PisigService) (bool, error) {
	groupExists := false

	if psr.IsServiceExistsInList(pisigService.ServiceName()) {
		return false, errors.New("service already exists in the registry")
	}

	if psr.IsGroupExistsInList(pisigService.GroupName()) {
		groupExists = true
	}

	psr.mRWLock.Lock()
	defer psr.mRWLock.Unlock()

	value, ok := psr.mStore[pisigService.GroupName()]

	if !ok {
		value = make(InterfaceMap)
		psr.mStore[pisigService.GroupName()] = value
	}

	value[pisigService.ServiceName()] = pisigService

	psr.mServiceList = append(psr.mServiceList, pisigService.ServiceName())
	sort.Strings(psr.mServiceList)

	if !groupExists {
		psr.mGroupList = append(psr.mGroupList, pisigService.GroupName())
		sort.Strings(psr.mGroupList)
	}

	// update all available services
	if services, ok := psr.mAllAvailableServices[pisigService.GroupName()]; ok {
		psr.mAllAvailableServices[pisigService.GroupName()] = append(services, pisigService.ServiceName())
	} else {
		psr.mAllAvailableServices[pisigService.GroupName()] = []string{pisigService.ServiceName()}
	}

	return true, nil
}

func (psr *PisigServiceRegistry) IsServiceExists(groupName string, serviceName string) bool {
	psr.mRWLock.RLock()
	defer psr.mRWLock.RUnlock()

	group, ok := psr.mStore[groupName]

	if !ok {
		return false
	}

	_, ok1 := group[serviceName]

	if !ok1 {
		return false
	}

	return true
}

func (psr *PisigServiceRegistry) GetService(groupName string, serviceName string) (PisigService, error) {
	psr.mRWLock.RLock()
	defer psr.mRWLock.RUnlock()

	group, ok := psr.mStore[groupName]

	if !ok {
		return nil, errors.New("service group does not exists")
	}

	s, ok1 := group[serviceName]

	if !ok1 {
		return nil, errors.New("service does not exists")
	}

	return s.(PisigService), nil
}

func NewPisigServiceRegistry() *PisigServiceRegistry {
	return &PisigServiceRegistry{
		mStore:                make(map[string]InterfaceMap),
		mTopicListener:        make(map[string][]interface{}),
		mDefaultServiceMap:    make(map[string]string),
		mAllAvailableServices: make(map[string][]string),
		mRWLock:               &sync.RWMutex{},
	}
}
