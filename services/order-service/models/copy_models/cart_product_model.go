package copy_models

type CartProductModel struct {
	ID        uint
	CartID    uint
	Cart      CartModel
	ProductId uint
	Quantity  uint
	Value     float64
}
