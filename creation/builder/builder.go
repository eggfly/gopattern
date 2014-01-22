package builder

import (
	"fmt"
)

type IHouseBuilder interface {
	BuildWallCeilingAndFloor()
	BuildDoorAndWindows()
	BuildDecorations()
}

type YoungHouseBuilder struct{}

func (this YoungHouseBuilder) BuildWallCeilingAndFloor() {
	fmt.Println("Build vanward wall ceiling and floor.")
}

func (this YoungHouseBuilder) BuildDoorAndWindows() {
	fmt.Println("Build vanward door and windows.")
}

func (this YoungHouseBuilder) BuildDecorations() {
	fmt.Println("Build vanward decorations.")
}
