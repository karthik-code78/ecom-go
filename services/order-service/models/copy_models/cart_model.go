package copy_models

type CartModel struct {
	ID     uint
	Name   string
	UserID uint
	Value  float64
	Items  []CartProductModel
}
