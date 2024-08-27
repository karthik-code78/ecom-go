package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/karthik-code78/ecom/shared/auth"
	"github.com/karthik-code78/ecom/shared/logging"
	"github.com/karthik-code78/ecom/shared/utils/http_utils"
	"github.com/karthik-code78/ecom/shared/utils/json_utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"user-service/models"
	"user-service/repository"
)

var db *gorm.DB

var usersRepo repository.UserRepository

func SetDatabase(database *gorm.DB) {
	logging.Log.Info(database)
	db = database

	usersRepo = repository.NewUserRepository(db)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json_utils.JsonDecode(r, &credentials)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error decoding the credentials", http.StatusBadRequest)
		return
	}
	user, err := usersRepo.FindByEmailId(credentials.Email)
	if err != nil {
		http_utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	credsCompare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if credsCompare != nil {
		logging.Log.Error(" Err in comparing : ", err)
		logging.Log.Error("credsCompare error : ", credsCompare)
		http_utils.SendErrorResponse(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.Email)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	err = json.NewEncoder(w).Encode(map[string]string{"token": token})
	if err != nil {
		http_utils.SendErrorResponse(w, "Error encoding the token", http.StatusInternalServerError)
	}
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := map[string]string{
		"firstname": r.URL.Query().Get("firstname"),
		"lastname":  r.URL.Query().Get("lastname"),
		"sortBy":    r.URL.Query().Get("sortBy"),
		"sortDir":   r.URL.Query().Get("sortDir"),
	}

	users, err := usersRepo.FindAll(queryParams)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error getting users", http.StatusInternalServerError)
	}

	err = json_utils.JsonEncode(w, users)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error getting users", http.StatusInternalServerError)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http_utils.SendErrorResponse(w, "Error parsing id", http.StatusBadRequest)
		return
	}
	user, err := usersRepo.FindById(uint(id))
	if err != nil {
		http_utils.SendErrorResponse(w, "Error getting users", http.StatusInternalServerError)
		return
	}
	err = json_utils.JsonEncode(w, user)
	if err != nil {
		logging.Log.Error("Error encoding users", err.Error())
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	err := json_utils.JsonDecode(r, &user)
	if err != nil {
		logging.Log.Error("Error encoding users", err)
		http_utils.SendErrorResponse(w, "Error decoding user", http.StatusBadRequest)
		return
	}
	logging.Log.Info(user.Email)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http_utils.SendErrorResponse(w, "Password hashing failed, please check", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)

	tx := db.Begin()

	if err = tx.Create(&user).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			http_utils.SendErrorResponse(w, "Email is already registered", http.StatusConflict)
			return
		}
		http_utils.SendErrorResponse(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	var userResp models.UserResponse

	userResp.ID = user.ID
	userResp.Firstname = user.Firstname
	userResp.Lastname = user.Lastname

	cartCreated, err := createCartForUser(userResp)

	if err != nil || !cartCreated {
		tx.Rollback()
		http_utils.SendErrorResponse(w, "Failed to create cart", http.StatusInternalServerError)
		return
	}

	tx.Commit()
	http_utils.SendSuccessResponse(w, "User and cart created successfully", http.StatusCreated)
}

func createCartForUser(userResp models.UserResponse) (bool, error) {
	// Define the payload
	marshalledUser, err := json.Marshal(userResp)
	if err != nil {
		return false, err
	}

	url := "http://localhost:8083/cart/withUserBody"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(marshalledUser))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return false, fmt.Errorf("cart service returned error")
	}

	return true, nil
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http_utils.SendErrorResponse(w, "Error parsing id", http.StatusBadRequest)
		return
	}

	user, err := usersRepo.FindById(uint(id))
	if err != nil {
		http_utils.SendErrorResponse(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	if user.ID == 0 {
		http_utils.SendErrorResponse(w, "User with Id"+string(rune(id))+"not found", http.StatusExpectationFailed)
		return
	}

	err = json_utils.JsonDecode(r, &user)
	if err != nil {
		http_utils.SendErrorResponse(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	err = usersRepo.Update(user)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	err = json_utils.JsonEncode(w, user)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error encoding the user", http.StatusInternalServerError)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logging.Log.Error("Error parsing id", http.StatusBadRequest)
		return
	}
	user, err := usersRepo.FindById(uint(id))
	if err != nil {
		http_utils.SendErrorResponse(w, "Error getting user", http.StatusInternalServerError)
		return
	}
	err = usersRepo.Delete(user)
	if err != nil {
		http_utils.SendErrorResponse(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
