package global

import (
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/replicainfo"
)

var (


	replicaConnInstance *replicainfo.ReplicaInfoConnStruct
	
)

var (
	HandshakeConn *net.Conn
)




func InitReplicaConnInfo() {
	log.Println("Initialize replica  C O N N info at master side")
	replicaConnInstance = replicainfo.NewReplicaConnInfo()
	
}




func GetReferenceToReplicaConnInfo() *replicainfo.ReplicaInfoConnStruct {
	return replicaConnInstance
}




// S E N D Request to master






// func AddReplicaData(addressToAdd string) {
// 	replicaInstance.Mu.Lock()
// 	defer replicaInstance.Mu.Unlock()
// 	replicaInstance.Slaveadress = append(replicaInstance.Slaveadress, addressToAdd)
// }





