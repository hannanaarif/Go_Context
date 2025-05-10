package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hannanaarif/crudapi/db"
	"github.com/hannanaarif/crudapi/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateTodo handles the creation of a new todo
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	// Using request context instead of background context
	result, err := db.GetTodoCollection().InsertOne(r.Context(), todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todo.ID = result.InsertedID.(primitive.ObjectID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

// GetTodos returns all todos
func GetTodos(w http.ResponseWriter, r *http.Request) {
	var todos []models.Todo
	cursor, err := db.GetTodoCollection().Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	if err = cursor.All(r.Context(), &todos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// GetTodo returns a single todo by ID
func GetTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var todo models.Todo
	err = db.GetTodoCollection().FindOne(r.Context(), bson.M{"_id": objID}).Decode(&todo)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// UpdateTodo updates a todo by ID
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"title":       todo.Title,
			"description": todo.Description,
			"completed":   todo.Completed,
			"updatedAt":   todo.UpdatedAt,
		},
	}

	result := db.GetTodoCollection().FindOneAndUpdate(
		r.Context(),
		bson.M{"_id": objID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if result.Err() != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	if err := result.Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// DeleteTodo deletes a todo by ID
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := db.GetTodoCollection().DeleteOne(r.Context(), bson.M{"_id": objID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	// w.WriteHeader("User deleted Successfully",http.StatusNoContent)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("User deleted Successfully")

}
