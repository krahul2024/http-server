package main

import (
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
		connCounter -= 1
	}
}

func handleConn(conn net.Conn, connCounter int) {
	defer conn.Close()

	b := make([]byte, 2048) // reading 2KB
	for {
		bytesRead, err := conn.Read(b)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error = %v\n", err.Error())
			}

			break
		}

		fmt.Print (bytesRead, " bytes read\nMessage = ", string(b[:bytesRead]))

		writeMsg := "Status: MSG_READ\nLast Message:" + string(b[:bytesRead])
		bytesWrite, err := conn.Write([]byte(writeMsg))
		if err != nil {
			log.Printf("Error = %v\n", err.Error())
			return
		}

		fmt.Println (bytesWrite, "bytes Written to connection ", connCounter)
	}

	fmt.Println ("Closed the connection[", connCounter, "]")
}

//TODO: add a queue to handle connection limit
