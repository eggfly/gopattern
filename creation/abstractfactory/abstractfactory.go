package abstractfactory

// 对创建新的工厂的开放,面向产品族

import (
	"fmt"
)

type IDesktopProduct interface {
	PlayDesktop()
}

type ILaptopProduct interface {
	PlayLaptop()
}

type IComputerFactory interface {
	CreateDesktopProduct() IDesktopProduct
	CreateLaptopProduct() ILaptopProduct
}

type IBMFactory struct{}
type IBMDesktopComputer struct{}
type IBMLaptopComputer struct{}

func (this IBMDesktopComputer) PlayDesktop() { fmt.Println("IBMDesktopComputer.PlayDesktop") }
func (this IBMLaptopComputer) PlayLaptop()   { fmt.Println("IBMLaptopComputer.PlayLaptop") }

func (this IBMFactory) CreateDesktopProduct() IDesktopProduct {
	return IBMDesktopComputer{}
}
func (this IBMFactory) CreateLaptopProduct() ILaptopProduct {
	return IBMLaptopComputer{}
}

type AppleFactory struct{}
type MacPro struct{}
type MacBookAir struct{}

func (this MacPro) PlayDesktop()    { fmt.Println("MacPro.PlayDesktop") }
func (this MacBookAir) PlayLaptop() { fmt.Println("MacBookAir.PlayLaptop") }
func (this AppleFactory) CreateDesktopProduct() IDesktopProduct {
	return MacPro{}
}
func (this AppleFactory) CreateLaptopProduct() ILaptopProduct {
	return MacBookAir{}
}

type DellFactory struct{} // TODO
