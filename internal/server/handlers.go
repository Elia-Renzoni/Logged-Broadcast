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
	"fmt"
)

func addNodeToCluster() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var msg model.BasicJoinMessage
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

			fmt.Println(msg.NodeCompleteAddress)
			length := len(msg.NodeCompleteAddress)
			if length > 1 {
				if err := cluster.AddMembers(msg.NodeCompleteAddress); err != nil {
					nack(w, err)
					return
				}
			} else {
				if err := cluster.AddMember(msg.NodeCompleteAddress[length - 1]); err != nil {
					nack(w, err)
					return
				}
			}


			if msg.SeedFlagOption {
				dataToSpread := body
				if ok := cluster.HasMoreElements(); ok {
					list := cluster.GetFullMembershipList()
					dataToSpread, _ = json.Marshal(model.BasicJoinMessage{
						SeedFlagOption:      false,
						NodeCompleteAddress: list,
					})
				}
				majorityReached := broadcaster.DoBroadcast(dataToSpread, ADD_NODE)
				if !majorityReached {
					nack(w, errors.New("operation aborted: quorum not reached"))
					return
				}
			}

			ack(w, []byte("Join Approved"))
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
				nack(w, errors.New("empty payload elements"))
				return
			}

			volatileBucketer.SetBucket(msg.Key, msg.Value)
			data, maErr := json.Marshal(model.BasicPositiveAck{Succ:"executed SET operation!"})
			if maErr != nil {
				nack(w, maErr)
				return
			}

			dbErr := buffer.WriteMessage(model.PersistentMessage{
				Sinfo: model.Sender{
					Addr: r.RemoteAddr, 
					Port: "undefined",
				},
				Cinfo: model.MessageContent{
					Endpoint: SET_DATA, 
					Key: msg.Key,
					Value: msg.Value,
				},
			}, 0)
			if dbErr != nil {
				nack(w, dbErr)
				return
			}

			majorityReached := broadcaster.DoBroadcast(body, SET_DATA)
			if !majorityReached {
				nack(w, errors.New("operation aborted: quorum not reached"))
				return
			}

			if err := buffer.ChangeStatus(msg.Key); err != nil {
				nack(w, err)
				return
			}
			
			ack(w, data)
		},
	)
}

func removeKvBucket(volatileBucketer cache.MemoryCache, buffer db.Storage) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			splitted := strings.Split(r.URL.Path, "/")
			key := splitted[2]
			
			if err := volatileBucketer.DeleteBucket(key); err != nil {
				nack(w, err)
				return
			}
			majorityReached := broadcaster.DoBroadcast(nil, DELETE_DATA)
			if !majorityReached {
				nack(w, errors.New("opreation aborted: quorum not reached"))
				return
			}

			if err := buffer.DeleteMessage(key); err != nil {
				nack(w, err)
				return
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
			nack(w, errors.New("key Not found"))
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
