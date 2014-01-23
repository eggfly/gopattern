package prototype

// deepcopy is here: https://code.google.com/p/rog-go/source/browse/exp/deepcopy/deepcopy.go
type Bullet struct {
	X, Y, Speed int
}

func (this *Bullet) Init() {
	this.X = 1
	this.Y = 1
	this.Speed = 100
}

// prototype clone function
func (this *Bullet) CopyFrom(bullet Bullet) {
	*this = bullet
}

func (this *Bullet) Clone() Bullet {
	return *this
}
