package server

import (
	"net/http"
	"net/url"
	"log-b/cache"
)

const (
	ADD_NODE = "/add-node"
	SET_DATA = "/set-data"
	DELETE_DATA = "/delete-data"
	GET_DATA = "/get-data"
)

type Router map[string]http.Handler

func (ro Router) ServeRequest(w http.ResponseWriter, r *http.Request) {
	endpoint := getEndpoint(r.URL)
	if ok := isEndpointLegit(endpoint); !ok {
		http.Error(w, "Invalid API Endpoint", http.StatusBadRequest)
		return
	}

	matchedHandler, ok := ro[endpoint]
	if !ok {
		http.Error(w, "The RPC Match is Not Possible!", http.StatusInternalServerError)
		return
	}
	matchedHandler.ServeHTTP(w, r)
}

func InitRouter(bucketer cache.MemoryCache) Router {
	return Router{
		ADD_NODE: addNodeToCluster(),
		SET_DATA: setKVBucket(bucketer),
		DELETE_DATA: removeKvBucket(bucketer),
		GET_DATA: fetchKvBucket(bucketer),
	}
}

func getEndpoint(u *url.URL) string {
	return u.Path
}

func isEndpointLegit(endpoint string) bool {
	switch endpoint {
		case ADD_NODE, SET_DATA, DELETE_DATA, GET_DATA:
			return true
	}
	return false
}
