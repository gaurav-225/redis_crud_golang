package handlecmd

import (
	"log"

	
)

func (cmd *Command) ReplconfReceived() bool {


	log.Println(cmd.Args)

	if cmd.Args[1] == "GETACK" && cmd.Args[2] == "*" {
		log.Println("As a Replica server sending ACK to master")
		dataToSend := "*3\r\n$8\r\nREPLCONF\r\n$3\r\nACK\r\n$1\r\n0\r\n"
		cmd.Conn.Write([]byte(dataToSend))
	}else if cmd.Args[1] == "ACK" {

		log.Println("Got ACK for Replconf GetAck as ", cmd.Args[2])
		
	}else {
		cmd.Conn.Write([]byte("+OK\r\n"))
	}
	
	// cmd.Conn.Write([]byte("+OK\r\n"))
	
	return true
}