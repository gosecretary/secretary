package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"secretary/alpha/internal/config"
	"secretary/alpha/internal/domain"
	"secretary/alpha/internal/handlers"
	"secretary/alpha/internal/middleware"
	"secretary/alpha/internal/repository"
	"secretary/alpha/internal/service"
	"secretary/alpha/pkg/utils"
)

func main() {
	// Parse command line arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "server":
		runServer()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func generateRandomPassword(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func runServer() {
	// Parse server-specific flags
	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	devMode := serverCmd.Bool("dev", false, "Run in development mode with admin user")
	serverCmd.Parse(os.Args[2:])

	// Print banner
	bannerBytes, err := os.ReadFile("banner.txt")
	if err == nil {
		fmt.Println(string(bannerBytes))
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := repository.InitDB(cfg.Database.Driver, cfg.Database.FilePath)
	if err != nil {
		utils.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	resourceRepo := repository.NewResourceRepository(db)
	credentialRepo := repository.NewCredentialRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	accessRequestRepo := repository.NewAccessRequestRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	resourceService := service.NewResourceService(resourceRepo)
	credentialService := service.NewCredentialService(credentialRepo)
	permissionService := service.NewPermissionService(permissionRepo)
	accessRequestService := service.NewAccessRequestService(accessRequestRepo)

	// Create admin user in development mode
	if *devMode {
		// Generate random password
		password, err := generateRandomPassword(12)
		if err != nil {
			utils.Fatalf("Failed to generate random password: %v", err)
		}

		adminUser := &domain.User{
			Username: "admin",
			Email:    "admin@example.com",
			Password: password,
			Name:     "Admin User",
			Role:     "admin",
		}

		err = userService.CreateUser(context.Background(), adminUser)
		if err != nil {
			utils.Warnf("Failed to create admin user: %v", err)
		} else {
			utils.Info("----------------------------------------")
			utils.Info("Development mode: Created admin user")
			utils.Info("Username: admin")
			utils.Info("Password: " + password)
			utils.Info("Email: admin@example.com")
			utils.Info("----------------------------------------")
		}
	}

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	resourceHandler := handlers.NewResourceHandler(resourceService)
	credentialHandler := handlers.NewCredentialHandler(credentialService)
	permissionHandler := handlers.NewPermissionHandler(permissionService)
	accessRequestHandler := handlers.NewAccessRequestHandler(accessRequestService)

	// Initialize router
	router := handlers.NewRouter()
	router.RegisterHandlers(
		userHandler,
		resourceHandler,
		credentialHandler,
		permissionHandler,
		accessRequestHandler,
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
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Channel to listen for errors coming from the listener
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

	// Channel to listen for an interrupt or terminate signal from the OS
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown
	select {
	case err := <-serverErrors:
		utils.Errorf("Server error: %v", err)
	case sig := <-shutdown:
		utils.Infof("Received signal %v, shutting down server...", sig)
	}

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		utils.Errorf("Server forced to shutdown: %v", err)
	}

	utils.Info("Server exiting")
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  secretary server [--dev]  Start the server")
	fmt.Println("\nOptions:")
	fmt.Println("  --dev  Run in development mode with admin user")
}
