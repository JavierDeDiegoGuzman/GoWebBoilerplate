package routes

import (
    "github.com/go-chi/chi/v5"
    "GoWebBoilerplate/internal/handlers"
)

// TodoRoutes configura las rutas para los todos
func TodoRoutes() chi.Router {
    r := chi.NewRouter()
    
    // Rutas para todos
    r.Get("/", handlers.ListTodos)
    r.Post("/", handlers.CreateTodo)
    r.Put("/{id}/toggle", handlers.ToggleTodo)  // Ruta más específica primero
    r.Delete("/{id}", handlers.DeleteTodo)
    
    return r
} 
