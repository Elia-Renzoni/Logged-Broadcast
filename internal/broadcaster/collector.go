package broadcaster

import (
	"log-b/cluster"
)

type ackCounter int

func (a *ackCounter) inc() {
	*a += 1
}

func (a ackCounter) isMajorityQourumReached() bool {
	group := cluster.GetFullMembershipList()
	quorum := len(group) / 2

	if a > quorum {
		return true
	}
	return false
}
