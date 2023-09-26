# TCP-IP_SimulatorInGo

Repo for the handin 2 in BSDISYS1KU in 2023, simulating TCP/IP in Golang to understand the workings and importance of the TCP protocol

# Questions for the handin

a) What are packages in your implementation? What data structure do you use to transmit data and meta-data?

    In our implementation we have simulated a network with data exchanged between a client and server.

    In our case the data and metadata are combined into strings and sent over channels between the server and client using the send and receive channels.

    These strings contain the relevant data and metadata separated by spaces. For establishing connection that is SYN (synchronisation requests), ACK (acknowledgement) messages (this includes the server and client ISN). For the data sent after the connection has been established this means a sequence number as well as the ‘actual’ data

b) Does your implementation use threads or processes? Why is it not realistic to use threads?

    Our implementation uses goroutines which are like threads.
    Using threads does not make for a realistic simulation of a web connection since they lack the uncertainty and delay which exists in the real web. This means that it is not necessary to take into account packages/messages that get lost or arrive out of order in a simulation using threads since threads operate in a predictable way. Web connections involve dynamic factors such as network latency, delays and variable and unpredictable response times from servers.

c) In case the network changes the order in which messages are delivered, how would you handle message re-ordering?

    Because messages may be delivered out of order (or not at all) in a real web connection, it is necessary to think about how to deal with these cases.
    Message acknowledgements and sequence numbers:
    TCP divides your data into packets, and assigns them labels (for example, packet 1, packet 2, etc.). The receiver then sends back an ACK to the sender, once they have received these packets. If the packets arrive in wrong order (for example, packet 3, packet 1, packet 2), then they are reordered in the incoming message buffer, before being sent to the application layer.
    In our implementation ee use message acknowledgements and sequence numbers to keep track of which messages sent by the client have been received by the server. This gives us a way to consider reordering messages out of order, become aware of lost messages.
    In the current implementation we check that the sequence number of the received acknowledgement matches the expected sequence number. We could implement functionality to order the data if it is out of order based on the sequence numbers by temporarily storing it until the missing messages arrive. If the sequence number is too far from the expected we could consider it lost and request it be sent again.
    We have (attempted to) implement a timeout logic which will make the client aware if the server has not recieved a message, so that this can be dealt with. If an acknowledgement of a message with a specific sequence number is expected but doesn't arrive within a certain time frame, we should consider the message lost and we could request a transmission of the same message again.

d) In case messages can be delayed or lost, how does your implementation handle message loss?

    The current implementation does not really handle this case … However:
    Our implementation makes sure an ACK message is sent back from the receiver to the sender, to confirm the successful transformation of data. If this message isn't received by the sender, we would then resend the message.

e) Why is the 3-way handshake important?

    The 3-way handshake is an important part of establishing reliable and stable connection between two parties over a server or network. Between a client and a server, the 3-way handshake can ensure that both the client and server are aware of and expecting to communicate.
    An alternate approach is UDP, where no initial handshake is attempted to ensure a stable connection and a responsive and aware reciever. Unlike UDP, which is a connectionless model, the 3-way handshake mechanism of TCP makes sure that someone will be on the receiving end of sent data.
    The proces involved:
        1. The client sends SYN (synchronisation) request to the server to initiate the connection
        2. The server acknowledges the SYN request and sends back a SYN-ACK (synchronize-acknowledgment) message/packet
        3. The client sends back an ACK message/packet to acknowledge the servers response, and both parties can start exchanging data.
    By exchanging these initial packets while establishing the connection, both the client and server can confirm that the other party is reachable and responsive. If any of these packets are lost or not received, or not responded to in the expected way the connection is not established, and the sender can retry the handshake. The 3-way handshake is also used to establish the sequence numbers which help guarantee that the packets during the subsequent communication arrive in the correct order.
