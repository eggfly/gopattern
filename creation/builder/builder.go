package builder

import (
	"fmt"
	"log"
)

type IHouseBuilder interface {
	BuildWallCeilingAndFloor()
	BuildDoorAndWindows()
	BuildDecorations()
	getHouse() *House
}
type House struct {
	WallCeilingAndFloor string
	DoorAndWindows      string
	Decorations         string
}
type YoungHouseBuilder struct {
	house House
}

func (this *YoungHouseBuilder) BuildWallCeilingAndFloor() {
	obj := "vanward wall ceiling and floor"
	this.house.WallCeilingAndFloor = obj
	fmt.Println(obj)
}

func (this *YoungHouseBuilder) BuildDoorAndWindows() {
	obj := "vanward door and windows"
	this.house.DoorAndWindows = obj
	fmt.Println(obj)
}

func (this *YoungHouseBuilder) BuildDecorations() {
	obj := "decorations"
	this.house.Decorations = obj
	fmt.Println(obj)
}

func (this *YoungHouseBuilder) getHouse() *House {
	return &this.house
}

type BuildHouseDirector struct {
	houseBuilder IHouseBuilder
}

func (this *BuildHouseDirector) SetHouseBuilder(b IHouseBuilder) {
	this.houseBuilder = b
}

func (this BuildHouseDirector) BuildHouse() {
	if this.houseBuilder == nil {
		log.Fatalln("houseBuilder is not set")
	}
	this.houseBuilder.BuildWallCeilingAndFloor()
	this.houseBuilder.BuildDoorAndWindows()
	this.houseBuilder.BuildDecorations()
}

func (this BuildHouseDirector) GetHouse() *House {
	if this.houseBuilder == nil {
		log.Fatalln("houseBuilder is not set")
	}
	return this.houseBuilder.getHouse()
}
