package cache

import (
	"log"

	lru "github.com/hashicorp/golang-lru"
)

var Cache *lru.Cache

func InitCache(size int) {
	var err error
	Cache, err = lru.New(size)
	if err != nil {
		log.Fatal("Failed to create cache: &v", err)
	}
}
