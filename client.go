package main

import (
	"errors"
	"fmt"
	"log"
	"net"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type loginInfos struct {
	pseudo   string
	password string
}

type customSendMessageEntry struct {
	widget.Entry
	OnTypedKey func(key *fyne.KeyEvent)
}

func NewCustomSendMessageEntry() *customSendMessageEntry {
	e := &customSendMessageEntry{}
	e.Wrapping = fyne.TextTruncate
	e.ExtendBaseWidget(e)
	return e
}

func (m *customSendMessageEntry) TypedKey(key *fyne.KeyEvent) {
	if key.Name == "Return" {
		sendMessage()
	}
}

var conn net.Conn
var a fyne.App
var w fyne.Window
var loginButton *widget.Button
var loginInfosVar loginInfos
var serverAddressEntry *widget.Entry
var serverPortEntry *widget.Entry
var pseudoEntry *widget.Entry
var passwordEntry *widget.Entry
var loginWin fyne.Window
var messagesBox *fyne.Container
var sendButton *widget.Button
var sendMessageEntry *customSendMessageEntry
var messagesBoxScroll *container.Scroll
var chatWin fyne.Window

func main() {
	a = app.New()
	w = a.NewWindow("Chat en temps réel")

	loginButton = widget.NewButton("Se connecter", loginFunction)
	w.SetContent(loginButton)
	w.CenterOnScreen()
	w.ShowAndRun()
}

func loginFunction() {
	if loginButton.Text == "..." {
		return
	}
	loginWin = a.NewWindow("Se connecter")
	loginButton.SetText("...")
	loginButton.Disable()

	loginWin.SetOnClosed(func() {
		loginButton.SetText("Se connecter")
		loginButton.Enable()
	})

	serverAddressEntry = widget.NewEntry()
	serverAddressEntry.SetPlaceHolder("localhost")
	serverPortEntry = widget.NewEntry()
	serverPortEntry.SetPlaceHolder("8888")
	pseudoEntry = widget.NewEntry()
	pseudoEntry.SetPlaceHolder("Pseudo")
	passwordEntry = widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Mot de passe")

	form := widget.NewForm(
		widget.NewFormItem("Adresse du serveur :", serverAddressEntry),
		widget.NewFormItem("Port du serveur :", serverPortEntry),
		widget.NewFormItem("Pseudo :", pseudoEntry),
		widget.NewFormItem("Mot de passe :", passwordEntry),
	)
	form.OnCancel = loginWin.Close
	form.OnSubmit = submited

	loginWin.SetContent(form)
	form.Resize(fyne.NewSize(400, form.Size().Height))
	loginWin.Resize(form.Size())
	loginWin.SetFixedSize(true)
	loginWin.CenterOnScreen()
	loginWin.Show()
}

func checkInfos() (bool, string) {
	if pseudoEntry.Text == "" {
		dialog.NewError(errors.New("veuillez entrer un pseudo"), loginWin).Show()
		return false, ""
	}
	if passwordEntry.Text == "" {
		dialog.NewError(errors.New("veuillez entrer un mot de passe"), loginWin).Show()
		return false, ""
	}
	loginInfosVar = loginInfos{pseudo: pseudoEntry.Text, password: passwordEntry.Text}
	return true, loginInfosVar.pseudo + ":" + loginInfosVar.password
}

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
	}
	if serverPortEntry.Text == "" {
		port = serverPortEntry.PlaceHolder
	}
	conn, err = net.Dial("tcp", fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		dialog.NewError(err, loginWin).Show()
		return
	}
	_, err = conn.Write([]byte(infos))
	if err != nil {
		dialog.NewError(err, loginWin).Show()
		return
	}
	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		dialog.NewError(err, loginWin).Show()
		return
	}
	var stringResponse string = string(response[:n])
	if stringResponse == "no" {
		dialog.NewError(errors.New("informations de connexion invalides"), loginWin).Show()
		return
	}
	if stringResponse != "yes" {
		return
	}
	loginWin.Close()
	modifyWindowToBeAbleToSendMessages()
}

func modifyWindowToBeAbleToSendMessages() {
	chatWin = a.NewWindow("Chat")
	//Bottom
	sendMessageEntry = NewCustomSendMessageEntry()
	sendMessageEntry.SetPlaceHolder("Votre message ici")
	sendButton = widget.NewButton("Envoyer", sendMessage)
	sendMessageContainer := container.NewBorder(nil, nil, nil, sendButton, sendMessageEntry)
	//Top
	messageEntry := widget.NewEntry()
	messageEntry.SetText("test\ntest\nmdr")
	messageEntry.Disable()
	messageEntry.Wrapping = fyne.TextTruncate
	messageEntry.ExtendBaseWidget(messageEntry)
	messagesBox = container.NewVBox(messageEntry)
	messagesBoxScroll = container.NewScroll(messagesBox)
	all := container.NewBorder(nil, sendMessageContainer, nil, nil, messagesBoxScroll)
	w.Close()
	chatWin.SetContent(all)
	chatWin.Resize(fyne.NewSize(720, 480))
	chatWin.Show()
	listenForMessages()
}

func sendMessage() {
	_, err := conn.Write([]byte(sendMessageEntry.Text))
	if err != nil {
		dialog.NewError(err, loginWin).Show()
		return
	}
	sendMessageEntry.SetText("")
	chatWin.Canvas().Focus(sendMessageEntry)
	messagesBoxScroll.ScrollToBottom()
}

func listenForMessages() {
	for {
		slice := make([]byte, 1024)
		n, err := conn.Read(slice)
		if err != nil {
			dialog.NewError(err, loginWin).Show()
			continue
		}
		msgString := string(slice[:n])
		msgEntry := widget.NewEntry()
		msgEntry.Disable()
		msgEntry.SetText(msgString)
		msgEntryMax := container.NewMax(msgEntry)
		messagesBox.Add(msgEntryMax)
		log.Print(msgString)
	}
}
