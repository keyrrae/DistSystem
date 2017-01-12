package main

type LamportClock struct{
	clock, proc_id int64
}

func NewLamportClock(proc_id int64) *LamportClock{
	lamportClock := LamportClock{1, proc_id}
	return &lamportClock
}

//func (this *LamportClock)