package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func getResponse(recieve chan string, seq int) []string {
	var serverresponse = <-recieve
	var split = strings.Split(serverresponse, " ")
	var recievedseq, err = strconv.Atoi(split[0])
	if err != nil {
		fmt.Print("OH NO! The sequence was not recieved correctly")
	}
	if seq != (recievedseq - 1) {
		fmt.Print("Oh dear! The recieved sequence is NOT the same!")
	}
	return split
}

func client(send chan string, recieve chan string) {
	fmt.Println("CLIENT MADE")

	var seq = rand.Int()
	var ack int
	var err error

	//send sequence
	send <- strconv.Itoa(seq)

	//recieve acknowledgement and send acknowledgement + data
	var response = getResponse(recieve, seq)
	seq, err = strconv.Atoi(response[0])
	if err != nil {
		fmt.Print("Oh dear! Something went wrong when recieving seq from the server")
	}
	ack, err = strconv.Atoi(response[1])
	if err != nil {
		fmt.Print("Oh dear! Something went wrong when recieving acknowledgement from the server")
	}
	ack = ack + 1
	//data send
	for {
		var data string
		data = "someting"
		send <- strconv.Itoa(ack) + " " + strconv.Itoa(seq) + " " + data

		response = getResponse(recieve, seq)
		seq, err = strconv.Atoi(response[0])
		if err != nil {
			fmt.Print("Oh dear! Something went wrong when recieving seq from the server")
		}
		ack, err = strconv.Atoi(response[1])
		ack++
		fmt.Println(ack)
	}
}

func server(send chan string, recieve chan string) {
	fmt.Println("SERVER MADE")
	//recieve sequence and send acknowledgement
	var ack int
	var err error
	var seq = rand.Int()
	ack, err = strconv.Atoi(<-recieve)
	if err != nil {
		fmt.Print("OH NO! The sequence was not recieved correctly by the server")
	}
	ack++
	send <- strconv.Itoa(ack) + " " + strconv.Itoa(seq)

	for {
		time.Sleep(time.Millisecond * time.Duration(1000))
		response := getResponse(recieve, seq)
		seq, err = strconv.Atoi(response[0])
		if err != nil {
			fmt.Print("Oh dear! Something went wrong when recieving seq from the client")
		}
		ack, err = strconv.Atoi(response[1])
		if err != nil {
			fmt.Print("Oh dear! Something went wrong when recieving ack from the client")
		}
		ack++
		fmt.Println(response[1], response[2])
		send <- strconv.Itoa(ack) + " " + strconv.Itoa(seq)
	}
}

func main() {
	var ch_cts = make(chan string)
	var ch_stc = make(chan string)
	go server(ch_cts, ch_stc)
	go client(ch_stc, ch_cts)
	select {}
}
