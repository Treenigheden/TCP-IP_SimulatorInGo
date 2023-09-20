package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func client(send chan string, recieve chan string) {
	fmt.Println("SERVER MADE")
	//send sequence
	var seq = 1234
	send <- strconv.Itoa(seq)

	//recieve acknowledgement and send acknowledgement + data

	var serverresponse = <-recieve
	var split = strings.Split(serverresponse, " ")

	var recievedseq int
	var err error
	recievedseq, err = strconv.Atoi(split[0])
	if seq != (recievedseq - 1) {
		fmt.Print("Oh dear! The recieved sequence is NOT the same!")
		return
	}

	var ack int
	ack, err = strconv.Atoi(split[1])
	if err != nil {
		fmt.Print("Oh dear! Something went wrong when recieving acknowledgement from the server")
	}
	ack = ack + 1
	//data send
	for {
		var data string
		data = "someting"
		send <- strconv.Itoa(ack) + " " + strconv.Itoa(seq) + " " + data

		serverresponse = <-recieve
		split = strings.Split(serverresponse, " ")
		recievedseq, err = strconv.Atoi(split[0])
		if seq != (recievedseq - 1) {
			fmt.Print("Oh dear! The recieved sequence is NOT the same!")
			return
		}
		ack, err = strconv.Atoi(split[1])
		if err != nil {
			fmt.Print("Oh dear! Something went wrong when recieving acknowledgement from the server")
		}
		ack = ack + 1
		fmt.Println(ack)
	}
}

func server(send chan string, recieve chan string) {
	fmt.Println("SERVER MADE")
	//recieve sequence and send acknowledgement
	var ack int
	var err error
	var seq = 4321
	ack, err = strconv.Atoi(<-recieve)
	if err != nil {
		fmt.Print("OH NO! The sequence was not recieved correctly by the server")
	}
	ack = ack + 1
	send <- strconv.Itoa(ack) + " " + strconv.Itoa(seq)

	for {
		time.Sleep(time.Millisecond * time.Duration(1000))
		var clientresponse = <-recieve
		var split = strings.Split(clientresponse, " ")
		var recievedseq int
		recievedseq, err = strconv.Atoi(split[0])
		if seq != (recievedseq - 1) {
			fmt.Print("Oh dear! The recieved sequence is NOT the same!")
			return
		}
		ack, err = strconv.Atoi(split[1])
		ack++
		fmt.Println(split[1], split[2])
		send <- strconv.Itoa(ack) + " " + strconv.Itoa(seq)
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
