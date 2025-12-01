package main

import (
	"net/http"
	"os"
	"log"
	"time"
	"log-b/model"
	"encoding/json"
	"bytes"
)

func main() {
	var message  = model.BasicMessage{
		Endpoint: "set-data",
		Key:      "foo",
		Value:    "bar",
	}

	process := http.Client{Timeout: 3 * time.Second}
	// contact seed node for semplicity
	data, mErr := json.Marshal(message)
	if mErr != nil {
		log.Fatalf("%s", mErr.Error())
		os.Exit(1)
	}
	req, err := http.NewRequest(
		http.MethodPost, 
		"http://localhost:6767/set-data", 
		bytes.NewBuffer(data),
	)
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
