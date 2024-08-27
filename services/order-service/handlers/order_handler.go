package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/karthik-code78/ecom/shared/logging"
	"github.com/karthik-code78/ecom/shared/utils/http_utils"
	"github.com/karthik-code78/ecom/shared/utils/json_utils"
	"net/http"
	"order-service/models"
	"order-service/models/copy_models"
)

func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	var orders []models.Order
	orders, err := OrdersRepo.FindAll()
	if err != nil {
		http_utils.SendErrorResponse(w, "Error while getting orders", http.StatusInternalServerError)
		return
	}
	err = json_utils.JsonEncode(w, orders)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error while encoding orders", http.StatusInternalServerError)
	}
}

func checkIfCartExists(cartId int) (bool, error, copy_models.CartModel) {
	url := fmt.Sprintf("http://localhost:8083/cart/%d", cartId)
	var cart copy_models.CartModel
	resp, err := http.Get(url)
	if err != nil {
		logging.Log.Error("error while getting resp", err)
		return false, err, cart
	}
	defer resp.Body.Close()
	logging.Log.Info(cart)
	err = json.NewDecoder(resp.Body).Decode(&cart)
	if err != nil {
		return false, err, cart
	}
	return true, nil, cart
}

func getProductsByIds(ids []uint) ([]copy_models.ProductModel, error) {
	url := fmt.Sprintf("http://localhost:8081/products/getByIds")
	logging.Log.Info(url)
	var products []copy_models.ProductModel
	marshalledIds, err := json.Marshal(ids)
	if err != nil {
		logging.Log.Error("error while marshalling ids", err)
		return products, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(marshalledIds))
	if err != nil {
		logging.Log.Error("error while getting products", err)
		return products, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		logging.Log.Error("error while decoding products", err)
		return products, err
	}
	logging.Log.Info(products)
	return products, nil
}

func updateProductQtys(productsAndOrderQtys map[uint]uint) (bool, error) {
	marshalledQtys, err := json.Marshal(productsAndOrderQtys)
	if err != nil {
		return false, err
	}
	url := fmt.Sprintf("http://localhost:8081/products/updateQty")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(marshalledQtys))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, err
	}
	return true, nil
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	err := json_utils.JsonDecode(r, &order)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error while decoding order", http.StatusInternalServerError)
	}
	logging.Log.Info(order.CartID)

	cartExists, err, cart := checkIfCartExists(int(order.CartID))
	if !cartExists || err != nil {
		http_utils.SendErrorResponse(w, "Error while checking cart exists", http.StatusInternalServerError)
		return
	}

	if cart.UserID != order.UserID {
		http_utils.SendErrorResponse(w, "User ID conflict", http.StatusInternalServerError)
		return
	}
	tx := db.Begin()

	err = tx.Create(&order).Error

	if err != nil {
		tx.Rollback()
		http_utils.SendErrorResponse(w, "Error while creating order", http.StatusInternalServerError)
		return
	}

	order.OrderValue = cart.Value
	var orderProducts []models.OrderProduct
	var productsIds []uint
	for i := 0; i < len(cart.Items); i++ {
		cartProduct := cart.Items[i]
		productsIds = append(productsIds, cartProduct.ProductId)
		var orderProduct models.OrderProduct
		orderProduct.OrderID = order.ID
		orderProduct.ProductID = cartProduct.ProductId
		orderProduct.Quantity = cartProduct.Quantity
		orderProducts = append(orderProducts, orderProduct)
	}

	products, err := getProductsByIds(productsIds)
	if err != nil {
		tx.Rollback()
		logging.Log.Error("error in create order: getProductsByIds", err)
		http_utils.SendErrorResponse(w, "Error while getting products", http.StatusInternalServerError)
		return
	}

	productsAndOrderQtys := make(map[uint]uint)
	for i := 0; i < len(products); i++ {
		product := products[i]
		for j := 0; j < len(orderProducts); j++ {
			orderProduct := orderProducts[j]
			if orderProduct.ProductID == product.ID {
				if product.Quantity < orderProduct.Quantity {
					tx.Rollback()
					http_utils.SendErrorResponse(w, "Product out of stock", http.StatusInternalServerError)
					return
				}
				productsAndOrderQtys[product.ID] = orderProduct.Quantity
				break
			}
		}
	}

	tx.Commit()

	logging.Log.Info(len(orderProducts))
	err = BatchAssignProductsFromCart(orderProducts)
	if err != nil {
		tx.Rollback()
		http_utils.SendErrorResponse(w, "Error while assigning products", http.StatusInternalServerError)
		return
	}

	updateQtysSuccess, err := updateProductQtys(productsAndOrderQtys)
	if err != nil || !updateQtysSuccess {
		tx.Rollback()
		http_utils.SendErrorResponse(w, "Error while updating product qtys in createOrder", http.StatusInternalServerError)
		return
	}

	http_utils.SendSuccessResponse(w, "Order created successfully", http.StatusOK)
}
