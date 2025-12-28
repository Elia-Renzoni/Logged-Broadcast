package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"log-b/model"
	"net/http"
	"os"
	"time"
)

func main() {
	var message = model.BasicMessage{
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
	resData, e := io.ReadAll(res.Body)
	if e != nil {
		os.Exit(1)
	}
	log.Print(string(resData))
	res.Body.Close()

	getReq, getErr := http.NewRequest(
		http.MethodGet,
		"http://localhost:6767/get-data/foo",
		nil,
	)
	if getErr != nil {
		log.Fatal(getErr)
		os.Exit(1)
	}
	getRes, getResErr := process.Do(getReq)
	if getResErr != nil {
		log.Fatal(getResErr)
		os.Exit(1)
	}
	getData, getDataErr := io.ReadAll(getRes.Body)
	if getDataErr != nil {
		log.Fatal(getDataErr)
		os.Exit(1)
	}

	log.Print(string(getData))
	getRes.Body.Close()

	request, erro := http.NewRequest(
		http.MethodDelete,
		"http://localhost:6767/delete-data/foo",
		nil,
	)
	if erro != nil {
		log.Fatalf("%s", erro.Error())
		os.Exit(1)
	}
	response, clientErro := process.Do(request)
	if clientErro != nil {
		log.Fatalf("%s", clientErro.Error())
		os.Exit(1)
	}
	respData, resErr := io.ReadAll(response.Body)
	if resErr != nil {
		log.Fatalf("%s", resErr.Error())
		os.Exit(1)
	}
	log.Print(string(respData))
	response.Body.Close()
}
