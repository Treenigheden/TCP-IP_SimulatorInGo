package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func client(send chan string, receive chan string) {

	fmt.Println("CLIENT MADE")
	var connected = false

	var isn = 100
	var err error

	// first establish connection to server using 3-way handshake
	if !connected {
		time.Sleep(time.Millisecond * time.Duration(1000))
		fmt.Print("CLIENT: Trying to connect to server...")

		// 1. send synchronization request
		message := "SYN " + strconv.Itoa(isn)
		send <- message
		isn += 100

		// 2. receive synchronization acknowledgement from server
		var serverresponse = <-receive
		var split = strings.Split(serverresponse, " ")

		// 3. send back ACK message to server to acknowledge server response
		if split[0] == "SYN" && split[1] == "ACK" {
			isn, err = strconv.Atoi(split[2])
			if err == nil {
				isn += 100
				message := "ACK " + strconv.Itoa(isn)
				send <- message
				connected = true
				fmt.Print("CLIENT: Connected to server!")
			}
		}
	}

}

func server(send chan string, recieve chan string) {
	fmt.Println("SERVER MADE")

	var isn = 100
	var err error

	// wait for connection with client...
	var connection = false
	for connection == false {
		time.Sleep(time.Millisecond * time.Duration(1000))
		fmt.Println("SERVER: Waiting for client to connect...")

		var syncReq = <-recieve
		var split = strings.Split(syncReq, " ")

		// client tries to connect to server using SYNchronization
		if split[0] == "SYN" {
			isn, err = strconv.Atoi(split[1])
			if err == nil {
				isn += 100
				// acknowledge synchronization attempt
				send <- "SYN ACK " + strconv.Itoa(isn)
			}
		}

		syncReq = <-recieve
		split = strings.Split(syncReq, " ")

		// if we receive this, then the 3-way handshake connection has been a success
		if split[0] == "ACK" {
			isn, err = strconv.Atoi(split[1])
			if err == nil {
				isn += 100
				connection = true
				fmt.Println("SERVER: Client connected! ")
			}
		}
	}

}

func main() {
	var ch_cts = make(chan string)
	var ch_stc = make(chan string)
	go server(ch_cts, ch_stc)
	go client(ch_stc, ch_cts)
	for {
	}
}
