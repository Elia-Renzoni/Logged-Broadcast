package main

import (
	"net/http"
	"os"
	"log"
	"time"
)

func main() {
	process := http.Client{Timeout: 3 * time.Second}
	// contact seed node for semplicity
	req, err := http.NewRequest(http.MethodPost, "http://localhost:6767/set-data", nil)
	if err != nil {
		log.Fatalf("%s", err.Error())
		os.Exit(1)
	}
	res, clientErr := process.Do(req)
	if clientErr != nil {
		log.Fatalf("%s", clientErr.Error())
		os.Exit(1)
	}
	log.Print(res)
}
