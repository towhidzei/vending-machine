package domain

const (
	Five    Coin = 5
	Ten     Coin = 10
	Twenty  Coin = 20
	Fifty   Coin = 50
	Hundred Coin = 100
)

type Coin uint

var Coins []uint = []uint{5, 10, 20, 50, 100}

func (c Coin) IsValid() bool {
	switch c {
	case 5:
	case 10:
	case 20:
	case 50:
	case 100:
	default:
		return false
	}
	return true
}
