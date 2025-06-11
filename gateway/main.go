package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"secretary/alpha/api"
	"secretary/alpha/internal"
	"secretary/alpha/internal/config"
	"secretary/alpha/internal/middleware"
	"secretary/alpha/internal/server"
	"secretary/alpha/storage"
)

func main() {
	// Parse command line flags
	var (
		healthCheck = flag.Bool("health-check", false, "perform health check and exit")
	)
	flag.Parse()

	// Load configuration
	cfg := config.Load()

	// Health check mode
	if *healthCheck {
		if !storage.DatabaseHealthCheck() {
			log.Printf("Health check failed: database connection error")
			os.Exit(1)
		}
		log.Println("Health check passed")
		os.Exit(0)
	}

	// Initialize database
	if !storage.DatabaseInit() {
		log.Fatal("Failed to initialize database")
	}

	// Initialize middleware
	middleware.InitializeMiddleware()

	// Run fixtures (initial data setup)
	internal.RunFixtures()

	// Show banner
	internal.ShowBanner("./banner.txt")

	// Create server
	srv := server.NewServer(cfg)

	// Setup routes with security middleware
	mux := http.NewServeMux()

	// Public endpoints (no authentication required)
	mux.Handle("/", middleware.PublicEndpoint(http.HandlerFunc(api.HomeAPI)))
	mux.Handle("/api/hz", middleware.PublicEndpoint(http.HandlerFunc(api.HealthCheckAPI)))
	mux.Handle("/api/user/login", middleware.PublicEndpoint(http.HandlerFunc(api.LoginAPI)))

	// Protected endpoints (authentication required)
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/asksfor", api.AskAPI)
	protectedMux.HandleFunc("/api/credential", api.CredentialAPI)
	protectedMux.HandleFunc("/api/resource", api.ResourceAPI)
	protectedMux.HandleFunc("/api/resource/credential", api.ResourceCredentialAPI)
	protectedMux.HandleFunc("/api/user", api.UserAPI)
	protectedMux.HandleFunc("/api/user/self", api.SelfAPI)
	protectedMux.HandleFunc("/api/user/logout", api.LogoutAPI)

	// Apply authentication and CSRF middleware to protected endpoints
	protectedHandler := middleware.AuthenticationMiddleware(
		middleware.CSRFMiddleware(protectedMux),
	)

	mux.Handle("/api/", http.StripPrefix("", protectedHandler))

	// Configure server with security settings
	srv.SetupRoutes(mux)

	// Start server
	log.Printf("Starting Secretary server...")
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
