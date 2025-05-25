package handlers

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "GoWebBoilerplate/internal/models"
    "GoWebBoilerplate/internal/views"
)

// ListTodos maneja la lista de todos
func ListTodos(w http.ResponseWriter, r *http.Request) {
    todo := &models.Todo{}
    todos, err := todo.GetAll()
    if err != nil {
        http.Error(w, "Error al obtener todos", http.StatusInternalServerError)
        return
    }

    views.TodoPage(todos).Render(r.Context(), w)
}

// CreateTodo maneja la creación de un nuevo todo
func CreateTodo(w http.ResponseWriter, r *http.Request) {
    task := r.FormValue("task")
    if task == "" {
        http.Error(w, "Task is required", http.StatusBadRequest)
        return
    }

    todo := models.NewTodo(task)
    if err := todo.Save(); err != nil {
        http.Error(w, "Error al crear todo", http.StatusInternalServerError)
        return
    }

    // Renderizar solo el nuevo todo
    views.TodoItem(todo).Render(r.Context(), w)
}

// UpdateTodo maneja la actualización de un todo
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    if id == "" {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    todo := &models.Todo{}
    if err := todo.GetByID(id); err != nil {
        http.Error(w, "Todo not found", http.StatusNotFound)
        return
    }

    task := r.FormValue("task")
    if task != "" {
        todo.Task = task
    }

    if err := todo.Save(); err != nil {
        http.Error(w, "Error al actualizar todo", http.StatusInternalServerError)
        return
    }

    views.TodoItem(todo).Render(r.Context(), w)
}

// DeleteTodo maneja la eliminación de un todo
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    if id == "" {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    todo := &models.Todo{ID: id}
    if err := todo.Delete(); err != nil {
        http.Error(w, "Error al eliminar todo", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

// ToggleTodo maneja el cambio de estado de un todo
func ToggleTodo(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    if id == "" {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    todo := &models.Todo{ID: id}
    if err := todo.GetByID(id); err != nil {
        http.Error(w, "Todo not found", http.StatusNotFound)
        return
    }

    todo.Done = !todo.Done
    
    if err := todo.Save(); err != nil {
        http.Error(w, "Error al actualizar todo", http.StatusInternalServerError)
        return
    }

    views.TodoItem(todo).Render(r.Context(), w)
} 
