package handleconnection

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/global"

	"github.com/codecrafters-io/redis-starter-go/serverstatus"
)

func HandShakeWithMaster(s *serverstatus.ServerConf) bool {

	masterAddres := s.MasterHost
	masterPort := s.MasterPort

	addressToConnect:= fmt.Sprintf("%s:%s", masterAddres, masterPort)
	fmt.Println("\nChecking connectivity to master ", addressToConnect)

	conn, err := net.Dial("tcp", addressToConnect)

	global.HandshakeConn = &conn

	if err != nil {
		log.Println("Failed to connect master server ", addressToConnect)
		os.Exit(1)
	}

	// -----------------------------------Tester is expecting to continue with Handshake Conn for Write ops forwarding

	// -----------------------------------So, commenting out below func

	// defer func() {
	// 	// time.Sleep(10* time.Second)
	// 	log.Println("Handshake completed with master ", addressToConnect)
	// 	log.Println("Closing connection of handshaking with master ", addressToConnect)
	// 	conn.Close()
	// }()

	// -----------------------------------------------------------------------------------------------------------------

	fmt.Println("Successfully connected to master", conn.RemoteAddr().String())
	conn.Write([]byte("*1\r\n$4\r\nping\r\n"))

	// Initial buffer size 
	buffer := make([]byte, 1024) 
	var response []byte
	sizeOfBuffereReceived, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return true
	}

	// Expecting to receive PONG
	response = append(response, buffer[:sizeOfBuffereReceived]...)
	if bytes.Contains(response, []byte("+PONG\r\n")) {
		log.Println("Received: +PONG response from ", conn.RemoteAddr().String())
	}
	
	serverConfValueOfInstance := global.GetSerConf()
	dataExtracted := serverConfValueOfInstance.Port

	dataToSend := fmt.Sprintf("*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$%d\r\n%s\r\n", len(dataExtracted), dataExtracted)
	conn.Write([]byte(dataToSend))

	sizeOfBuffereReceived, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return true
	}

	// Expecting to receive +OK
	response = append(response, buffer[:sizeOfBuffereReceived]...)
	if bytes.Contains(response, []byte("+OK\r\n")) {
		log.Println("Received: +OK response from ", conn.RemoteAddr().String())
	}

	log.Println("Sending REPLCONF capa psync2")
	conn.Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n"))

	sizeOfBuffereReceived, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return true
	}

	// Expecting to receive second +OK
	response = append(response, buffer[:sizeOfBuffereReceived]...)
	if bytes.Contains(response, []byte("+OK\r\n")) {
		log.Println("Received: Second +OK response from ", conn.RemoteAddr().String())
	}

	//

	log.Println("Sending to master: PSYNC ? -1")
	conn.Write([]byte("*3\r\n$5\r\nPSYNC\r\n$1\r\n?\r\n$2\r\n-1\r\n"))

	// sizeOfBuffereReceived, err = conn.Read(buffer)

	// Will use same reader below for RDB file
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')

	// Expected value is +FULLRESYNC <REPL_ID> 0\r\n
	// message = +FULLRESYNC <REPL_ID> 0
	message = message[:len(message)-2]

	log.Println(message)

	if err != nil {
		fmt.Println("Error reading:", err)
		return true
	}

	// Expecting to ReplId and Offset from master
	
	// response = append(response, buffer[:sizeOfBuffereReceived]...)


	//--------------------------------------------Start: Handle +FulResync from master and get MasterReplId and OffSet value
	response = append(response, []byte(message)...)
	if bytes.Contains(response, []byte("+FULLRESYNC")) {
		log.Println("Received: +FULLRESYNC response from ", conn.RemoteAddr().String())
		// dataInString := string(buffer[:sizeOfBuffereReceived])
		dataInString := message

		serverConfValueOfInstance.MasterReplid =  strings.Split(dataInString, " ")[1]

		
		offsetValue, err := strconv.Atoi(strings.ReplaceAll(strings.Split(dataInString, " ")[2], "\r\n", "")) 

		log.Println("OffSetValue is ", offsetValue)

		// storing offvalue in global variable
		serverConfValueOfInstance.MasterReplOffset = offsetValue

		if err != nil {
			log.Panicln(err)
		}



		
		
	}

	//--------------------------------------------End: Handle +FulResync from master and get MasterReplId and OffSet value

	

	// *******************************************Now Replica is Expecting to receive RDB file**********************
	
	//+++++++++++Comment since already declared at above
	// reader = bufio.NewReader(conn)

	// #####################RDB file should be sent like this: $<length>\r\n<contents>

	// ---------------------Trying to get length of rdb file and read with same buffer size from conn
	message, _ = reader.ReadString('\n')

	
	message = message[1:len(message)-2]
	log.Println("Length of RDB file going to be received is ", message)
	rdbFileSize, errWhileGetrdbsize := strconv.Atoi(message)

	if errWhileGetrdbsize != nil {
		log.Panicln(errWhileGetrdbsize)
	}



	emptyRDBFileRead := make([]byte, rdbFileSize)

	sizeOFFile, errOnRDB := reader.Read(emptyRDBFileRead)
	// sizeOFFile, errOnRDB := conn.Read(emptyRDBFileRead)

	if errOnRDB != nil {
		log.Println(errOnRDB)
	}

	log.Println(string(emptyRDBFileRead[:sizeOFFile]))

	/*

	response = append(response, emptyRDBFileRead[:sizeOFFile]...)
	if bytes.Contains(response, []byte("GETACK")) {
		log.Println("Received: Master sent GETACK", conn.RemoteAddr().String())
	}
	
	*/
	
	// -----------------------------------------------------------------


	// Give Response back for GETACK from Master

	// sizeOfBuffereReceived, err = conn.Read(buffer)
	// if err != nil {
	// 	fmt.Println("Error reading:", err)
	// 	return true
	// }

	
	// response = append(response, buffer[:sizeOfBuffereReceived]...)
	// if bytes.Contains(response, []byte("getack")) {
	// 	log.Println("Received: Master sent GETACK", conn.RemoteAddr().String())
	// }

	/*
	dataToSendACK := "*3\r\n$8\r\nREPLCONF\r\n$3\r\nACK\r\n$1\r\n0\r\n"
	conn.Write([]byte(dataToSendACK))

	*/
	// Printing server details
	global.PrintServerConf()

	// It is presented at server.go file
	// go Handleconnection(conn)
	return true


}

