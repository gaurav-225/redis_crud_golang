package handlecmd

import (
	"log"
	"net"
	"strings"
)


type Command struct {
	Conn net.Conn
	Args []string
}

func (cmd Command) Handle() bool {
	switch strings.ToUpper(cmd.Args[0]) {
	case "PING":
		return cmd.Ping()
	case "ECHO":
		return cmd.Echo()
	case "SET":
		return cmd.SetDB()
	case "GET":
		return cmd.GetDB()
	case "INFO":
		return cmd.InfoOfInstance()
	case "OK":
		return cmd.Okcmd()
	case "REPLCONF":
		return cmd.ReplconfReceived()
	case "PSYNC":
		return cmd.PsyncReceived()
	default:
		log.Println("Command not supported", cmd.Args[0])
		cmd.Conn.Write([]byte("-ERR unknown command '" + cmd.Args[0] + "'\r\n"))
	}
	return true
}

