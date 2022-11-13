package main

import (
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
var network *widget.RadioGroup
var conversationsContainer *container.DocTabs
var messagesBoxTabItem *container.TabItem

func main() {
	a = app.New()
	w = a.NewWindow("Chat en temps réel")

	loginButton = widget.NewButton("Se connecter", loginFunction)
	w.SetContent(loginButton)
	w.CenterOnScreen()
	w.ShowAndRun()
}

func displayErrToLoginWin(err error) {
	dialog.NewError(err, loginWin).Show()
}

func loginWinClosed() {
	loginButton.SetText("Se connecter")
	loginButton.Enable()
}

func createLoginForm() *widget.Form {
	network = widget.NewRadioGroup([]string{"TCP", "UDP"}, nil)
	network.SetSelected("TCP")
	serverAddressEntry = widget.NewEntry()
	serverAddressEntry.SetPlaceHolder("localhost")
	serverPortEntry = widget.NewEntry()
	serverPortEntry.SetPlaceHolder("8888")
	pseudoEntry = widget.NewEntry()
	pseudoEntry.SetPlaceHolder("Pseudo")
	passwordEntry = widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Mot de passe")
	form := widget.NewForm(
		widget.NewFormItem("Réseau :", network),
		widget.NewFormItem("Adresse du serveur :", serverAddressEntry),
		widget.NewFormItem("Port du serveur :", serverPortEntry),
		widget.NewFormItem("Pseudo :", pseudoEntry),
		widget.NewFormItem("Mot de passe :", passwordEntry),
	)
	form.OnCancel = loginWin.Close
	form.OnSubmit = submited
	return form
}

func createLoginWin() {
	loginWin = a.NewWindow("Se connecter")
	loginWin.SetOnClosed(loginWinClosed)

	form := createLoginForm()

	loginWin.SetContent(form)
	form.Resize(fyne.NewSize(400, form.Size().Height))
	loginWin.Resize(form.Size())
	loginWin.SetFixedSize(true)
	loginWin.CenterOnScreen()
	loginWin.Show()
}

func loginFunction() {
	if loginButton.Text == "..." {
		return
	}

	loginButton.SetText("...")
	loginButton.Disable()

	createLoginWin()
}

func createDocTabs() *container.DocTabs {
	messagesBoxScroll = createMessageBoxScroll()
	messagesBoxTabItem = container.NewTabItem("Général", messagesBoxScroll)
	convContainer := container.NewDocTabs(
		messagesBoxTabItem,
	)
	return convContainer
}

func displayGeneralConv() {
	for _, value := range conversationsContainer.Items {
		if value.Text == "Général" {
			return
		}
	}
	conversationsContainer.Append(messagesBoxTabItem)
}

func createLeftPanel() *fyne.Container {
	leftPanel := container.NewVBox(
		widget.NewLabel("Conversations"),
		widget.NewButton("Général", displayGeneralConv),
	)
	return leftPanel
}

func displayChatWin() {
	chatWin = a.NewWindow("Chat en temps réel")
	//Bottom
	sendMessageContainer := createSendMessageContainer()
	//Top
	conversationsContainer = createDocTabs()
	all := container.NewBorder(nil, sendMessageContainer, createLeftPanel(), nil, conversationsContainer)
	w.Close()
	chatWin.SetContent(all)
	chatWin.Resize(fyne.NewSize(720, 480))
	chatWin.Show()
	listenForMessages()
}

func sendMessage() {
	_, err := conn.Write([]byte(sendMessageEntry.Text))
	if err != nil {
		displayErrToLoginWin(err)
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
			dialog.NewError(err, chatWin).Show()
			continue
		}
		msgString := string(slice[:n])
		msgEntry := widget.NewMultiLineEntry()
		msgEntry.Disable()
		msgEntry.SetText(msgString)
		msgEntry.Wrapping = fyne.TextWrapWord
		messagesBox.Add(msgEntry)
		log.Print(msgString)
	}
}
