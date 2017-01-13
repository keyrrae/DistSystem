package main

type LamportClock struct {
	logicalClock int64  `json:"logical_clock"`
	procId       int    `json:"proc_id"`
}

func NewLamportClock(procId int) *LamportClock {
	lamportClock := LamportClock{1, procId}
	return &lamportClock
}

//func (this *LamportClock)
