package broadcaster

import (
	"log-b/internal/cluster"
)

type ackCounter int

func (a *ackCounter) inc() {
	*a += 1
}

func (a ackCounter) isMajorityQuorumReached() bool {
	group := cluster.GetFullMembershipList()
	var quorum = len(group) / 2

	return int(a) > quorum
}
