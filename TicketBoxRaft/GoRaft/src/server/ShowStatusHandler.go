package main

type ShowStatusRequest struct {

}

type ShowStatusReply struct {
	NumTickets        int  // current term, for candidate to update itself
	Logs []LogEntry // true means candidate received vote
}

func (t *ClientComm) ShowStatusHandler(req *ShowStatusRequest, reply *ShowStatusReply) error {
	reply.NumTickets = self.StateParam.RemainingTickets
	for i := 0; i <= self.StateParam.CommitIndex; i++ {
		reply.Logs = append(reply.Logs, self.StateParam.Logs[i])
	}
	return nil
}
