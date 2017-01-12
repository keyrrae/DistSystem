package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
	"strings"
	"net/rpc"
	"log"
)

func printUsage() {
	fmt.Println("To buy tickets, enter:")
	fmt.Println("       buy [amount of tickets]")
	fmt.Println("e.g    buy 5")
}


func newRPCclient(protocol string, address string) *rpc.Client {
	client, err := rpc.DialHTTP(protocol, address)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	
	return client
}

func handleUserInput(command string) {
	
	// Parse a command from user
	tokens := strings.Fields(command)
	fmt.Println(tokens)
	
	if len(tokens) == 0 || len(tokens) > 2 {
		return
	}
	
	switch tokens[0] {
	case "h":
		fallthrough
	case "help":
		printUsage()
	case "b":
		fallthrough
	case "buy":
		amount, err := strconv.ParseInt(tokens[1], 10, 32)
		if err != nil {
			printUsage()
			break
		}
		buyTicket(amount)
	case "e":
		fallthrough
	case "exit":
		fallthrough
	case "q":
		fallthrough
	case "quit":
		os.Exit(0)
	}
}


func buyTicket(amount int){
	// Synchronous call
	args := Args{7}
	var reply int
	err := rpcClient.Call("Mutex.Decrease", args, &reply)
	if err != nil {
		log.Fatal("Error:", err)
	}
	fmt.Println("Remaining tickets:", reply)
	//time.Sleep(100 * time.Millisecond)
}

func waitUserInput() {
	for {
		// command line user interface
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		command, _ := reader.ReadString('\n')
		handleUserInput(command)
		time.Sleep(80 * time.Millisecond)
	}
}

type Args struct{
	BuyTickets int
}

var rpcClient *rpc.Client
var serverAddress string

func init(){
	serverAddress = ReadConfig()
	
	rpcClient = newRPCclient("tcp", serverAddress)
}

func main() {
	
	
	fmt.Println("ServerAddress:", serverAddress)
	fmt.Println()

	fmt.Println("Starting Ticket Services")
	fmt.Println("Services started, please enter your command")
	waitUserInput()
}
