package main

type LamportClock struct {
	logicalClock    int64
	procId          int
}

func NewLamportClock(procId int) *LamportClock {
	lamportClock := LamportClock{1, procId}
	return &lamportClock
}

//func (this *LamportClock)
