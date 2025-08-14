package server

import (
	"net/http"
)

func addNodeToCluster() http.Handler {
	return http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
		},
	)
}

func setKVBucket() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
		},
	)
}

func removeKvBucket() http.Handler {
	return http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
		},
	)
}

func fetchKvBucket() http.Handler {
	return http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
		},
	)
}
