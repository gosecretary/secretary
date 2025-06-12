package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpapi "secretary/alpha/api/http"
	"secretary/alpha/internal/config"
	"secretary/alpha/internal/middleware"
	"secretary/alpha/internal/repository"
	"secretary/alpha/internal/service"
)

func main() {
	// Display banner
	bannerBytes, err := ioutil.ReadFile("banner.txt")
	if err == nil {
		fmt.Println(string(bannerBytes))
	}

	// Load configuration
	cfg := config.Load()
	if cfg == nil {
		log.Fatalf("Failed to load configuration")
	}

	// Initialize database
	db, err := repository.InitDB(cfg.Database.Driver, cfg.Database.FilePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	resourceRepo := repository.NewResourceRepository(db)
	credentialRepo := repository.NewCredentialRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	accessRequestRepo := repository.NewAccessRequestRepository(db)
	ephemeralCredentialRepo := repository.NewEphemeralCredentialRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	resourceService := service.NewResourceService(resourceRepo)
	credentialService := service.NewCredentialService(credentialRepo)
	permissionService := service.NewPermissionService(permissionRepo)
	sessionService := service.NewSessionService(sessionRepo)
	accessRequestService := service.NewAccessRequestService(accessRequestRepo)
	ephemeralCredentialService := service.NewEphemeralCredentialService(ephemeralCredentialRepo)

	// Initialize router with middleware
	router := httpapi.NewRouter(
		userService,
		resourceService,
		credentialService,
		permissionService,
		sessionService,
		accessRequestService,
		ephemeralCredentialService,
	)

	router.Use(middleware.Logger)
	router.Use(middleware.Recovery)
	router.Use(middleware.CORS)

	// Initialize server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s", serverAddr)
		if cfg.Server.TLSCertPath != "" && cfg.Server.TLSKeyPath != "" {
			log.Printf("Using TLS with certificate: %s", cfg.Server.TLSCertPath)
			if err := server.ListenAndServeTLS(cfg.Server.TLSCertPath, cfg.Server.TLSKeyPath); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Failed to start secure server: %v", err)
			}
		} else {
			log.Printf("WARNING: Running in HTTP mode. For production, use HTTPS.")
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Failed to start server: %v", err)
			}
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
