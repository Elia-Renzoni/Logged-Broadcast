package broadcaster

import (
	"log-b/cluster"
	"log-b/server"
	"net/http"
)

func DoBroadcast(message []byte, methodRouter string) {
	var memberlist = cluster.GetFullMembershipList()
	var c ackCounter

	for _, node := range memberlist {
	
	}
}

func send(addr string, msg []byte, methodRouter string) {
	var res *http.Response
	switch methodRouter {
	case ADD_NODE:
		res = http.Post()
	case SET_DATA:
		res = http.Post()
	case DELETE_DATA:
		res = http.Get()
	}
}


