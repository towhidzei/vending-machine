package requests

// we use it for update product too
type AddProduct struct {
	Name  string `json:"name" validate:"required,alphanum"`
	Count uint   `json:"count" validate:"required,gt=0"`
	Price uint   `json:"price" validate:"required,gte=5"`
}

type Buy struct {
	// map of product id => count
	Cart map[uint]uint `json:"cart" validate:"required"`
}
