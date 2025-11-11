package cache_test


import (
	"testing"
	memory "log-b/internal/cache"
)

func TestFetchBucket(t *testing.T) {
	var cache = memory.Bcache{}
	cache.OpenDB()

	cache.SetBucket("testKey", "testValue")

	value := cache.FetchBucket("testKey")
	if value != "testValue" {
		t.Fail()
	}

	cache.CloseDB()
}


func TestDeleteBucket(t *testing.T) {
	var cache = memory.Bcache{}
	cache.OpenDB()

	cache.SetBucket("testK", "testV")
	err := cache.DeleteBucket("testK")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	cache.CloseDB()
}
