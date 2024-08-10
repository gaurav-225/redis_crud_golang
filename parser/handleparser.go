package parser

import (
	"bufio"
	"errors"
	"log"
	"net"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/handlecmd"
)

type Parser struct {
	Connec net.Conn
	R *bufio.Reader
	Line []byte
	Pos int
}


func NewParser(conn net.Conn) *Parser {
	return &Parser{
		Connec: conn,
		R: bufio.NewReader(conn),
		Line: make([]byte, 0),
		Pos: 0,
	}
}


func (p *Parser) Advance() {
	p.Pos++
}

func (p *Parser) AtEnd() bool {
	return p.Pos >= len(p.Line)
}

func (p *Parser) Current() byte {
	if p.AtEnd() {
		return '\r'
	}

	return p.Line[p.Pos]
}


func (p *Parser) ConsumeString() (s []byte, err error) {
	for p.Current() != '"' && !p.AtEnd() {
		cur := p.Current()
		p.Advance()
		next := p.Current()
		if cur == '\\' && next == '"' {
			s = append(s, '"')
			p.Advance()
		} else {
			s = append(s, cur)
		}
	}
	if p.Current() != '"' {
		return nil, errors.New("unbalanced quotes in request")
	}
	p.Advance()
	return
}

func (p *Parser) ReadLine() ([]byte, error) {
	line, err := p.R.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	// Trim the trailing \r\n
	if len(line) >= 2 && line[len(line)-2] == '\r' && line[len(line)-1] == '\n' {
		line = line[:len(line)-2]
	}

	return line, nil
}



func (p *Parser) ConsumeArg() (s string, err error) {
	for p.Current() == ' '{
		p.Advance()
	}

	for !p.AtEnd() && p.Current() != ' ' && p.Current() != '\r' {
		s+= string(p.Current())
		p.Advance()
	}
	return
}


func (p *Parser) Command() (handlecmd.Command, error) {
	b, err := p.R.ReadByte()

	if err != nil {
		return handlecmd.Command{}, err
	}

	log.Println("Value of first byte is: ", string(b))
	if b == '*' {
		log.Println("RESP array received from client")
		return p.respArray()
		
	} else if b == '+' {
		log.Println("Received a simple String: ", b)
		return p.simpleString()
	} else {

		line, err := p.ReadLine()

		log.Println("After doing p.Readline: ", string(b)+string(line))
		if err != nil {
			return handlecmd.Command{}, err
		}
		p.Pos = 0
		p.Line = append([]byte{}, b)
		p.Line = append(p.Line, line...)
		return p.inline()

	}
}


func (p *Parser) simpleString() (handlecmd.Command, error) {
	for p.Current() == ' '{
		p.Advance()
	}

	cmd := handlecmd.Command { Conn : p.Connec}

	elementsStr, err := p.R.ReadBytes('\n')
	
	if err != nil {
		return cmd, err
	}
	log.Println(string(elementsStr))
	elementsStr = elementsStr[:len(elementsStr)-2]

	cmd.Args = append(cmd.Args, string(elementsStr))
	log.Println("Simple String received: ", cmd.Args)
	return cmd, nil



}
func (p *Parser) inline() (handlecmd.Command, error) {
	for p.Current() == ' '{
		p.Advance()
	}

	cmd := handlecmd.Command { Conn : p.Connec}
	for !p.AtEnd() {
		arg,err := p.ConsumeArg()
		if err != nil {
			return cmd, nil
		}

		if arg != "" {
			cmd.Args = append(cmd.Args, arg)
		}
	}
	log.Println("Inline Return of cmd.Args", cmd.Args)
	return cmd, nil
}

// func (p *Parser) respArray() (handlecmd.Command, error) {
// 	cmd := handlecmd.Command{}
// 	elementsStr, err := p.ReadLine()
// 	if err != nil {
// 		return cmd, err
// 	}
// 	elements, _ := strconv.Atoi(string(elementsStr))
// 	log.Println("Number of Elements: ", elements)
// 	for i := 0; i < elements; i++ {
// 		tp, err := p.R.ReadByte()
// 		if err != nil {
// 			return cmd, err
// 		}
// 		switch tp {
// 		case ':':
// 			arg, err := p.ReadLine()
// 			if err != nil {
// 				return cmd, err
// 			}
// 			cmd.Args = append(cmd.Args, string(arg))
// 		case '$':
// 			arg, err := p.ReadLine()
// 			if err != nil {
// 				return cmd, err
// 			}
// 			length, _ := strconv.Atoi(string(arg))
// 			text := make([]byte, 0)
// 			for i := 0; len(text) <= length; i++ {
// 				line, err := p.ReadLine()
// 				if err != nil {
// 					return cmd, err
// 				}
// 				text = append(text, line...)
// 			}
// 			cmd.Args = append(cmd.Args, string(text[:length]))
// 		case '*':
// 			next, err := p.respArray()
// 			if err != nil {
// 				return cmd, err
// 			}
// 			cmd.Args = append(cmd.Args, next.Args...)
// 		}
// 	}
// 	return cmd, nil
// }

func (p *Parser) respArray() (handlecmd.Command, error) {
	cmd := handlecmd.Command{Conn: p.Connec}
	
	// Read the array header
	line, err := p.ReadLine()
	if err != nil {
		return cmd, err
	}
	// if line[0] != '*' {
	// 	return cmd, errors.New("expected '*', got " + string(line[0]))
	// }
	elements, err := strconv.Atoi(string(line[0]))
	if err != nil {
		return cmd, err
	}
	log.Println("Number of Elements: ", elements)

	// Iterate over the elements in the array
	for i := 0; i < elements; i++ {
		tp, err := p.R.ReadByte()
		if err != nil {
			return cmd, err
		}
		switch tp {
		case '$': // Bulk String
			lengthStr, err := p.ReadLine()
			if err != nil {
				return cmd, err
			}
			length, err := strconv.Atoi(string(lengthStr))
			if err != nil {
				return cmd, err
			}
			text := make([]byte, length)
			_, err = p.R.Read(text)
			if err != nil {
				return cmd, err
			}
			// Read and discard the trailing CRLF
			if _, err := p.R.ReadBytes('\n'); err != nil {
				return cmd, err
			}
			cmd.Args = append(cmd.Args, string(text))
		default:
			return cmd, errors.New("unexpected type byte: " + string(tp))
		}
	}
	log.Println(cmd.Args)
	return cmd, nil
}

