package main

import (
	"bufio"
	"fmt"
	"gopkg.in/square/go-jose.v1/json"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"time"
)

func printUsage() {
	fmt.Println("To buy tickets, enter:  buy/b [amount of tickets]")
	fmt.Println("e.g    buy 5")
	fmt.Println()
	fmt.Println("To show the datacener status: show")
	fmt.Println()
	fmt.Println("To change cluster configuration: change [dc1] [dc2] ...")
	fmt.Println()
	fmt.Println("To exit, enter: e/exit/q/quit")
	fmt.Println()
	fmt.Println("For help, enter: help/h")
}

func handleUserInput(command string) {

	// Parse a command from user
	tokens := strings.Fields(command)

	if len(tokens) == 0 {
		return
	}

	switch len(tokens) {
	case 1:
		{
			switch tokens[0] {
			case "h":
				fallthrough
			case "help":
				printUsage()
			case "show":
				showStatus()
			case "c":
				servers := server.NewConfig
				changeConfig(servers)
			case "e":
				fallthrough
			case "exit":
				fallthrough
			case "q":
				fallthrough
			case "quit":
				os.Exit(0)
			default:
				printUsage()
			}
		}
	case 2:
		{
			switch tokens[0] {
			case "b":
				fallthrough
			case "buy":
				amount, err := strconv.ParseInt(tokens[1], 10, 32)
				if err != nil {
					printUsage()
					break
				}
				log.Print("Sent BUY TICKET request to the data center.")
				log.Print("Waiting for the datacenter's reply....")
				buyTicket(int(amount))
			case "c":
				fallthrough
			case "change":
				confChange(tokens[1:])
			default:
				printUsage()
			}

		}
	default:
		{
			fmt.Println("change")

			switch tokens[0] {
			case "c":
				fallthrough
			case "change":
				confChange(tokens[1:])
			default:
				printUsage()
			}
		}

	}
}

type BuyTicketRequest struct {
	NumTickets int
}

type BuyTicketReply struct {
	Success bool
	Remains int
}

func confChange(tokens []string) {
	fmt.Println(tokens)
	var newConf []Peer

	for _, token := range tokens{
		if val, ok := ConfigMap[token]; ok {
			newConf = append(newConf, val)
		} else {
			fmt.Println("Wrong data center name!")
			return
		}
	}
	fmt.Println(newConf)
	changeConfig(newConf)
}

func buyTicket(amount int) {
	// Synchronous call
	args := BuyTicketRequest{NumTickets: amount}
	reply := new(BuyTicketReply)
	err := rpcClient.Call("ClientComm.BuyTicketHandler", args, &reply)
	if err != nil {
		rpcClient = tryToConnectToServer("tcp", server)
	}

	if reply.Success {
		fmt.Printf("You have successfully bought %v tickets.\n", amount)
	} else {
		fmt.Println("Failed. No enough tickets to buy.")
	}
	fmt.Println("Remaining tickets:", reply.Remains)
	//time.Sleep(100 * time.Millisecond)
}

// Configuration change request, reply, and func
type ChangeConfigRequest struct {
	Servers []byte
}

type ChangeConfigReply struct {
	Success bool
}

func changeConfig(servers []Peer) {

	serversJson, err := json.Marshal(servers)
	if err != nil {
		return
	}

	args := ChangeConfigRequest{Servers: serversJson}
	reply := new(ChangeConfigReply)
	err = rpcClient.Call("ClientComm.ChangeConfigHandler", args, &reply)
	if err != nil {
		log.Println(err)
		rpcClient = tryToConnectToServer("tcp", server)
	}

	if reply.Success {
		fmt.Println("Configuration successfully changed.")
	} else {
		fmt.Println("Configuration change failed")
	}
	//fmt.Println("Remaining tickets:", reply.Remains)
}

type ShowStatusRequest struct {
}

type ShowStatusReply struct {
	NumTickets int
	Logs       []LogEntry
}

type LogEntry struct {
	Num                   int `json:"value"`
	Term                  int `json:"term"`
	IsConfigurationChange bool
	NewConfig             string
}

func showStatus() {

	args := ShowStatusRequest{}
	reply := new(ShowStatusReply)
	err := rpcClient.Call("ClientComm.ShowStatusHandler", args, &reply)
	if err != nil {
		log.Println(err)
		rpcClient = tryToConnectToServer("tcp", server)
	}

	fmt.Println("Remaining Tickets:", reply.NumTickets)
	if len(reply.Logs) == 0 {
		fmt.Println("[]")
	} else {
		for _, log := range reply.Logs {
			jsonlog, err := json.Marshal(log)
			if err != nil {
				fmt.Println("Parsing log error")
				fmt.Println(err)
				return
			}
			fmt.Println(string(jsonlog))
		}
	}
}

func waitUserInput() {
	printUsage()
	for {
		// command line user interface
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		command, _ := reader.ReadString('\n')
		handleUserInput(command)
	}
}

var rpcClient *rpc.Client
var server Server

func tryToConnectToServer(protocol string, server Server) *rpc.Client {
	var client *rpc.Client
	var err error
	var i int
	for i = 0; i < server.MaxAttempts; i++ {
		client, err = rpc.DialHTTP(protocol, server.Address)
		if err != nil {
			log.Println("dialing:", err.Error()+", retrying...")
			time.Sleep(1000 * time.Millisecond)
		} else {
			break
		}
	}

	if i == server.MaxAttempts {
		log.Fatal("Maximum attempts, cannot connect to the server")
	}
	return client
}

func init() {
	server = ReadConfig()

	rpcClient = tryToConnectToServer("tcp", server)
}

func main() {

	fmt.Println("ServerAddress:", server.Address)
	fmt.Println()

	fmt.Println("Starting Ticket Services")
	fmt.Println("Services started, please enter your command")
	waitUserInput()
}
