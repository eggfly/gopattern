package factorymethod

import (
	"fmt"
)

type IButton interface {
	Click()
}

type IButtonFactory interface {
	CreateButton() IButton
}

type WindowsButton struct{}
type MacButton struct{}

func (this WindowsButton) Click() { fmt.Println("WindowsButton.Click") }
func (this MacButton) Click()     { fmt.Println("MacButton.Click") }

type WindowsButtonFactory struct{}
type MacButtonFactory struct{}

func (this WindowsButtonFactory) CreateButton() IButton {
	return WindowsButton{}
}

func (this MacButtonFactory) CreateButton() IButton {
	return MacButton{}
}
