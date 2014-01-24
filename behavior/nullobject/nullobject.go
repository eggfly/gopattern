package nullobject

type IBase interface {
	Do()
}

type ConcreteImpl struct{}

func (this ConcreteImpl) Do() {
	println("ConcreteImpl.Do")
}

type NullImpl struct{}

func (this NullImpl) Do() {
	// Do nothing
}

var Null IBase = NullImpl{}
