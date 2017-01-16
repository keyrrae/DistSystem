package main

import (
	"container/heap"
)

// An Item is something we manage in a priority queue.
type Request struct {
	Request int          `json:"request"` // The value of the item; arbitrary.
	Clock   LamportClock `json:"clock"`   // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	Index int `json:"index"` // The index of the item in the heap.
}

func (a Request) equalsTo(b Request) bool {
	logicalClockEquality := a.Clock.LogicalClock == b.Clock.LogicalClock
	procIdEquality := a.Clock.ProcId == b.Clock.ProcId
	return logicalClockEquality && procIdEquality
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Request

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].Clock.smallerThan(pq[j].Clock)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Request)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) Peek() *Request {
	//n := len(*pq)
	//item := (*pq)[n-1]
	item := (*pq)[0]
	return item
}

// update modifies the priority and value of an Request in the queue.
func (pq *PriorityQueue) update(item *Request, value int, priority LamportClock) {
	item.Request = value
	item.Clock = priority
	heap.Fix(pq, item.Index)
}
