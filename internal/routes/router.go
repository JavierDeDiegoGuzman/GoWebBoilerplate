package routes

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"GoWebBoilerplate/internal/handlers"
	"GoWebBoilerplate/internal/utils"
)

// LoggerMiddleware es un middleware personalizado que usa nuestro sistema de logging
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		
		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		utils.LogRequest(r.Method, r.URL.Path, ww.Status(), duration, nil)
	})
}

// SetupRouter configura el router principal
func SetupRouter() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Use(middleware.GetHead)
	r.Use(LoggerMiddleware)

	// Rutas estáticas
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Rutas de la aplicación
	r.Get("/", handlers.Home)
	
	// Rutas de todos
	r.Mount("/todos", TodoRoutes())

	return r
} 