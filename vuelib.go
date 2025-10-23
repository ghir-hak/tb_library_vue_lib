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

type Todo struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Done  bool   `json:"done"`
}

//export addTodo
func addTodo(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

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

	err = db.Put(todo.ID, todoBytes)
	if err != nil {
		return fail(h, err, 500)
	}

	return 0
}

//export getTodo
func getTodo(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

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

//export listTodos
func listTodos(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	db, err := database.New("/todo/list")
	if err != nil {
		return fail(h, err, 500)
	}

	// Get all keys (this is a simplified approach)
	// In a real app, you'd want to implement proper listing
	keys := []string{"todo1", "todo2", "todo3"} // Simplified for demo
	
	var todos []Todo
	for _, key := range keys {
		value, err := db.Get(key)
		if err == nil {
			var todo Todo
			json.Unmarshal(value, &todo)
			todos = append(todos, todo)
		}
	}

	todosJson, err := json.Marshal(todos)
	if err != nil {
		return fail(h, err, 500)
	}

	h.Write(todosJson)
	h.Return(200)

	return 0
}
