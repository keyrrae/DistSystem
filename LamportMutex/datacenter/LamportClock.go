package main

type LamportClock struct {
	LogicalClock int64  `json:"logical_clock"`
	ProcId       int    `json:"proc_id"`
}

func NewLamportClock(procId int) *LamportClock {
	lamportClock := LamportClock{1, procId}
	return &lamportClock
}

//func (this *LamportClock)
