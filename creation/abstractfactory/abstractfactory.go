package abstractfactory

type IDesktopProduct interface {
	PlayDesktop()
}

type ILaptopProduct interface {
	PlayLaptop()
}

type IComputerFactory interface {
	createDesktopProduct() IDesktopProduct
	createLaptopProduct() ILaptopProduct
}

type IBMFactory struct{}

func (this IBMFactory) createDesktopProduct() IDesktopProduct {
	return nil
}
func (this IBMFactory) createLaptopProduct() ILaptopProduct {
	return nil
}

type AppleFactory struct{}
type DellFactory struct{}
