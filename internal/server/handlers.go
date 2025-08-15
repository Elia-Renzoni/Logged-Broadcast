package server

import (
	"net/http"
	"log-b/cache"
)

func addNodeToCluster() http.Handler {
	return http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
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
		func (w http.ResponseWriter, r *http.Request) {
		},
	)
}

func fetchKvBucket(volatileBucketer cache.MemoryCache) http.Handler {
	return http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
		},
	)
}
