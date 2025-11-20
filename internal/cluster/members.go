package cluster

import (
	"sync"
	"errors"
	"slices"
)

var (
	pGroup = []string{}
        lock sync.Mutex
)

func AddMember(addr string) error {
	lock.Lock()
	defer lock.Unlock()

	if addr == "" {
		return errors.New("empty node address")
	}

	pGroup = append(pGroup, addr)
	return nil
}

func RemoveMember(addr string) error {
	lock.Lock()
	defer lock.Unlock()

	if addr == "" {
		return errors.New("empty node address")
	}

	for index := range pGroup {
		if pGroup[index] == addr {
			if index < len(pGroup) {
				pGroup = slices.Delete(pGroup, index, index + 1)
				break
			}
		}
	}

	return nil
}

func GetFullMembershipList() []string  {
	lock.Lock()
	defer lock.Unlock()
	return pGroup
}
