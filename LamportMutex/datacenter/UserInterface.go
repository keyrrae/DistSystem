package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func printUsage() {
	fmt.Println()
	fmt.Println("pc/config - print configuration")
	fmt.Println("pq/queue - print wait queue")
	fmt.Println("pv/value - print amount of tickets")
	fmt.Println("e/exit/q/quit - exit")
	fmt.Println("rst/reset - reset tickets and clock")
	fmt.Println()
	fmt.Println("For help, enter: help/h")
}

func handleUserInput(command string) {
	// Parse a command from user
	tokens := strings.Fields(command)

	if len(tokens) == 0 {
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

	case "config":
		fallthrough
	case "pc":
		confJson, _ := json.MarshalIndent(&conf, "", "    ")
		fmt.Println(string(confJson))

	case "value":
		fallthrough
	case "pv":
		fmt.Println("Remaining tickets:", conf.RemainingTickets)

	case "queue":
		fallthrough
	case "pq":
		// Take the items out; they arrive in decreasing priority order.
		// TODO: find a better way to print a priority queue
		for _, item := range waitQueue {
			itemJson, _ := json.MarshalIndent(item, "", "    ")
			fmt.Println(string(itemJson))
		}

	case "time":
		fallthrough
	case "pt":
		clockJson, _ := json.MarshalIndent(&lamClock, "", "    ")
		fmt.Println(string(clockJson))

	case "reset":
		fallthrough
	case "rst":
		lamClock.LogicalClock = 1
		conf.RemainingTickets = conf.InitialTktNum

	default:
		printUsage()
	}
}


func waitUserInput() {
	for {
		if allConnected {
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}
	printUsage()

	for {
		// command line user interface
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')
		handleUserInput(command)
		time.Sleep(80 * time.Millisecond)
	}
}
