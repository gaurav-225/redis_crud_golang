package handleconnection

import (
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/parser"
)

func Handleconnection(conn net.Conn) {
	
	defer func ()  {
		log.Println("Closing Connection for ", conn.RemoteAddr())
		conn.Close()
	}()
	
	p := parser.NewParser(conn)

	for {
		cmd, err := p.Command()
		
		if err != nil {
			log.Println("Error", err)
			conn.Write([]byte("-ERR " + err.Error() + "\r\n"))
			break
		}

		if !cmd.Handle() {
			break
		}
	}

	
}