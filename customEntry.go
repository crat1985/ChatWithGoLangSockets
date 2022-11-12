package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

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
