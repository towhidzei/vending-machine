package requests

type Deposit struct {
	Coin uint `json:"coin" validate:"required,oneof=5 10 20 50 100"`
}
