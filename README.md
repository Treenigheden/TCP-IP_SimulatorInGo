# TCP-IP_SimulatorInGo
Repo for the handin 2 in BSDISYS1KU in 2023, simulating TCP/IP in Golang to understand the workings and importance of the TCP protocol

# Questions for the handin
a) What are packages in your implementation? What data structure do you use to transmit data and meta-data?

b) Does your implementation use threads or processes? Why is it not realistic to use threads?

c) In case the network changes the order in which messages are delivered, how would you handle message re-ordering?

TCP divides your data into packets, and assigns them labels (for example, packet 1, packet 2, etc.). The receiver then sends back an ACK to the sender, once they have received these packets. If the packets arrive in wrong order (for example, packet 3, packet 1, packet 2), then they are reordered in the incoming message buffer, before being sent to the application layer.  

d) In case messages can be delayed or lost, how does your implementation handle message loss?

Our implementation makes sure an ACK message is sent back from the receiver to the sender, to confirm the successful transformation of data. If this message isn't received by the sender, we would then resend the message. 

e) Why is the 3-way handshake important?

The 3-way handshake makes sure that a successful connection has been made between a client and a server. Unlike UDP, which is a connectionless model, the 3-way handshake mechanism of TCP makes sure that someone will be on the receiving end of sent data. 