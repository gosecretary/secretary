package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"secretary/alpha/internal/config"
	"secretary/alpha/internal/middleware"
	"secretary/alpha/pkg/utils"
)

type Server struct {
	httpServer *http.Server
	config     *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

func (s *Server) setupTLS() *tls.Config {
	if s.config.Server.TLSCertPath == "" || s.config.Server.TLSKeyPath == "" {
		return nil
	}

	tlsConfig := &tls.Config{
		MinVersion:       tls.VersionTLS12,
		MaxVersion:       tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		PreferServerCipherSuites: true,
	}

	return tlsConfig
}

func (s *Server) SetupRoutes(mux *http.ServeMux) {
	// Apply rate limiting to all routes
	rateLimitedMux := middleware.RateLimitMiddleware(mux)

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port),
		Handler:      rateLimitedMux,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
		IdleTimeout:  s.config.Server.IdleTimeout,
		TLSConfig:    s.setupTLS(),
		// Security headers
		ErrorLog: utils.GetStandardLogger(),
	}

	// Disable HTTP/2 for now due to potential security issues
	s.httpServer.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
}

func (s *Server) Start() error {
	utils.Logger("info", fmt.Sprintf("Starting server on %s", s.httpServer.Addr))

	// Create a channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		if s.config.Server.TLSCertPath != "" && s.config.Server.TLSKeyPath != "" {
			utils.Logger("info", "Starting HTTPS server with TLS")
			serverErr <- s.httpServer.ListenAndServeTLS(s.config.Server.TLSCertPath, s.config.Server.TLSKeyPath)
		} else {
			utils.Logger("warn", "Starting HTTP server without TLS - this is insecure for production!")
			serverErr <- s.httpServer.ListenAndServe()
		}
	}()

	// Wait for interrupt signal or server error
	select {
	case err := <-serverErr:
		if err != http.ErrServerClosed {
			return fmt.Errorf("server failed to start: %w", err)
		}
	case sig := <-quit:
		utils.Logger("info", fmt.Sprintf("Received signal %v, shutting down server...", sig))
		return s.Shutdown()
	}

	return nil
}

func (s *Server) Shutdown() error {
	utils.Logger("info", "Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := s.httpServer.Shutdown(ctx); err != nil {
		utils.Logger("err", fmt.Sprintf("Server shutdown error: %v", err))
		return err
	}

	utils.Logger("info", "Server gracefully stopped")
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.httpServer.Handler.ServeHTTP(w, r)
}
