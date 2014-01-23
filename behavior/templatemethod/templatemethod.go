package templatemethod

import (
	"fmt"
)

type IGame interface {
	initializeGame()
	// makePlay(int)
	// endOfGame()
	// printWinner()
	// playOneGame()
}

type AbstractGame struct {
	IGame
	playersCount int
}

func (this AbstractGame) PlayOneGame(playersCount int) {
	this.playersCount = playersCount
	this.IGame.initializeGame()
}

func (this AbstractGame) initializeGame() {
	fmt.Println("AbstractGame.initializeGame")
}

type Monopoly struct {
	AbstractGame
}

func (this Monopoly) initializeGame() {
	fmt.Println("Monopoly.initializeGame")
}

type Chess struct {
	AbstractGame
}

//  below is not used
type privateClass struct{}

func CreatePrivate() privateClass {
	return privateClass{}
}
