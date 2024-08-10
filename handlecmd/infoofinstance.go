package handlecmd

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/global"
)




func (cmd *Command) InfoOfInstance() bool {
	
	serverConfValueOfInstance := global.GetSerConf()

	dataExtracted := fmt.Sprintf("role:%s\nmaster_replid:%s\nmaster_repl_offset:%d", serverConfValueOfInstance.Role, serverConfValueOfInstance.MasterReplid, serverConfValueOfInstance.MasterReplOffset)


	dataToSend := fmt.Sprintf("$%d\r\n%s\r\n", len(dataExtracted), dataExtracted)
	cmd.Conn.Write([]byte(dataToSend))
	return true
}