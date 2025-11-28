package cluster_test


import (
	"testing"
	"log-b/internal/cluster"
	"slices"
)

/*
func TestAddMember(t *testing.T) {
	e1 := cluster.AddMember("127.0.0.1:7979")
	e2 := cluster.AddMember("127.0.0.1:8080")
	e3 := cluster.AddMember("127.0.0.1:8081")
	
	switch {
	case e1 != nil, e2 != nil, e3 != nil:
		return
	}

	var list = cluster.GetFullMembershipList()
	ok := slices.Contains(list, "127.0.0.1:7979")
	if !ok {
		t.Fail()
	}
}*/

func TestAddMembers(t *testing.T) {
	membs := []string{
		"127.0.0.1:123",
		"127.0.0.1:124",
		"127.0.0.1:125",
	}
	err := cluster.AddMembers(membs)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	var list = cluster.GetFullMembershipList()
	if ok := slices.Contains(list, "127.0.0.1:125"); !ok {
		t.Fail()
	}
}

func TestAddMembersWithIdempotency(t *testing.T) {
	membs := []string{
		"127.0.0.1:126",
		"127.0.0.1:126",
	}

	err := cluster.AddMembers(membs)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if ok := cluster.HasMoreElements(); ok {
		t.Fail()
	}
}

func TestRemoveMember(t *testing.T) {
	e1 := cluster.AddMembers([]string{"127.0.0.1:6767"})
	e2 := cluster.AddMembers([]string{"127.0.0.1:5400"})

	if e1 != nil || e2 != nil {
		return
	}

	if err := cluster.RemoveMember("127.0.0.1:6767"); err != nil {
		t.Fatalf("%s", err.Error())
	}
	var list = cluster.GetFullMembershipList()
	ok := slices.Contains(list, "127.0.0.1:6767")
	if ok {
		t.Fail()
	}


	e3 := cluster.AddMembers([]string{"127.0.0.1:2222"})
	e4 := cluster.AddMembers([]string{"127.0.0.1:5555"})
	if e3 != nil || e4 != nil {
		return
	}

	if err := cluster.RemoveMember("127.0.0.1:5555"); err != nil {
		t.Fatalf("%s", err.Error())
	}

	var finalList = cluster.GetFullMembershipList()
	result := slices.Contains(finalList, "127.0.0.1:5555")
	if result {
		t.Fail()
	}
}
