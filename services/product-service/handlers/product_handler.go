package handlers

import (
	"github.com/gorilla/mux"
	"github.com/karthik-code78/ecom/shared/logging"
	"github.com/karthik-code78/ecom/shared/utils/http_utils"
	"github.com/karthik-code78/ecom/shared/utils/json_utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"product-service/models"
	"product-service/repository"
	"strconv"
)

var db *gorm.DB

var productsRepo repository.ProductRepository

func SetDatabase(database *gorm.DB) {
	logging.Log.Info(database)
	db = database

	productsRepo = repository.NewProductRepository(db)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	queryParams := map[string]string{
		"name":        r.URL.Query().Get("name"),
		"description": r.URL.Query().Get("description"),
		"minPrice":    r.URL.Query().Get("minPrice"),
		"maxPrice":    r.URL.Query().Get("maxPrice"),
		"minQuantity": r.URL.Query().Get("minQuantity"),
		"maxQuantity": r.URL.Query().Get("maxQuantity"),
		"sortBy":      r.URL.Query().Get("sortBy"),
		"sortDir":     r.URL.Query().Get("sortDir"),
	}

	products, err := productsRepo.FindAll(queryParams)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error getting products", http.StatusInternalServerError)
	}

	err = json_utils.JsonEncode(w, products)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error getting products", http.StatusInternalServerError)
	}
}

func GetProductsByIds(w http.ResponseWriter, r *http.Request) {
	var ids []uint
	err := json_utils.JsonDecode(r, &ids)
	if err != nil {
		logging.Log.Error(err)
		http_utils.SendErrorResponse(w, "Error decoding product Ids", http.StatusInternalServerError)
		return
	}
	logging.Log.Info(ids)
	var products []models.Product
	products, err = productsRepo.FindByIds(ids)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error getting products", http.StatusInternalServerError)
		return
	}
	err = json_utils.JsonEncode(w, products)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error encoding products", http.StatusInternalServerError)
	}
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	logging.Log.Info(uint(id))
	if err != nil {
		http_utils.SendErrorResponse(w, "Error parsing id", http.StatusBadRequest)
		return
	}
	product, err := productsRepo.FindById(uint(id))
	logging.Log.Info(err)
	logging.Log.Info(product)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error getting product", http.StatusInternalServerError)
		return
	}
	err = json_utils.JsonEncode(w, product)
	if err != nil {
		logging.Log.Error("Error encoding product", err.Error())
		http_utils.SendErrorResponse(w, "Error encoding product", http.StatusInternalServerError)
	}
}

func UpdateProductQuantity(w http.ResponseWriter, r *http.Request) {
	var products map[uint]uint
	err := json_utils.JsonDecode(r, &products)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error decoding products", http.StatusInternalServerError)
		return
	}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	for productId, orderQty := range products {
		var product models.Product
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("ID = ?", productId).First(&product).Error
		if err != nil {
			tx.Rollback()
			http_utils.SendErrorResponse(w, "Error getting products", http.StatusInternalServerError)
			return
		}
		logging.Log.Info(product.Description)

		if orderQty > product.Quantity {
			orderQty = product.Quantity
		}

		logging.Log.Info(orderQty)

		err = tx.Model(&product).Where("id = ?", productId).Update("quantity", (product.Quantity - orderQty)).Error
		if err != nil {
			logging.Log.Error(err)
			tx.Rollback()
			http_utils.SendErrorResponse(w, "Error modifying quantity", http.StatusInternalServerError)
			return
		}
		tx.Commit()
	}

	http_utils.SendSuccessResponse(w, "Successfully updated the Qtys", http.StatusOK)
}

func GetProductByQuantity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http_utils.SendErrorResponse(w, "Error parsing id", http.StatusBadRequest)
		return
	}
	qty, err := strconv.Atoi(vars["qty"])
	if err != nil {
		http_utils.SendErrorResponse(w, "Error parsing qty", http.StatusBadRequest)
		return
	}
	product, err := productsRepo.FindByIdAndQuantity(uint(id), uint(qty))
	if err != nil {
		http_utils.SendErrorResponse(w, "Error getting product", http.StatusInternalServerError)
		return
	}

	err = json_utils.JsonEncode(w, product)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error encoding product", http.StatusInternalServerError)
		return
	}
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json_utils.JsonDecode(r, &product)
	if err != nil {
		logging.Log.Error("Error encoding products", err)
		http_utils.SendErrorResponse(w, "Error decoding product", http.StatusBadRequest)
		return
	}
	err = productsRepo.Create(&product)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error creating product", http.StatusInternalServerError)
		return
	}
	err = json_utils.JsonEncode(w, product)
	if err != nil {
		logging.Log.Error("Error encoding product", err.Error())
	}
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http_utils.SendErrorResponse(w, "Error parsing id", http.StatusBadRequest)
		return
	}

	product, err := productsRepo.FindById(uint(id))
	if err != nil {
		http_utils.SendErrorResponse(w, "Error getting product", http.StatusInternalServerError)
		return
	}

	if product.ID == 0 {
		http_utils.SendErrorResponse(w, "Product with Id"+string(rune(id))+"not found", http.StatusExpectationFailed)
		return
	}

	err = json_utils.JsonDecode(r, &product)
	if err != nil {
		http_utils.SendErrorResponse(w, "Invalid product data", http.StatusBadRequest)
		return
	}

	err = productsRepo.Update(product)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error updating product", http.StatusInternalServerError)
		return
	}

	err = json_utils.JsonEncode(w, product)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error encoding the product", http.StatusInternalServerError)
	}
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logging.Log.Error("Error parsing id", http.StatusBadRequest)
		return
	}
	product, err := productsRepo.FindById(uint(id))
	if err != nil {
		http_utils.SendErrorResponse(w, "Error getting product", http.StatusInternalServerError)
		return
	}
	err = productsRepo.Delete(product)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error deleting product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
