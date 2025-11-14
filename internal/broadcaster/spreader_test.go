package broadcaster_test


import (
	"testing"
	"log-b/internal/cluster"
	"log-b/internal/broadcaster"
	"net/http"
	"time"
)


func formCluster() {
	cluster.AddMember("127.0.0.1:8080")
	cluster.AddMember("127.0.0.1:8081")
	cluster.AddMember("127.0.0.1:8082")
	cluster.AddMember("127.0.0.1:8083")
}

func TestDoBroadcast(t *testing.T) {
	formCluster()
	go startServers()
	time.Sleep(1 * time.Second)
	result := broadcaster.DoBroadcast([]byte("mockbytes"), broadcaster.SET_DATA, "/addbk")
	
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
			mux.HandleFunc("/addbk", handleMessage)

			http.ListenAndServe(peer, mux)
		}(node)
	}
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
