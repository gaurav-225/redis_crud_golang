package main

import (
	"flag"
	"fmt"
	"strings"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/global"
	handleconnection "github.com/codecrafters-io/redis-starter-go/handleConnection"
)

// [URL]handle map using blog (https://reintech.io/blog/implementing-caching-in-go-using-sync-map-package)


func main() {
	
	fmt.Println("Logs from your program will appear here!")

	fmt.Println("initialize DB cache for storing value")
	global.InitCacheDB()

	global.InitServerConf()


	// https://bulutbatuhan.medium.com/how-to-parse-command-line-arguments-in-go-cf78f353d587 [ Blog I followed to implement ]
	portNum := flag.String("port", "6379", "Provide valid port number")
	replicaOf := flag.String("replicaof", "", "Provide master server address: localhost 6379" )
	flag.Parse()

	if *replicaOf != "" {
		valuesOfMasterAndPort := strings.Split(*replicaOf, " ")

		master_replid := ""
		master_repl_offset := 0

		//ownAddress for now default is localhost

		ownhost := "localhost"

		fmt.Printf("It is  Slave server with Master address: %s and Port number: %s", valuesOfMasterAndPort[0], valuesOfMasterAndPort[1] )
		// fmt.Printf("It is  Slave server will run on port %s and  with Master address: %s and Port number: %s", *portNum, valuesOfMasterAndPort[0], flag.Arg(0) )
		
		global.SetSerConfForFirstTime("slave", *portNum, valuesOfMasterAndPort[0], valuesOfMasterAndPort[1], master_replid, master_repl_offset, ownhost)
		// global.SetSerConfForFirstTime("slave", *portNum, valuesOfMasterAndPort[0], flag.Arg(0), master_replid, master_repl_offset)

		handleconnection.HandShakeWithMaster(global.GetSerConf())

		fmt.Println("Hanshake conn is here ", (*global.HandshakeConn).RemoteAddr().String())


		//maintaining handshake connection for being contact with master for Write ops
		go handleconnection.Handleconnection(*global.HandshakeConn)





	} else {

		// reference to global variable containing information of about the server

		fmt.Printf("It is  Master server and will run on port %s ", *portNum)
		master_replid := "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
		master_repl_offset := 0

		//ownAddress for now default is localhost

		ownhost := "localhost"

		global.SetSerConfForFirstTime("master", *portNum, "", "", master_replid, master_repl_offset, ownhost)



		//Will be used via tester
		global.InitReplicaConnInfo()
	}


	addressWithPort := fmt.Sprintf("0.0.0.0:%s", *portNum)
	fmt.Println("\nAddress for server starting: ", addressWithPort)

	// Start listening on port and throw error if not started
	l, err := net.Listen("tcp", addressWithPort)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}


	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleconnection.Handleconnection(conn)
	}
	
}
