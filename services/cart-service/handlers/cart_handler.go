package handlers

import (
	"cart-service/models"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/karthik-code78/ecom/shared/logging"
	"github.com/karthik-code78/ecom/shared/utils/http_utils"
	"github.com/karthik-code78/ecom/shared/utils/json_utils"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func ifUserExists(userId int) (bool, error) {
	url := fmt.Sprintf("http://localhost:8082/users/%d", userId)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return true, nil
}

func CreateCart(w http.ResponseWriter, r *http.Request) {
	var cart models.Cart
	err := json_utils.JsonDecode(r, &cart)
	if err != nil {
		http_utils.SendErrorResponse(w, "error decoding cart", http.StatusInternalServerError)
	}

	err = CartRepo.Create(&cart)
	if err != nil {
		http_utils.SendErrorResponse(w, "error creating cart", http.StatusInternalServerError)
	}

	err = json_utils.JsonEncode(w, cart)
	if err != nil {
		logging.Log.Error("error encoding cart", err.Error())
	}
}

func CreateCartFromUser(w http.ResponseWriter, r *http.Request) {
	var cart models.Cart
	var tempUserStruct struct {
		ID        uint   `json:"id"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	}

	err := json_utils.JsonDecode(r, &tempUserStruct)
	if err != nil {
		logging.Log.Error("error decoding temp user struct from Body", err.Error())
		http_utils.SendErrorResponse(w, "error decoding temp user struct from Body", http.StatusInternalServerError)
		return
	}
	cart.UserID = tempUserStruct.ID
	cart.Name = tempUserStruct.Firstname

	err = CartRepo.Create(&cart)
	if err != nil {
		logging.Log.Error("error creating cart", err.Error())
		http_utils.SendErrorResponse(w, "error creating cart", http.StatusInternalServerError)
		return
	}
	http_utils.SendSuccessResponse(w, "Cart creation success", http.StatusCreated)
}

func GetAllCarts(w http.ResponseWriter, r *http.Request) {
	var carts []models.Cart
	carts, err := CartRepo.FindAll()
	if err != nil {
		http_utils.SendErrorResponse(w, "error getting carts", http.StatusInternalServerError)
		return
	}

	err = json_utils.JsonEncode(w, carts)
	if err != nil {
		logging.Log.Error("error encoding carts", err.Error())
	}
}

func GetCartByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	logging.Log.Info(id)
	if err != nil {
		logging.Log.Error(err)
		http_utils.SendErrorResponse(w, "error getting id", http.StatusBadRequest)
		return
	}

	cart, err := CartRepo.FindById(uint(id))
	if err != nil {
		http_utils.SendErrorResponse(w, "error getting cart", http.StatusInternalServerError)
		return
	}

	logging.Log.Info(cart)
	err = json_utils.JsonEncode(w, &cart)
	if err != nil {
		http_utils.SendErrorResponse(w, "error encoding cart", http.StatusInternalServerError)
	}
}

func UpdateCartValueByCartId(id uint) (string, int) {
	var message = "Cart value updated successfully"
	cart, err := CartRepo.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message = "Record not found"
			return message, http.StatusNotFound
		} else {
			message = "Unable to find cart"
			return message, http.StatusInternalServerError
		}
	}
	for i := 0; i < len(cart.Items); i++ {
		cart.Value += cart.Items[i].Value
	}
	err = CartRepo.Update(cart)
	if err != nil {
		message = "Unable to update cart"
		return message, http.StatusInternalServerError
	}

	return message, http.StatusOK
}
