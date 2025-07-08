package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"todoApp/auth"
	"todoApp/database"
	"todoApp/models"
	"todoApp/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	//collect user as request body
	var req *models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//verify that user doesnt exist
	var user models.User
	err = database.Db.Where("username = ?", req.Username).First(&user).Error
	if err == nil {
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}

	//hash the password
	HashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.Password = HashPassword
	myuuid := uuid.NewString()
	req.ID = myuuid

	// add the user to the db
	err = database.Db.Create(&req).Error
	if err != nil {
		http.Error(w, "unable to create user", http.StatusInternalServerError)
		return
	}

	// send response back to client
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)

}

func Login(w http.ResponseWriter, r *http.Request) {
	var login models.User
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if user exists
	var user models.User
	err = database.Db.Where("username = ?", login.Username).First(&user).Error
	if err != nil {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	//compare passwords
	err = utils.ComparePassword(user.Password, login.Password)
	if err != nil {
		http.Error(w, "passwords do not match/incorrect", http.StatusBadRequest)
		return
	}

	//generate token
	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "unable to generate token", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(string)
    
    var todo models.Todo
    err := json.NewDecoder(r.Body).Decode(&todo)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    todo.ID = uuid.NewString()
    todo.UserID = userID
    todo.CreatedAt = time.Now()
    
    err = database.Db.Create(&todo).Error
    if err != nil {
        http.Error(w, "unable to create todo", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(todo)
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(string)
    
    var todos []models.Todo
    err := database.Db.Where("user_id = ?", userID).Find(&todos).Error
    if err != nil {
        http.Error(w, "unable to fetch todos", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(todos)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(string)
    todoID := mux.Vars(r)["id"]
    
    var todo models.Todo
    err := database.Db.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error
    if err != nil {
        http.Error(w, "todo not found", http.StatusNotFound)
        return
    }
    
    var updateData models.Todo
    err = json.NewDecoder(r.Body).Decode(&updateData)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    if updateData.Title != "" {
        todo.Title = updateData.Title
    }
    todo.Completed = updateData.Completed
    
    err = database.Db.Save(&todo).Error
    if err != nil {
        http.Error(w, "unable to update todo", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(string)
    todoID := mux.Vars(r)["id"]
    
    result := database.Db.Where("id = ? AND user_id = ?", todoID, userID).Delete(&models.Todo{})
    if result.Error != nil {
        http.Error(w, "unable to delete todo", http.StatusInternalServerError)
        return
    }
    
    if result.RowsAffected == 0 {
        http.Error(w, "todo not found", http.StatusNotFound)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "todo deleted successfully"})
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	user := "pana"

	message := fmt.Sprintf("Welcome, user %s!", user)
	w.Write([]byte(message))
}

