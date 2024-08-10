package handlecmd

import (
	"fmt"
	"strconv"
	"strings"
)

func (cmd *Command) Echo() bool {
	if len(cmd.Args) == 1 {
		cmd.Conn.Write([]byte("-ERR wrong number of arguments for '" + cmd.Args[0] + "' command\r\n"))
		return true
	}

	dataToSend := fmt.Sprint(cmd.Args[1:])
	dataToSend = strings.Trim(dataToSend, "[]")
	lengthOfString := strconv.Itoa(len(dataToSend))
	dataToSend = "$" + lengthOfString + "\r\n" + dataToSend + "\r\n"
	cmd.Conn.Write([]byte(dataToSend))
	return true



}