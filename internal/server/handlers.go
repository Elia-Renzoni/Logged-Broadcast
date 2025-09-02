package server

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
	"log-b/internal/broadcaster"
	"log-b/internal/cache"
	"log-b/internal/cluster"
	"log-b/model"
	"net/http"
	"log-b/internal/db"
)

const (
	POST_ADD_NODE string = "/join"
	POST_SET_BUCKET string = "/addbk"
	DELETE_BUCKET string = "/delbk"
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
			broadcaster.DoBroadcast(body, ADD_NODE, POST_ADD_NODE)
		},
	)
}

func setKVBucket(volatileBucketer cache.MemoryCache, buffer db.Storage) http.Handler {
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

			dbErr := buffer.WriteMessage(model.PersistentMessage{
				Sinfo: model.Sender{Addr: r.RemoteAddr, Port: "undefined"},
				Cinfo: model.MessageContent{Endpoint: SET_DATA, Key: msg.Key, Value: msg.Value},
			}, 0)
			if dbErr != nil {
				nack(w, dbErr)
				return
			}
			ack(w, data)

			majorityReached := broadcaster.DoBroadcast(body, SET_DATA, POST_SET_BUCKET)
			
			// if the majority quorum is reached
			// change the status to DELIVERED
			if majorityReached {
				// TODO -> change status
			}
		},
	)
}

func removeKvBucket(volatileBucketer cache.MemoryCache, buffer db.Storage) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			splitted := strings.Split(r.URL.Path, "/")
			param := splitted[2]
			
			if err := volatileBucketer.DeleteBucket(param); err != nil {
				nack(w, err)
				return
			}
			majorityReached := broadcaster.DoBroadcast(nil, DELETE_DATA, DELETE_BUCKET)
			if majorityReached {
				// delete from persistent storage...
			}
			ack(w, []byte("Bucket Succesfully Removed!"))
		},
	)
}

func fetchKvBucket(volatileBucketer cache.MemoryCache) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			splitted := strings.Split(r.URL.Path, "/")
			param := splitted[2]

			if value := volatileBucketer.FetchBucket(param); value != "" {
				ack(w, []byte(value))
				return
			}
			nack(w, errors.New("Key Not Found!"))
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