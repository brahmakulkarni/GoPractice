package main

import (
	"container/heap"
	"fmt"
)

type Item struct {
	Name string
	Priority int
}

type MaxHeap []Item

func (h MaxHeap) Len() int {
	return len(h)
}

func (h MaxHeap) Less(i, j int) bool {
	if h[i].Priority == h[j].Priority {
		return h[i].Name < h[j].Name
	}
	return h[i].Priority > h[j].Priority
}

func (h MaxHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MaxHeap) Push(item any) {
	*h = append(*h, item.(Item))
}

func (h *MaxHeap) Pop() any {
	old := *h
	n := len(old)

	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	maxHeap := &MaxHeap{}

	heap.Init(maxHeap)

	heap.Push(maxHeap, Item{Name: "John", Priority: 1})
	heap.Push(maxHeap, Item{Name: "Abby", Priority: 1})
	heap.Push(maxHeap, Item{Name: "Tony", Priority: 3})

	for range len(*maxHeap) {
		fmt.Println(heap.Pop(maxHeap).(Item))
	}
}