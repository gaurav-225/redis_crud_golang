package replicainfo

import (
	"net"
	"sync"
)



//added  mutex to avoid race condition



type ReplicaInfoConnStruct struct {
	SlaveHandshakeList []*net.Conn
	MuForConn sync.Mutex
}

func NewReplicaConnInfo() *ReplicaInfoConnStruct {
	return &ReplicaInfoConnStruct{}
}




