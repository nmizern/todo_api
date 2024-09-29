package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type TODO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

var tasks = make([]TODO, 0)

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "appliction/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo TODO
	json.NewDecoder(r.Body).Decode(&newTodo)
	newTodo.ID = len(tasks) + 1
	tasks = append(tasks, newTodo)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTodo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
	}

	var updateTodo TODO
	json.NewDecoder(r.Body).Decode(&updateTodo)

	for i, todo := range tasks {
		if todo.ID == id {
			tasks[i].Name = updateTodo.Name
			tasks[i].Done = updateTodo.Done
			json.NewEncoder(w).Encode(tasks[i])

			return
		}

	}
	http.Error(w, "Task ID not found", http.StatusNotFound)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for i, todo := range tasks {
		if todo.ID == id {
			tasks = append(tasks[:i], tasks[:i+1]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	http.Error(w, "Task ID not found", http.StatusNotFound)
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the TODO API"))
	})

	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			updateTodo(w, r)
		} else if r.Method == "DELETE" {
			deleteTodo(w, r)
		}
	})

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getTodos(w, r)
		} else if r.Method == "POST" {
			createTodo(w, r)
		}
	})

	log.Println("Server is running")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
