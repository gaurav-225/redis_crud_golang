package handlecmd

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/codecrafters-io/redis-starter-go/global"
)

func (cmd *Command) PsyncReceived() bool {
	if len(cmd.Args) == 3 && cmd.Args[1] == "?" && cmd.Args[2] == "-1" {
		log.Println(cmd.Args)
		replId := global.GetSerConf().MasterReplid
		
		dataToSend := fmt.Sprintf("+FULLRESYNC %s 0\r\n", replId)
		log.Println("Data to sent", dataToSend)
		cmd.Conn.Write([]byte(dataToSend))

		//send RDB file


		emptyRDB, err := hex.DecodeString("524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2")
		if err != nil {
			fmt.Println(err)
		}

		log.Println("Sending RDB file")
		cmd.Conn.Write(append([]byte(fmt.Sprintf("$%d\r\n", len(emptyRDB))), emptyRDB...))

		/*
		log.Println("Sending GETACK")
		cmd.Conn.Write([]byte("3\r\n$8\r\nreplconf\r\n$6\r\ngetack\r\n$1\r\n*\r\n"))
		*/

		accessToReferenceOfStringSlaveConnList := global.GetReferenceToReplicaConnInfo()

		accessToReferenceOfStringSlaveConnList.SlaveHandshakeList = append(accessToReferenceOfStringSlaveConnList.SlaveHandshakeList, &cmd.Conn)

	}

	return true
}