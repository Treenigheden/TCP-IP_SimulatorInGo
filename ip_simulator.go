package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func client(send chan string, recieve chan string) {

	fmt.Println("CLIENT MADE")
	var connected = false

	var isn = generateISNForClient()

	//wait to make sure the server is ready to recieve messages before sending
	time.Sleep(time.Millisecond * time.Duration(1000))

	//first establish connection to server using 3-way handshake
	if !connected {
		time.Sleep(time.Millisecond * time.Duration(3000))
		fmt.Println("CLIENT: Trying to connect to server...")

		//1: Send synchronization request
		var message = "SYN " + strconv.Itoa(isn)
		send <- message
		isn += 1

		//2: recieve synchronization acknowledgement from server
		var serverresponse = <-recieve
		var split = strings.Split(serverresponse, " ")

		// 3: Send back acknowledgment of server response with ACK message.
		if split[0] == "SYN" && split[1] == "ACK" && split[2] == strconv.Itoa(isn) {
			serverISN, err := strconv.Atoi(split[3])
			if err == nil {
				message = "ACK " + strconv.Itoa(serverISN+1) + " " + strconv.Itoa(isn)
				send <- message
				connected = true
				fmt.Println("CLIENT: Connected to server")
			}
		}
	}
	//After connection is established, we send data!
	var sequenceNumber = isn
	var listeningForAck = true
	for connected {
		//generate sequence number string + data (random number between 0 and 100) as string and send it
		stringSequenceNumber := strconv.Itoa(sequenceNumber)
		data := strconv.Itoa(rand.Intn(100))
		//increment sequence number
		sequenceNumber++

		//create and send message
		var message = "SEQUENCE NUMBER: " + stringSequenceNumber + " DATA: " + data
		send <- message

		//wait for ACK response from server for sent message (with same sequence number) or timeout

		listeningForAck = true
		for listeningForAck {
			select {
			case ack := <-recieve:
				//we add a delay to simulate network connection delay
				time.Sleep(time.Millisecond * time.Duration(500))
				if strings.Contains(ack, "SEQUENCE NUMBER: "+stringSequenceNumber) {
					// Received the expected ACK from server
					fmt.Println("CLIENT: RECIEVED ACKNOWLEDGEMENT FROM SERVER FOR RECIEVED MESSAGE WITH SEQUENCE NUMBER: ", stringSequenceNumber, ack)
					listeningForAck = false
					break
				} else {
					//handle out of order/missing ACKs somehow
					listeningForAck = false
					break
				}
			case <-time.After(5000 * time.Millisecond):
				// Timeout occurred.
				fmt.Println("CLIENT: TIMEOUT OCCURRED WAITING FOR ACKNOWLEDGEMENT FROM SERVER FOR MESSAGE WITH SEQUENCE NUMBER", stringSequenceNumber)
				listeningForAck = false
				break
			}
		}
	}
}

func server(send chan string, recieve chan string) {
	fmt.Println("SERVER MADE")

	var isn = generateISNForServer()
	var err error

	var connected = false

	//wait for connection with client ...
	for connected == false {
		time.Sleep(time.Millisecond * time.Duration(3000))
		fmt.Println("SERVER: Waiting for client to connect...")

		//Recieve synchronization request
		var syncReq = <-recieve
		var split = strings.Split(syncReq, " ")

		// client tries to connect to server using SYNchronization
		if split[0] == "SYN" {
			var clientISN, err = strconv.Atoi(split[1])
			if err == nil {
				//acknowledge synchronization attempt:
				send <- "SYN ACK " + strconv.Itoa(clientISN+1) + " " + strconv.Itoa(isn)
				isn += 1
			}
		}

		var ackReq = <-recieve
		split = strings.Split(ackReq, " ")

		// We validate ack message from client. If everything matches, then the 3-way handshake has been a success.
		//A connection is established!
		if split[0] == "ACK" {

			if err == nil && split[1] == strconv.Itoa(isn) {
				connected = true
				fmt.Println("SERVER: CLIENT connected!")
			} else {
				fmt.Print("FAILURE TO CONNECT!")
			}
		}
	}

	//We simulate data reception and response
	//wait before beginning to send data...
	time.Sleep(time.Millisecond * time.Duration(3000))
	for connected {

		//we add a delay before sending response
		time.Sleep(time.Millisecond * time.Duration(1000))

		//recieve request
		var request = <-recieve
		fmt.Println("SERVER: Recieved request from client:", request)

		//we add a delay before sending response
		time.Sleep(time.Millisecond * time.Duration(1000))

		//Generate response and send back to client
		var response = "SERVER: Response on request: " + request
		send <- response

	}
}

// To ensure that a unique sequence number is generated for both client and server we use
// seperate functions to generate ISN for client and server.
func generateISNForClient() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(1000)
}
func generateISNForServer() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(1000)
}

func main() {
	var ch_cts = make(chan string)
	var ch_stc = make(chan string)
	go server(ch_cts, ch_stc)
	go client(ch_stc, ch_cts)
	for {
	}
}
