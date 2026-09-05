package main

import (
	"fmt"
)

type LRUCache struct {
	capacity int
	cache map[int]int
	order []int
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache: make(map[int]int),
	}
}

func (lru *LRUCache) get(key int) int {
	val, exists := lru.cache[key]
	if !exists {
		return -1
	}

	lru.markUsed(key)
	return val
}

func (lru *LRUCache) put(key, value int) {
	if _, exists := lru.cache[key]; exists {
		lru.cache[key] = value
		lru.markUsed(key)
		return
	}
	if len(lru.cache) == lru.capacity {
		removeKey := lru.order[0]
		lru.order = lru.order[1:]
		delete(lru.cache, removeKey)
	}
	lru.cache[key] = value
	lru.order = append(lru.order, key)
}

func (lru *LRUCache) markUsed(key int) {
	for ind, k := range lru.order {
		if k == key {
			lru.order = append(lru.order[:ind], lru.order[ind+1:]...)
			break
		}
	}
	lru.order = append(lru.order, key)
}

func main() {
	lru := NewLRUCache(2)

	lru.put(1, 10)
	lru.put(2, 20)

	fmt.Println("get 1: ", lru.get(1))
	lru.put(3,30)

	fmt.Println("get 2: ", lru.get(2))
	fmt.Println("get 3: ", lru.get(3))

	lru.put(4,40)

	fmt.Println("get 1: ", lru.get(1))
	fmt.Println("get 3: ", lru.get(3))
	fmt.Println("get 4: ", lru.get(4))
}