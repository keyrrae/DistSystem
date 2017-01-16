package main

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
