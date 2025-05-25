package handlers

import (
	//"GoWebBoilerplate/internal/views"
	"fmt"
	"net/http"
)

// Home maneja la página principal
func Home(w http.ResponseWriter, r *http.Request) {
    //views.TodoPage().Render(r.Context(), w)
	fmt.Println("home page??")
	fmt.Fprintln(w, "¡Hola desde mi servidor web!")
} 
