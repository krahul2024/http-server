package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	fmt.Println("This is some hello text")
	connect()
}

func connect() {
	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatalf("Error = %v\n", err.Error())
	}

	log.Printf ("Listening for connections on localhost%v", ":4000")

	connCounter := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error = %v\n", err.Error())
			break
		}

		connCounter += 1
		go handleConn(conn, connCounter)
	}
}

func handleConn(conn net.Conn, connCounter int) {
	defer conn.Close()

	for {
		msg, err := readMsg(conn)
		if err != nil {
			log.Printf ("Conn[%v]: Read Error = %v\n", connCounter, err.Error())
			return;
		}

		fmt.Printf ("Conn[%v]: Message = %v", connCounter, msg)

		writeMsg        := "Status: MSG_READ\nLast Message:" + msg
		bytesWrite, err := conn.Write([]byte(writeMsg))
		if err != nil {
			log.Printf("Error = %v\n", err.Error())
			return
		}

		fmt.Println (bytesWrite, "bytes Written to connection ", connCounter)
	}
}

func readMsg (conn net.Conn) (string, error) {
	readSize := 4 * 1024; // 4KB as of now
	b        := make([]byte, readSize)
	var buf bytes.Buffer

	for {
		n, err := conn.Read(b)
		if n > 0 {
			buf.Write(b[:n])
		}

		if err != nil {
			if err == io.EOF {
				break;
			}

			return buf.String(), err
		}
	}

	return buf.String(), nil
}

//TODO: add a queue to handle connection limit
