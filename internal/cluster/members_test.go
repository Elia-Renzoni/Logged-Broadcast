package cluster_test


import (
	"testing"
	"log-b/internal/cluster"
	"slices"
)

func TestAddMember(t *testing.T) {
	cluster.AddMember("127.0.0.1:7979")
	cluster.AddMember("127.0.0.1:8080")
	cluster.AddMember("127.0.0.1:8081")

	var list = cluster.GetFullMembershipList()
	ok := slices.Contains(list, "127.0.0.1:7979")
	if !ok {
		t.Fail()
	}
}

func TestRemoveMember(t *testing.T) {
	cluster.AddMember("127.0.0.1:6767")
	cluster.AddMember("127.0.0.1:5400")

	if err := cluster.RemoveMember("127.0.0.1:6767"); err != nil {
		t.Fatalf("%s", err.Error())
	}
	var list = cluster.GetFullMembershipList()
	ok := slices.Contains(list, "127.0.0.1:6767")
	if ok {
		t.Fail()
	}
}
