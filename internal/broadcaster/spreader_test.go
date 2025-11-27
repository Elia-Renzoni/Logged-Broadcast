package broadcaster_test


import (
	"testing"
	"log-b/internal/cluster"
	"log-b/internal/broadcaster"
	"net/http"
	"time"
)


func formCluster() {
	e1 := cluster.AddMember("127.0.0.1:8080")
	e2 := cluster.AddMember("127.0.0.1:8081")
	e3 := cluster.AddMember("127.0.0.1:8082")
	e4 := cluster.AddMember("127.0.0.1:8083")

	switch {
	case e1 != nil, e2 != nil, e3 != nil, e4 != nil:
		return
	}
}

func TestDoBroadcast(t *testing.T) {
	formCluster()
	go startServers()
	time.Sleep(1 * time.Second)
	result := broadcaster.DoBroadcast([]byte("mockbytes"), broadcaster.SET_DATA)
	
	// fail the test if the majority quorum is 
	// not reached
	if !result {
		t.Fail()
	}
}

func startServers() {
	group := cluster.GetFullMembershipList()
	for _, node := range group {
		go func(peer string) {
			mux := http.NewServeMux()
			mux.HandleFunc("/set-data", handleMessage)

			err := http.ListenAndServe(peer, mux)
			if err != nil {
				return
			}
		}(node)
	}
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
