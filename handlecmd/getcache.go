package handlecmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/global"
)

func (cmd *Command) GetDB() bool {

	if len(cmd.Args) != 2 {
		cmd.Conn.Write([]byte("-ERR wrong number of arguments for '" + cmd.Args[0] + "' command\r\n"))
		return true
	}

	referenceToDB := global.GetCache()

	val, _ := referenceToDB.Get(cmd.Args[1])

	if val != nil {
		res, _ := val.(string)
		if strings.HasPrefix(res, "\"") {
			res, _ = strconv.Unquote(res)
		}
		log.Println("Response length", len(res))
		cmd.Conn.Write([]byte(fmt.Sprintf("$%d\r\n", len(res))))
		cmd.Conn.Write(append([]byte(res), []byte("\r\n")...))
	} else {
		cmd.Conn.Write([]uint8("$-1\r\n"))
	}
	return true


	


}


