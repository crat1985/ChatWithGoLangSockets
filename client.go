package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"net"
)

var conn net.Conn
var a fyne.App
var w fyne.Window
var loginButton *widget.Button

func main() {
	a = app.New()
	w = a.NewWindow("Hello World")

	loginButton = widget.NewButton("Se connecter", loginFunction)
	w.SetContent(loginButton)
	w.CenterOnScreen()
	w.ShowAndRun()
}

func loginFunction() {
	if loginButton.Text == "..." {
		return
	}
	loginWin := a.NewWindow("Se connecter")
	loginButton.SetText("...")
	loginButton.Disable()

	loginWin.SetOnClosed(func() {
		loginButton.SetText("Se connecter")
		loginButton.Enable()
	})

	pseudoEntry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()

	form := widget.NewForm(
		widget.NewFormItem("Pseudo :", pseudoEntry),
		widget.NewFormItem("Mot de passe :", passwordEntry),
	)

	loginWin.SetContent(form)
	form.Resize(fyne.NewSize(250, form.Size().Height))
	loginWin.Resize(form.Size())
	loginWin.SetFixedSize(true)
	loginWin.CenterOnScreen()
	loginWin.Show()
}
