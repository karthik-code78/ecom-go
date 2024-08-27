package handlers

import "order-service/models"

func BatchAssignProductsFromCart(orderProducts []models.OrderProduct) error {
	err := OrderProductsRepo.CreateBatch(orderProducts)
	if err != nil {
		return err
	}
	return nil
}
