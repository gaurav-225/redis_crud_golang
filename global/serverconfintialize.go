package global

import (
	"fmt"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/serverstatus"
)

var (
	serverConfVar *serverstatus.ServerConf
	onceForConfVar sync.Once
)


func InitServerConf() {
	onceForConfVar.Do(func() {
		serverConfVar = serverstatus.NewServerConf()
	})
}

func GetSerConf() *serverstatus.ServerConf {
	return serverConfVar
}

func SetSerConfForFirstTime(role string, port string, masterAddr string, masterPort string, master_replid string, master_repl_offset int, ownhost string) {
	serverConfVar.Role = role
	serverConfVar.Port = port
	serverConfVar.MasterHost = masterAddr
	serverConfVar.MasterPort = masterPort
	serverConfVar.MasterReplid = master_replid
	serverConfVar.MasterReplOffset = master_repl_offset
	serverConfVar.OwnHost = ownhost
}

func PrintServerConf() {
	fmt.Println("************************Server Conf*******************************************")
	fmt.Println("serverConfVar.Role: ", serverConfVar.Role)
	fmt.Println("serverConfVar.Port: ", serverConfVar.Port)
	fmt.Println("serverConfVar.MasterHost: ", serverConfVar.MasterHost)
	fmt.Println("serverConfVar.MasterPort: ", serverConfVar.MasterPort)
	fmt.Println("serverConfVar.MasterReplid: ", serverConfVar.MasterReplid)
	fmt.Println("serverConfVar.MasterReplOffset: ", serverConfVar.MasterReplOffset)
	fmt.Println("serverConfVar.OwnHost: ", serverConfVar.OwnHost)
	fmt.Println("******************************************************************************")
}
