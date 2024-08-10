package handlecmd

import "log"

func (cmd *Command) Ping() bool {
	if len(cmd.Args) != 1 {
		cmd.Conn.Write([]byte("-ERR wrong number of arguments for '" + cmd.Args[0] + "' command\r\n"))
		return true
	}
	log.Println("Ping received")
	cmd.Conn.Write([]byte("+PONG\r\n"))
	return true
}