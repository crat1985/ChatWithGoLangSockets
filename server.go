package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type socket struct {
	socket   net.Conn
	pseudo   string
	password string
}

var sockets []socket
var loginInfos []string

func main() {
	//Login infos to modify
	loginInfos = append(loginInfos, "admin:password")
	loginInfos = append(loginInfos, "example:example")
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
		go processClient(conn)
	}
}

func verifyLoginInfos(conn net.Conn) (socket, error) {
	slice := make([]byte, 1024)
	var infos string
	for {
		n, err := conn.Read(slice)
		if err != nil {
			log.Print(err)
			return socket{}, errors.New("cannot read from client " + conn.RemoteAddr().String())
		}
		infos = string(slice[:n])
		var contains bool
		for _, info := range loginInfos {
			if infos == info {
				contains = true
			}
		}
		if !contains {
			conn.Write([]byte("no"))
			continue
		}
		break
	}
	pseudo := strings.Split(infos, ":")[0]
	password := strings.Split(infos, ":")[1]
	if !isAlreadyConnected(pseudo) {
		usersocket := socket{socket: conn, pseudo: pseudo, password: password}
		_, err := conn.Write([]byte("yes"))
		if err != nil {
			log.Print(err)
			return socket{}, errors.New("error")
		}
		sockets = append(sockets, usersocket)
		return usersocket, nil
	}

	conn.Write([]byte("already connected"))
	return socket{}, errors.New(pseudo + " already connected")
}

func isAlreadyConnected(pseudo string) bool {
	for _, socket := range sockets {
		if socket.pseudo == pseudo {
			return true
		}
	}
	return false
}

func processClient(conn net.Conn) {
	usersocket, err := verifyLoginInfos(conn)
	if err != nil {
		log.Print(err)
		return
	}
	broadcast(usersocket.pseudo + " vient de se connecter au chat !")
	go listenMsg(usersocket.pseudo, conn)
}

func broadcast(msg string) {
	if strings.Split(msg, " : ")[0] == "serv" {
		msg = strings.Join(strings.Split(msg, " : ")[1:], "")
	}
	log.Print(msg)
	for _, socket := range sockets {
		socket.socket.Write([]byte(msg))
	}
}

func listenMsg(sender string, conn net.Conn) {
	message := make([]byte, 1024)
	for {
		n, err := conn.Read(message)
		if err != nil {
			removeElement(conn)
			broadcast("serv : " + sender + " vient de se déconnecter du chat !\n")
			break
		}
		messageString := string(message[:n])
		broadcast(fmt.Sprintf("%s : %s", sender, messageString))
	}
}

func removeElement(element net.Conn) {
	var index int = -1
	for i, socket := range sockets {
		if socket.socket == element {
			index = i
		}
	}
	if index == -1 {
		return
	}
	sockets = append(sockets[:index], sockets[index+1:]...)
}
