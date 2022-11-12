package main

import (
	"errors"

	"fyne.io/fyne/v2/dialog"
)

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
