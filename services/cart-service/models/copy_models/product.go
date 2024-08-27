package copy_models

type Product struct {
	ID          uint    `json:"id"`
	Name        string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    uint    `json:"quantity"`
}
