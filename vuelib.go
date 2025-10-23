package lib

import (
	"encoding/json"

	"github.com/taubyte/go-sdk/database"
	"github.com/taubyte/go-sdk/event"
	http "github.com/taubyte/go-sdk/http/event"
)

func fail(h http.Event, err error, code int) uint32 {
	h.Write([]byte(err.Error()))
	h.Return(code)
	return 1
}

func setCORSHeaders(h http.Event) {
	h.Headers().Set("Access-Control-Allow-Origin", "*")
	h.Headers().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	h.Headers().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

type Todo struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Done  bool   `json:"done"`
}

// POST /api/todo
//export addTodo
func addTodo(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	setCORSHeaders(h)

	// Create & Open the database
	db, err := database.New("/todo/list")
	if err != nil {
		return fail(h, err, 500)
	}

	// Decode the request body
	reqDec := json.NewDecoder(h.Body())
	defer h.Body().Close()

	// Decode the request body
	var todo Todo
	err = reqDec.Decode(&todo)
	if err != nil {
		return fail(h, err, 500)
	}

	// Put the todo into the database
	todoBytes, err := json.Marshal(todo)
	if err != nil {
		return fail(h, err, 500)
	}

	// Store with proper key prefix
	key := "/todo/" + todo.ID
	err = db.Put(key, todoBytes)
	if err != nil {
		return fail(h, err, 500)
	}

	h.Write([]byte("Todo added successfully"))
	h.Return(200)
	return 0
}

// GET /api/todo?id={id}
//export getTodo
func getTodo(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	setCORSHeaders(h)

	id, err := h.Query().Get("id")
	if err != nil {
		return fail(h, err, 400)
	}

	db, err := database.New("/todo/list")
	if err != nil {
		return fail(h, err, 500)
	}

	value, err := db.Get(id)
	if err != nil {
		return fail(h, err, 500)
	}

	h.Write(value)
	h.Return(200)

	return 0
}

// DELETE /api/todo?id={id}
//export deleteTodo
func deleteTodo(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	setCORSHeaders(h)

	id, err := h.Query().Get("id")
	if err != nil {
		return fail(h, err, 400)
	}

	db, err := database.New("/todo/list")
	if err != nil {
		return fail(h, err, 500)
	}

	// Delete the todo
	key := "/todo/" + id
	err = db.Delete(key)
	if err != nil {
		return fail(h, err, 500)
	}

	h.Write([]byte("Todo deleted successfully"))
	h.Return(200)

	return 0
}

// GET /api/todos
//export listTodos
func listTodos(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	setCORSHeaders(h)

	db, err := database.New("/todo/list")
	if err != nil {
		return fail(h, err, 500)
	}

	// List all keys with prefix
	keys, err := db.List("/todo/")
	if err != nil {
		// Return empty array on error instead of failing
		emptyArray := []Todo{}
		todosJson, _ := json.Marshal(emptyArray)
		h.Write(todosJson)
		h.Return(200)
		return 0
	}

	var todos []Todo
	for _, key := range keys {
		value, err := db.Get(key)
		if err == nil {
			var todo Todo
			json.Unmarshal(value, &todo)
			todos = append(todos, todo)
		}
	}

	// Always return an array, even if empty
	todosJson, err := json.Marshal(todos)
	if err != nil {
		// Fallback to empty array
		emptyArray := []Todo{}
		todosJson, _ = json.Marshal(emptyArray)
	}

	h.Write(todosJson)
	h.Return(200)

	return 0
}
