package main

import (
	"context"
	"fmt"
	"io/ioutil"
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
	"secretary/alpha/pkg/utils"
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
		utils.Fatalf("Failed to load configuration")
	}

	// Initialize database
	db, err := repository.InitDB(cfg.Database.Driver, cfg.Database.FilePath)
	if err != nil {
		utils.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	resourceRepo := repository.NewResourceRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	credentialRepo := repository.NewCredentialRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	accessRequestRepo := repository.NewAccessRequestRepository(db)
	ephemeralCredentialRepo := repository.NewEphemeralCredentialRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	resourceService := service.NewResourceService(resourceRepo)
	sessionService := service.NewSessionService(sessionRepo)
	credentialService := service.NewCredentialService(credentialRepo)
	permissionService := service.NewPermissionService(permissionRepo)
	accessRequestService := service.NewAccessRequestService(accessRequestRepo)
	ephemeralCredentialService := service.NewEphemeralCredentialService(ephemeralCredentialRepo)

	// Initialize HTTP handlers
	router := httpapi.NewRouter(
		userService,
		resourceService,
		credentialService,
		permissionService,
		sessionService,
		accessRequestService,
		ephemeralCredentialService,
	)

	// Add middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recovery)
	router.Use(middleware.CORS)

	// Create server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		utils.Infof("Server starting on %s", serverAddr)
		if cfg.Server.TLSCertPath != "" && cfg.Server.TLSKeyPath != "" {
			utils.Infof("Using TLS with certificate: %s", cfg.Server.TLSCertPath)
			if err := srv.ListenAndServeTLS(cfg.Server.TLSCertPath, cfg.Server.TLSKeyPath); err != nil && err != http.ErrServerClosed {
				utils.Fatalf("Failed to start secure server: %v", err)
			}
		} else {
			utils.Warn("WARNING: Running in HTTP mode. For production, use HTTPS.")
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				utils.Fatalf("Failed to start server: %v", err)
			}
		}
	}()

	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		utils.Errorf("Server error: %v", err)
	case sig := <-shutdown:
		utils.Infof("Received signal %v, shutting down server...", sig)
	}

	// Give outstanding requests a deadline for completion.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		utils.Errorf("Server forced to shutdown: %v", err)
	}

	utils.Info("Server exiting")
}
