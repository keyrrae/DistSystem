package main

import (
	"log"
	"errors"
)

type BuyTicketRequest struct {
	NumTickets int
}

type BuyTicketReply struct {
	Success bool
	Remains int
}

func (t *ClientComm) BuyTicketHandler(req *BuyTicketRequest, reply *BuyTicketReply) error {
	log.Println("received buy ticket from a client")

	if self.LeaderID != self.Conf.ProcessID {

		leader := self.Conf.PeersMap[self.LeaderID]
		leaderReply := new(BuyTicketReply)

		if leader == nil || leader.Comm == nil {
			reply.Success = false
			return errors.New("No leader in the cluster yet")
		}

		err := leader.Comm.Call("DataCenterComm.BuyTicketHandler", req, &leaderReply)
		if err != nil {
			leader.Comm = nil
			leader.Connected = false
			reply.Success = false
			reply.Remains = self.StateParam.RemainingTickets
			return err
		}

		reply.Remains = leaderReply.Remains
		reply.Success = leaderReply.Success
	} else {
		sendReqToFollowers(req, reply)
	}

	return nil
}

func (t *DataCenterComm) BuyTicketHandler(req *BuyTicketRequest, reply *BuyTicketReply) error {
	log.Println("received buy ticket from a follower redirection")

	if self.LeaderID == self.Conf.ProcessID {
		sendReqToFollowers(req, reply)
	} else {
		reply.Success = false
		reply.Remains = self.StateParam.RemainingTickets
	}

	return nil
}

func sendReqToFollowers(req *BuyTicketRequest, reply *BuyTicketReply) {
	// I'm the leader
	logEntry := LogEntry{
		Num:  req.NumTickets,
		Term: self.StateParam.CurrentTerm,
	}

	if req.NumTickets > self.StateParam.RemainingTickets {
		reply.Success = false
		reply.Remains = self.StateParam.RemainingTickets
		return
	}
	self.StateParam.Logs = append(self.StateParam.Logs, logEntry)

	leaderBehavior()

	reply.Success = true
	reply.Remains = self.StateParam.RemainingTickets
}
