package main

import (
	"container/heap"
	"fmt"
	"log"
	"net/rpc"
	"time"
)

type Args struct {
	BuyTickets int
}

type ClientComm struct {
	value int
}

const (
	ASK     = "ASK"
	RELEASE = "RELEASE"
)

func delay(){
	if conf.Delay == 0{
		return
	}
	time.Sleep(time.Duration(conf.Delay) * time.Second)
}

func (t *ClientComm) BuyTicketRequest(args *Args, reply *int) error {
	// simulating a message delay
	//delay()
	log.Println()
	log.Print("Received a BUY TICKET request from a client")
	
	// piggybacking time stamp with request
	changeRequest := &Request{
		Request: args.BuyTickets,
		Clock:   LamportClock{lamClock.LogicalClock, lamClock.ProcId},
	}

	// places the request on waitQueue
	heap.Push(&waitQueue, changeRequest)

	// sends request message to all sites
	// asking for permission to enter critical section
	req := DataCenterRequest{
		RequestType: ASK,
		RequestBody: *changeRequest,
	}

	count := 0
	done := make(chan bool)

	go func() {
		for {
			// has received messages with larger timestamps from all other sites
			allOtherSitesAgree := count == conf.NumOfServers()
			// request is at the top of waitQueue
			requestAtTop := changeRequest.Clock.equalsTo(waitQueue.Peek().Clock)

			if allOtherSitesAgree && requestAtTop {
				done <- true
				break
			}
			// TODO: replace this wait loop with channel

			time.Sleep(100 * time.Millisecond)
		}
	}()

	for _, server := range connections {

		dataCenterReply := new(DataCenterReply)
		divCall := server.Go("DataCenterComm.CriticalSectionRequest", req, dataCenterReply, nil)

		go func() {
			replyCall := <-divCall.Done
			gotReply := replyCall.Reply.(*DataCenterReply)

			lamClock.LogicalClock = max(gotReply.TimeStamp.LogicalClock, lamClock.LogicalClock) + 1
			log.Printf("Received a reply from %v with logical clock %v. My clock now: %v\n",
				gotReply.TimeStamp.ProcId, gotReply.TimeStamp.LogicalClock, lamClock.LogicalClock)
			
			// simulate delay of release request
			delay()

			// increase the counter if got larger time stamp
			if gotReply.TimeStamp.largerThan(changeRequest.Clock) {
				count++
			}

		}()
	}

	<-done // block the main thread

	var releaseDecAmount int
	if conf.RemainingTickets < changeRequest.Request {
		releaseDecAmount = 0
	} else {
		releaseDecAmount = changeRequest.Request
	}

	heap.Pop(&waitQueue)
	releaseRequest := &Request{
		Request: releaseDecAmount,
		Clock:   LamportClock{lamClock.LogicalClock, lamClock.ProcId},
	}

	conf.RemainingTickets -= releaseDecAmount
	*reply = conf.RemainingTickets
	count = 0

	go func() {
		for {
			// has received messages with larger timestamps from all other sites
			allOtherSitesAgree := count == conf.NumOfServers()

			if allOtherSitesAgree {
				done <- true
				break
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()

	for _, server := range connections {

		dataCenterReply := new(DataCenterReply)
		req = DataCenterRequest{
			RequestType: RELEASE,
			RequestBody: *releaseRequest,
		}
		releaseCall := server.Go("DataCenterComm.CriticalSectionRequest", req, dataCenterReply, nil)

		go func() {
			_ = <-releaseCall.Done
			delay()
			count++
		}()
	}
	<-done
	return nil
}

type DataCenterComm struct {
	value int
}

type DataCenterRequest struct {
	RequestType string
	RequestBody Request
}

type DataCenterReply struct {
	TimeStamp LamportClock
}

func max(a int64, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func (t *DataCenterComm) CriticalSectionRequest(req *DataCenterRequest, reply *DataCenterReply) error {
	
	// Simulating request delay
	delay()
	
	// upon receives a request, increase the logic clock
	lamClock.LogicalClock = max(req.RequestBody.Clock.LogicalClock, lamClock.LogicalClock) + 1

	fmt.Println()
	
	var article string = "a"
	if req.RequestType == ASK {
		article = "an"
	}
	
	log.Printf("Received %v %v request from process %v with logical clock %v. My clock now: %v\n",
		article,
		req.RequestType, req.RequestBody.Clock.ProcId,
		req.RequestBody.Clock.LogicalClock,
		lamClock.LogicalClock)

	switch req.RequestType {
	case ASK:
		// receives a request asking to enter critical section
		// reply with timestamp of this site
		reply.TimeStamp.LogicalClock = lamClock.LogicalClock
		reply.TimeStamp.ProcId = lamClock.ProcId

		// push request to waitQueue
		heap.Push(&waitQueue, &(req.RequestBody))
	case RELEASE:
		heap.Pop(&waitQueue)
		conf.RemainingTickets -= req.RequestBody.Request
		fmt.Println()
		fmt.Print("> ")
	default:

	}
	delay()
	if req.RequestType == ASK {
		log.Printf("Replied to process %v. My clock now: %v\n", req.RequestBody.Clock.ProcId, lamClock.LogicalClock)
	}
	return nil
}

var connectionStatus map[string]bool

// Reference to all the rpc clients
type Connections []*rpc.Client
var connections Connections

// Indicator -- whether this data center has connected to all the other data centers
var allConnected bool

func EstablishConnections() {
	// Establishing connections to other data centers
	connectionCounter := 0
	allConnected = false

	connectionStatus = make(map[string]bool)
	log.Println("number of servers:", conf.Servers)
	for _, server := range conf.Servers {
		connectionStatus[server.Address] = false
	}

	var i int
	for i = 0; i < conf.MaxAttempts; i++ {
		log.Printf("%v  ", i)
		for _, server := range conf.Servers {
			log.Println(connectionStatus)
			if connectionStatus[server.Address] {
				continue
			}

			client, err := rpc.DialHTTP("tcp", server.Address)
			if err != nil {
				log.Println("dialing:", err.Error()+", retrying...")
			} else {
				connectionStatus[server.Address] = true
				connectionCounter++
				connections = append(connections, client)
				break
			}
		}
		if connectionCounter == conf.NumOfServers() {
			fmt.Println("Number of clients:", conf.NumOfServers())
			allConnected = true
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}

	if i == conf.MaxAttempts {
		log.Fatal("Maximum attempts: cannot reach all the datacenters")
	}
}
