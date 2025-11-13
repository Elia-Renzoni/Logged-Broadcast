package cluster_test


import (
	"testing"
	"log-b/internal/cluster"
	"sync"
	"net/http"
	"fmt"
	"time"
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

	backoffRetry += 1
	mu.Unlock()
	time.Sleep(4 * time.Second)
}

func TestRegisterToSeed(t *testing.T) {
	go startMockSeed()
	time.Sleep(1 * time.Second)

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()

		cluster.RegisterToSeed(mockSeedAddress, "my-address")
	}()
	mu.Lock()
	if backoffRetry != backOffMaxRetries {
		t.Fail()
	}
	mu.Unlock()
}
