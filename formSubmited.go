package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

func submited() {
	infosOk, infos := checkInfos()
	if !infosOk {
		return
	}
	var err error
	var address string
	var port string
	if serverAddressEntry.Text == "" {
		address = serverAddressEntry.PlaceHolder
	} else {
		address = serverAddressEntry.Text
	}
	if serverPortEntry.Text == "" {
		port = serverPortEntry.PlaceHolder
	} else {
		address = serverPortEntry.Text
	}
	conn, err = net.Dial(strings.ToLower(network.Selected), fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		displayErrToLoginWin(err)
		return
	}
	_, err = conn.Write([]byte(infos))
	if err != nil {
		displayErrToLoginWin(err)
		conn.Close()
		return
	}
	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		displayErrToLoginWin(err)
		conn.Close()
		return
	}
	var stringResponse string = string(response[:n])
	if stringResponse == "no" {
		displayErrToLoginWin(errors.New("informations de connexion invalides"))
		conn.Close()
		return
	}
	if stringResponse == "already connected" {
		displayErrToLoginWin(errors.New("déjà connecté"))
		conn.Close()
		return
	}
	if stringResponse != "yes" {
		return
	}
	loginWin.Close()
	displayChatWin()
}
