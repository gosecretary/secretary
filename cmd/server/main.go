package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"secretary/api/http"
	"secretary/internal/config"
	"secretary/internal/middleware"
	"secretary/internal/repository"
	"secretary/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := repository.InitDB(cfg.Database.Driver, cfg.Database.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	resourceRepo := repository.NewResourceRepository(db)
	credentialRepo := repository.NewCredentialRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	resourceService := service.NewResourceService(resourceRepo)
	credentialService := service.NewCredentialService(credentialRepo)
	permissionService := service.NewPermissionService(permissionRepo)

	// Initialize router with middleware
	router := http.NewRouter(userService, resourceService, credentialService, permissionService)
	router.Use(middleware.Logger)
	router.Use(middleware.Recovery)
	router.Use(middleware.CORS)

	// Initialize server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	// TODO: Implement graceful shutdown
} 