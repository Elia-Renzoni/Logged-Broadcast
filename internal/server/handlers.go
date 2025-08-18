package server

import (
	"encoding/json"
	"io"
	"ioutil"
	"log-b/broadcaster"
	"log-b/cache"
	"log-b/cluster"
	"log-b/model"
	"net/http"
)

func addNodeToCluster() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var msg model.BasicMessage
			defer r.Body.Close()

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				nack(w)
				return
			}

			if err := json.Unmarhsal(body, &msg); err != nil {
				nack(w)
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

func nack(w io.Writer) {

}
