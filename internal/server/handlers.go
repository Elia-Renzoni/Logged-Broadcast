package server

import (
	"encoding/json"
	"errors"
	"io"
	"log-b/internal/broadcaster"
	"log-b/internal/cache"
	"log-b/internal/cluster"
	"log-b/model"
	"net/http"
)

func addNodeToCluster() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var msg model.BasicMessage
			defer r.Body.Close()

			body, err := io.ReadAll(r.Body)
			if err != nil {
				nack(w, err)
				return
			}

			if err := json.Unmarshal(body, &msg); err != nil {
				nack(w, err)
				return
			}

			cluster.AddMember(msg.Node)
			// maybe is not useful...
			broadcaster.DoBroadcast(body, ADD_NODE)
		},
	)
}

func setKVBucket(volatileBucketer cache.MemoryCache) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close() 
			
			body, err := io.ReadAll(r.Body)
			if err != nil {
				nack(w, err)
				return
			}

			var msg model.BasicMessage
			if jsonErr := json.Unmarshal(body, &msg); jsonErr != nil {
				nack(w, jsonErr)
				return
			}

			if msg.Key == "" || msg.Value == "" {
				nack(w, errors.New("Empty Payload Elements!"))
				return
			}

			volatileBucketer.SetBucket(msg.Key, msg.Value)
			data, maErr := json.Marshal(model.BasicPositiveAck{Succ:"executed SET operation!"})
			if maErr != nil {
				nack(w, maErr)
				return
			}
			ack(w, data)
		},
	)
}

func removeKvBucket(volatileBucketer cache.MemoryCache) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
		},
	)
}

func fetchKvBucket(volatileBucketer cache.MemoryCache) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
		},
	)
}

func nack(w io.Writer, err error) {
	rw, ok := w.(http.ResponseWriter)
	if ok {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
	}
	payload, err := json.Marshal(model.BasicError{
		Error: err.Error(),
	})

	if err != nil  {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(payload)
}

func ack(w io.Writer, payload []byte) {
	rw, ok := w.(http.ResponseWriter)
	if ok {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)
	}

	w.Write(payload)
}
