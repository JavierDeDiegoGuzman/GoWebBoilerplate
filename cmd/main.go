package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"GoWebBoilerplate/internal/routes"
	"GoWebBoilerplate/internal/utils"
)

func main() {
	// Inicializar logger
	if err := utils.InitLogger(); err != nil {
		panic(err)
	}

	// Inicializar assets
	if err := utils.InitAssets(); err != nil {
		utils.LogServer("initializing assets", err)
		os.Exit(1)
	}

	// Inicializar base de datos
	if err := utils.InitDB("todos.db"); err != nil {
		utils.LogDB("initializing database", err)
		os.Exit(1)
	}
	if err := utils.CreateTables(); err != nil {
		utils.LogDB("creating tables", err)
		os.Exit(1)
	}
	defer utils.CloseDB()

	// Crear el router usando SetupRouter
	r := routes.SetupRouter()

	// Configurar el servidor
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Manejar señales de terminación
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Iniciar el servidor en una goroutine
	go func() {
		utils.LogServer("server started", nil)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.LogServer("server error", err)
			os.Exit(1)
		}
	}()

	// Esperar señal de terminación
	<-stop
	utils.LogServer("shutting down server", nil)
}