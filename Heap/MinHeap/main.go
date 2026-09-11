package main

import (
	"container/heap"
	"fmt"
)

type Item struct {
	Name string
	Priority int
}

type MinHeap []Item

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Less(i, j int) bool {
	if (h[i].Priority == h[j].Priority) {
		return h[i].Name < h[j].Name
	}
	return h[i].Priority < h[j].Priority
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(item any) {
	*h = append(*h, item.(Item)) 
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)

	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	minHeap := &MinHeap{}

	heap.Init(minHeap)

	heap.Push(minHeap, Item{Name: "John", Priority: 1})
	heap.Push(minHeap, Item{Name: "Abby", Priority: 1})
	heap.Push(minHeap, Item{Name: "Tony", Priority: 3})

	for range len(*minHeap) {
		fmt.Println(heap.Pop(minHeap).(Item))
	}
}