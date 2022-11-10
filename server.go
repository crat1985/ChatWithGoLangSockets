package main

import (
	"log"
	"net"
)

func main() {
	log.Println("Starting server...")
	server, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Print(err)
		}
		go listenMsg(conn.RemoteAddr().String(), conn)
	}
}

func listenMsg(sender string, conn net.Conn) {
	message := make([]byte, 1024)
	n, err := conn.Read(message)
	if err != nil {
		log.Print(err)
	}
	messageString := string(message[:n])
	log.Printf("New message from %s : %s\n", sender, messageString)
}
