package main

import (
	"fmt"
	"net/rpc"
	"net/http"
)

type Args struct{
	BuyTickets int
}

type Mutex struct{
	value int
}

func (t *Mutex) Decrease(args *Args, reply *int) error {
	t.value -= args.BuyTickets
	*reply = t.value
	return nil
}

func main(){
	conf := ReadConfig()
	
	arith := new(Mutex)
	arith.value = 1000
	
	rpc.Register(arith)
	rpc.HandleHTTP()
	
	err := http.ListenAndServe(conf.Self, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
