package cluster

import (
	"sync"
	"errors"
	"slices"
)

var pGroup = []string{}
var lock sync.Mutex

func AddMember(addr string) error {
	lock.Lock()
	defer lock.Unlock()

	if addr == "" {
		return errors.New("Empty Node Address!")
	}

	pGroup = append(pGroup, addr)
	return nil
}

func RemoveMember(addr string) error {
	lock.Lock()
	defer lock.Unlock()

	if addr == "" {
		return errors.New("Empty Node Address!")
	}

	for index := range pGroup {
		pGroup = slices.Delete(pGroup, index, index + 1)
	}

	return nil
}

func GetFullMembershipList() []string  {
	return pGroup
}
