package handlecmd

import (
	"fmt"
	"log"

	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/global"
)


func (cmd *Command) SetDB() bool {
	if len(cmd.Args) < 3 || len(cmd.Args) > 6 {
		cmd.Conn.Write([]byte("-ERR wrong number of arguments for '" + cmd.Args[0] + "' command\r\n"))
		return true
	}

	log.Println("Handle SET")
	log.Println(cmd.Args)

	referenceToDB := global.GetCache()

	log.Println("Value length", len(cmd.Args[2]))
	pos := 3
	if len(cmd.Args) > 3 {
		
		option := strings.ToUpper(cmd.Args[pos])
		switch option {
		case "NX":
			log.Println("Handle NX")
			if _, ok := referenceToDB.Get(cmd.Args[1]); ok {
				if global.GetSerConf().Role == "master" {
					cmd.Conn.Write([]byte("$-1\r\n"))
				}
				return true
			}
			pos++
		case "XX":
			log.Println("Handle XX")
			if _, ok := referenceToDB.Get(cmd.Args[1]); !ok {
				if global.GetSerConf().Role == "master" {
					cmd.Conn.Write([]byte("$-1\r\n"))
				}
				return true
			}
			pos++
		}
	}

	if len(cmd.Args) > pos {
		option := strings.ToUpper(cmd.Args[pos])
		value, _ := strconv.Atoi(cmd.Args[pos+1])
		var duration time.Duration
		switch option {
		case "EX":
			duration = time.Second * time.Duration(value)
		case "PX":
			duration = time.Millisecond * time.Duration(value)
		default:
			if global.GetSerConf().Role == "master" {
				cmd.Conn.Write([]byte("-ERR expiration option is not valid" + "\r\n"))
			}
			return true
			
		}

		referenceToDB.SetWithTimeLimit(cmd.Args[1], cmd.Args[2], duration)
		if global.GetSerConf().Role == "master" {
			cmd.Conn.Write([]byte("+OK\r\n"))
		}
		cmd.sendDataToAllReplicaServer()
		return true

	}

	referenceToDB.Set(cmd.Args[1], cmd.Args[2] )
	cmd.Conn.Write([]uint8("+OK\r\n"))
	cmd.sendDataToAllReplicaServer()


	


	
	return true

}


func (cmd *Command) sendDataToAllReplicaServer() {
		// This is to forward WRITE ops to all Replica
		if global.GetSerConf().Role == "master" {
			dataTOSendToReplica := convertToRESPArray(cmd.Args)

			// dataTOSendToReplica := "*3\\r\\n$3\\r\\nset\\r\\n$1\\r\\n1\\r\\n"

			
			// dataTOSendToReplica := fmt.Sprintf("*3\r\n$3\r\nset\r\n$1\r\na\r\n$1\r\n1\r\n")
			// dataTOSendToReplica := fmt.Sprintf("*3\r\n$3\r\nSET\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", 1, "a", 1, "1")

	
			fmt.Println(dataTOSendToReplica)

			
	
		


			allReplicaConnInfoTOSend := global.GetReferenceToReplicaConnInfo()

			for _, replicaConn := range allReplicaConnInfoTOSend.SlaveHandshakeList {

				log.Println((*replicaConn).RemoteAddr(), "Sending Write Ops")
				(*replicaConn).Write([]byte(dataTOSendToReplica))
			}
		}
}

// func convertToRESPArray(input []string) string {
// 	// Remove square brackets and split the string by space
// 	elements := strings.Fields(strings.Trim(input, "[]"))

// 	fmt.Println(elements)

// 	// Create a RESP array
// 	respArray := ""
// 	respArray = respArray + "*"+fmt.Sprint(len(elements)) + "\r\n"

// 	for _, element := range elements {
// 		respArray = respArray + "$"+ fmt.Sprint(len(element)) + "\r\n" + element + "\r\n"
// 	}

// 	respArray = respArray[:len(respArray)-2]

// 	return respArray
// }


func convertToRESPArray(a []string) string {

	var b strings.Builder

	b.WriteString(fmt.Sprintf("*%d\r\n", len(a)))

	for _, s := range a{
		b.WriteString(BulkString(s))
	}

	return b.String()

}


func BulkString(s string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}


