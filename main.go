package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type TODO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

type TasksManager struct {
	mu    sync.Mutex
	tasks []TODO
}

func (tm *TasksManager) getTodos(w http.ResponseWriter, r *http.Request) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	w.Header().Set("Content-Type", "appliction/json")
	json.NewEncoder(w).Encode(tm.tasks)
}

func (tm *TasksManager) createTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo TODO
	json.NewDecoder(r.Body).Decode(&newTodo)

	tm.mu.Lock()
	defer tm.mu.Unlock()

	newTodo.ID = len(tm.tasks) + 1
	tm.tasks = append(tm.tasks, newTodo)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTodo)
}

func (tm *TasksManager) updateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updateTodo TODO
	json.NewDecoder(r.Body).Decode(&updateTodo)

	for i, todo := range tm.tasks {
		if todo.ID == id {
			tm.tasks[i].Name = updateTodo.Name
			tm.tasks[i].Done = updateTodo.Done
			json.NewEncoder(w).Encode(tm.tasks[i])

			return
		}

	}
	http.Error(w, "Task ID not found", http.StatusNotFound)
}

func (tm *TasksManager) deleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for i, todo := range tm.tasks {
		if todo.ID == id {
			tm.tasks = append(tm.tasks[:i], tm.tasks[:i+1]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	http.Error(w, "Task ID not found", http.StatusNotFound)
}

func main() {

	tasksManager := &TasksManager{
		tasks: make([]TODO, 0),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the TODO API"))
	})

	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			tasksManager.updateTodo(w, r)
		} else if r.Method == "DELETE" {
			tasksManager.deleteTodo(w, r)
		}
	})

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tasksManager.getTodos(w, r)
		} else if r.Method == "POST" {
			tasksManager.createTodo(w, r)
		}
	})

	log.Println("Server is running")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
