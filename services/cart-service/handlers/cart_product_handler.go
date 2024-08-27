package handlers

import (
	"cart-service/models"
	"cart-service/models/copy_models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/karthik-code78/ecom/shared/logging"
	"github.com/karthik-code78/ecom/shared/utils/http_utils"
	"github.com/karthik-code78/ecom/shared/utils/json_utils"
	"net/http"
	"strconv"
)

func IsProductExists(productID int) (bool, error, *copy_models.Product) {
	url := fmt.Sprintf("http://localhost:8081/products/%d", productID)
	logging.Log.Info(url)
	resp, err := http.Get(url)

	var product copy_models.Product
	if err != nil {
		return false, err, &product
	}
	defer resp.Body.Close()
	logging.Log.Info(resp.Body)
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		logging.Log.Error(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return false, nil, &product
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode), &product
	}

	return true, nil, &product
}

func GetAllByCartId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var cartProducts []models.CartProduct
	cartProducts, err = CartProductsRepo.FindAllByCartId(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(cartProducts) == 0 {
		http.Error(w, fmt.Sprintf("No linked carts with id %d found", id), http.StatusNotFound)
	}

	err = json_utils.JsonEncode(w, cartProducts)
	if err != nil {
		http.Error(w, "error encoding the cartProducts", http.StatusInternalServerError)
	}
}

func AddProductsToCart(w http.ResponseWriter, r *http.Request) {
	var cartProduct models.CartProduct
	err := json_utils.JsonDecode(r, &cartProduct)

	if err != nil {
		http_utils.SendErrorResponse(w, "Error decoding cartProduct", http.StatusInternalServerError)
		return
	}

	if cartProduct.Quantity < 1 {
		http_utils.SendErrorResponse(w, "Quantity must be greater than 1", http.StatusBadRequest)
		return
	}

	exists, err, product := IsProductExists(int(cartProduct.ProductId))
	if err != nil {
		http_utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		http_utils.SendErrorResponse(w, "Product not found - cannot add to cart", http.StatusNotFound)
		return
	}

	if cartProduct.Quantity > product.Quantity {
		http_utils.SendErrorResponse(w, "Product out of stock", http.StatusBadRequest)
		return
	} else {
		//cartProduct.AvailableQty = product.Quantity
		cartProduct.Value = product.Price * float64(cartProduct.Quantity)

		if err = CartProductsRepo.Update(&cartProduct); err != nil {
			http_utils.SendErrorResponse(w, "Unable to update cart-products", http.StatusInternalServerError)
			return
		}

		message, statusCode := UpdateCartValueByCartId(cartProduct.CartID)
		if statusCode != http.StatusOK {
			http_utils.SendErrorResponse(w, message, statusCode)
			return
		}
	}

	http_utils.SendSuccessResponse(w, "Successfully added the products", http.StatusOK)
}
