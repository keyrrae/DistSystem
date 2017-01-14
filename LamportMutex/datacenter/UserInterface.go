package main

import (
	"os"
	"fmt"
	"bufio"
	"time"
	"strings"
	"encoding/json"
	"container/heap"
)

func printUsage() {
	fmt.Println("pc - print configuration")
	fmt.Println("pq - print wait queue")
	fmt.Println("pv - print amount of tickets")
	fmt.Println("e/exit/q/quit - exit")
	fmt.Println()
	fmt.Println("For help, enter: help/h")
}

func handleUserInput(command string) {
	// Parse a command from user
	tokens := strings.Fields(command)
	fmt.Println(tokens)
	
	if len(tokens) == 0{
		return
	}
	
	if len(tokens) > 1 {
		printUsage()
		return
	}
	
	switch tokens[0] {
	case "h":
		fallthrough
	case "help":
		printUsage()
	case "e":
		fallthrough
	case "exit":
		fallthrough
	case "q":
		fallthrough
	case "quit":
		os.Exit(0)
	case "pc":
		confJson, _ := json.MarshalIndent(&conf, "", "    ")
		fmt.Println(string(confJson))
	case "pv":
		fmt.Println("Remaining tickets:", conf.RemainingTickets)
	case "pq":
		// Take the items out; they arrive in decreasing priority order.
		// TODO: find a better way to print a priority queue
		var tempList []*Request
		for waitQueue.Len() > 0 {
			item := heap.Pop(&waitQueue).(*Request)
			tempList = append(tempList, item)
			itemJson, _ := json.MarshalIndent(item, "", "    ")
			fmt.Println(string(itemJson))
		}
		for _, item:= range tempList{
			heap.Push(&PriorityQueue{}, item)
		}
	case "pt":
		clockJson, _ := json.MarshalIndent(&lamClock, "", "    ")
		fmt.Println(string(clockJson))
	
	default:
		printUsage()
	}
}

func waitUserInput() {
	for{
		if allConnected{
			break
		}
		time.Sleep(1000*time.Millisecond)
	}
	printUsage()
	fmt.Print("> ")
	
	for {
		// command line user interface
		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')
		handleUserInput(command)
		time.Sleep(80 * time.Millisecond)
		fmt.Print("> ")
	}
}