package main

import (
	"container/heap"
	"encoding/json"
	"fmt"
)

type LamportClock struct {
	LogicalClock int64 `json:"logical_clock"`
	ProcId       int   `json:"proc_id"`
}

func NewLamportClock(procId int) *LamportClock {
	lamportClock := LamportClock{1, procId}
	return &lamportClock
}

func (this *LamportClock) equalsTo(that LamportClock) bool {
	return this.LogicalClock == that.LogicalClock && this.ProcId == that.ProcId
}

func (this *LamportClock) smallerThan(that LamportClock) bool {
	if this.LogicalClock < that.LogicalClock {
		return true
	}

	if this.LogicalClock == that.LogicalClock {
		return this.ProcId < that.ProcId
	}
	return false
}

func (this *LamportClock) largerThan(that LamportClock) bool {
	return !this.equalsTo(that) && !this.smallerThan(that)
}

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

func TestWaitQueue() {
	var pq PriorityQueue
	heap.Init(&pq)

	item := &Request{
		Request: 3,
		Clock: LamportClock{
			LogicalClock: 1,
			ProcId:       3,
		},
		//Index: 0,
	}

	heap.Push(&pq, item)

	item = &Request{
		Request: 1,
		Clock: LamportClock{
			LogicalClock: 1,
			ProcId:       1,
		},
		//Index: 0,
	}

	heap.Push(&pq, item)

	item = &Request{
		Request: 2,
		Clock: LamportClock{
			LogicalClock: 1,
			ProcId:       2,
		},
		//Index: 0,
	}

	heap.Push(&pq, item)

	for _, item := range pq {
		//item := heap.Pop(&pq).(*Request)
		itemJson, _ := json.MarshalIndent(item, "", "    ")
		fmt.Println(string(itemJson))
	}

	fmt.Println("peek\n", pq.Peek())

	fmt.Println()
	fmt.Println()

	for len(pq) > 0 {
		item := heap.Pop(&pq).(*Request)
		itemJson, _ := json.MarshalIndent(item, "", "    ")
		fmt.Println(string(itemJson))
	}

}

func main() {
	TestWaitQueue()
}
