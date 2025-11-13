package cluster_test


import (
	"testing"
	"log-b/internal/cluster"
	"sync"
	"net/http"
)

const mockSeedAddress string = "127.0.0.1:5006"
const backOffMaxRetries int  = 7

var backoffRetry int = 0
var mu sync.Mutex

func startMockSeed() {
	mux := http.NewServeMux()
	mux.HandleFunc("/add-node", handleConn)

	http.ListenAndServe(mockSeedAddress, mux)
}

func handleConn(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	backoffRetry += 1
}

func TestRegisterToSeed(t *testing.T) {
	go startMockSeed()
	cluster.RegisterToSeed(mockSeedAddress, "my-address")

	defer func() {
		if r := recover(); r != nil {
			mu.Lock()
			if backoffRetry != backOffMaxRetries {
				t.Fail()
			}
			mu.Unlock()
		}
	}()
}
