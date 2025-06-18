package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"secretary/alpha/internal/domain"
	"secretary/alpha/pkg/utils"

	"github.com/google/uuid"
)

type proxyService struct {
	sessionCommandService   domain.SessionCommandService
	sessionRecordingService domain.SessionRecordingService
	securityAlertService    domain.SecurityAlertService
	activeConnections       map[string]*ProxyConnection
	mu                      sync.RWMutex
}

type ProxyConnection struct {
	ID            string
	SessionID     string
	UserID        string
	ResourceID    string
	Protocol      string
	LocalPort     int
	RemoteHost    string
	RemotePort    int
	Status        string
	listener      net.Listener
	connections   []net.Conn
	recordingFile *os.File
	cancel        context.CancelFunc
}

func NewProxyService(
	sessionCommandService domain.SessionCommandService,
	sessionRecordingService domain.SessionRecordingService,
	securityAlertService domain.SecurityAlertService,
) domain.ProxyService {
	return &proxyService{
		sessionCommandService:   sessionCommandService,
		sessionRecordingService: sessionRecordingService,
		securityAlertService:    securityAlertService,
		activeConnections:       make(map[string]*ProxyConnection),
	}
}

func (s *proxyService) CreateProxy(ctx context.Context, sessionID, protocol, remoteHost string, remotePort int) (*domain.ProxyConnection, error) {
	proxyID := uuid.New().String()

	// Find available local port
	localPort, err := s.findAvailablePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find available port: %w", err)
	}

	proxy := &domain.ProxyConnection{
		ID:           proxyID,
		SessionID:    sessionID,
		Protocol:     protocol,
		LocalPort:    localPort,
		RemoteHost:   remoteHost,
		RemotePort:   remotePort,
		Status:       "created",
		LastActivity: time.Now(),
		CreatedAt:    time.Now(),
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Store in active connections
	internalProxy := &ProxyConnection{
		ID:         proxyID,
		SessionID:  sessionID,
		Protocol:   protocol,
		LocalPort:  localPort,
		RemoteHost: remoteHost,
		RemotePort: remotePort,
		Status:     "created",
	}
	s.activeConnections[proxyID] = internalProxy

	utils.Infof("Created proxy %s for session %s: %s:%d -> %s:%d",
		proxyID, sessionID, "localhost", localPort, remoteHost, remotePort)

	return proxy, nil
}

func (s *proxyService) StartProxy(ctx context.Context, proxyID string) (int, error) {
	s.mu.Lock()
	proxy, exists := s.activeConnections[proxyID]
	if !exists {
		s.mu.Unlock()
		return 0, fmt.Errorf("proxy %s not found", proxyID)
	}
	s.mu.Unlock()

	// Start listening on local port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", proxy.LocalPort))
	if err != nil {
		return 0, fmt.Errorf("failed to start listener: %w", err)
	}

	proxy.listener = listener
	proxy.Status = "active"

	// Create context for this proxy
	proxyCtx, cancel := context.WithCancel(ctx)
	proxy.cancel = cancel

	// Start session recording
	recording, err := s.sessionRecordingService.StartRecording(proxyCtx, proxy.SessionID)
	if err != nil {
		utils.Warnf("Failed to start recording for session %s: %v", proxy.SessionID, err)
	} else {
		utils.Infof("Started recording for session %s: %s", proxy.SessionID, recording.RecordingPath)
	}

	// Start accepting connections
	go s.handleConnections(proxyCtx, proxy)

	utils.Infof("Started proxy %s on port %d", proxyID, proxy.LocalPort)
	return proxy.LocalPort, nil
}

func (s *proxyService) StopProxy(ctx context.Context, proxyID string) error {
	s.mu.Lock()
	proxy, exists := s.activeConnections[proxyID]
	if !exists {
		s.mu.Unlock()
		return fmt.Errorf("proxy %s not found", proxyID)
	}
	delete(s.activeConnections, proxyID)
	s.mu.Unlock()

	// Cancel context and close listener
	if proxy.cancel != nil {
		proxy.cancel()
	}
	if proxy.listener != nil {
		proxy.listener.Close()
	}

	// Close all active connections
	for _, conn := range proxy.connections {
		conn.Close()
	}

	// Stop session recording
	if err := s.sessionRecordingService.StopRecording(ctx, proxy.SessionID); err != nil {
		utils.Warnf("Failed to stop recording for session %s: %v", proxy.SessionID, err)
	}

	utils.Infof("Stopped proxy %s", proxyID)
	return nil
}

func (s *proxyService) GetActiveProxies(ctx context.Context) ([]*domain.ProxyConnection, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	proxies := make([]*domain.ProxyConnection, 0, len(s.activeConnections))
	for _, proxy := range s.activeConnections {
		proxies = append(proxies, &domain.ProxyConnection{
			ID:           proxy.ID,
			SessionID:    proxy.SessionID,
			Protocol:     proxy.Protocol,
			LocalPort:    proxy.LocalPort,
			RemoteHost:   proxy.RemoteHost,
			RemotePort:   proxy.RemotePort,
			Status:       proxy.Status,
			LastActivity: time.Now(),
			CreatedAt:    time.Now(),
		})
	}

	return proxies, nil
}

func (s *proxyService) GetProxyBySession(ctx context.Context, sessionID string) (*domain.ProxyConnection, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, proxy := range s.activeConnections {
		if proxy.SessionID == sessionID {
			return &domain.ProxyConnection{
				ID:           proxy.ID,
				SessionID:    proxy.SessionID,
				Protocol:     proxy.Protocol,
				LocalPort:    proxy.LocalPort,
				RemoteHost:   proxy.RemoteHost,
				RemotePort:   proxy.RemotePort,
				Status:       proxy.Status,
				LastActivity: time.Now(),
				CreatedAt:    time.Now(),
			}, nil
		}
	}

	return nil, fmt.Errorf("no proxy found for session %s", sessionID)
}

func (s *proxyService) UpdateProxyStats(ctx context.Context, proxyID string, bytesIn, bytesOut int64) error {
	// Implementation would update proxy statistics
	return nil
}

func (s *proxyService) handleConnections(ctx context.Context, proxy *ProxyConnection) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := proxy.listener.Accept()
			if err != nil {
				if ctx.Err() != nil {
					return // Context cancelled
				}
				utils.Errorf("Failed to accept connection on proxy %s: %v", proxy.ID, err)
				continue
			}

			// Handle connection based on protocol
			go s.handleConnection(ctx, proxy, conn)
		}
	}
}

func (s *proxyService) handleConnection(ctx context.Context, proxy *ProxyConnection, clientConn net.Conn) {
	defer clientConn.Close()

	// Connect to target server
	targetConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", proxy.RemoteHost, proxy.RemotePort))
	if err != nil {
		utils.Errorf("Failed to connect to target %s:%d: %v", proxy.RemoteHost, proxy.RemotePort, err)
		return
	}
	defer targetConn.Close()

	// Add to active connections
	s.mu.Lock()
	proxy.connections = append(proxy.connections, clientConn, targetConn)
	s.mu.Unlock()

	utils.Infof("Established connection through proxy %s: client -> %s:%d",
		proxy.ID, proxy.RemoteHost, proxy.RemotePort)

	// Handle different protocols
	switch strings.ToLower(proxy.Protocol) {
	case "ssh":
		s.handleSSHConnection(ctx, proxy, clientConn, targetConn)
	case "mysql":
		s.handleMySQLConnection(ctx, proxy, clientConn, targetConn)
	case "postgres", "postgresql":
		s.handlePostgreSQLConnection(ctx, proxy, clientConn, targetConn)
	default:
		// Generic TCP proxy with basic monitoring
		s.handleGenericConnection(ctx, proxy, clientConn, targetConn)
	}
}

func (s *proxyService) handleSSHConnection(ctx context.Context, proxy *ProxyConnection, clientConn, targetConn net.Conn) {
	// Create channels for data flow
	done := make(chan struct{}, 2)

	// Client to Server (commands)
	go func() {
		defer func() { done <- struct{}{} }()
		s.monitorSSHTraffic(ctx, proxy, clientConn, targetConn, "client_to_server")
	}()

	// Server to Client (responses)
	go func() {
		defer func() { done <- struct{}{} }()
		s.monitorSSHTraffic(ctx, proxy, targetConn, clientConn, "server_to_client")
	}()

	// Wait for either direction to close
	<-done
}

func (s *proxyService) handleMySQLConnection(ctx context.Context, proxy *ProxyConnection, clientConn, targetConn net.Conn) {
	// Create channels for data flow
	done := make(chan struct{}, 2)

	// Client to Server (SQL commands)
	go func() {
		defer func() { done <- struct{}{} }()
		s.monitorMySQLTraffic(ctx, proxy, clientConn, targetConn, "client_to_server")
	}()

	// Server to Client (results)
	go func() {
		defer func() { done <- struct{}{} }()
		s.monitorMySQLTraffic(ctx, proxy, targetConn, clientConn, "server_to_client")
	}()

	// Wait for either direction to close
	<-done
}

func (s *proxyService) handlePostgreSQLConnection(ctx context.Context, proxy *ProxyConnection, clientConn, targetConn net.Conn) {
	// Similar to MySQL but for PostgreSQL protocol
	done := make(chan struct{}, 2)

	go func() {
		defer func() { done <- struct{}{} }()
		s.monitorPostgreSQLTraffic(ctx, proxy, clientConn, targetConn, "client_to_server")
	}()

	go func() {
		defer func() { done <- struct{}{} }()
		s.monitorPostgreSQLTraffic(ctx, proxy, targetConn, clientConn, "server_to_client")
	}()

	<-done
}

func (s *proxyService) handleGenericConnection(ctx context.Context, proxy *ProxyConnection, clientConn, targetConn net.Conn) {
	// Simple bidirectional proxy with basic monitoring
	done := make(chan struct{}, 2)

	go func() {
		defer func() { done <- struct{}{} }()
		io.Copy(targetConn, clientConn)
	}()

	go func() {
		defer func() { done <- struct{}{} }()
		io.Copy(clientConn, targetConn)
	}()

	<-done
}

func (s *proxyService) monitorSSHTraffic(ctx context.Context, proxy *ProxyConnection, src, dst net.Conn, direction string) {
	scanner := bufio.NewScanner(src)
	writer := bufio.NewWriter(dst)

	for scanner.Scan() {
		line := scanner.Text()

		// Write to destination
		writer.WriteString(line + "\n")
		writer.Flush()

		// Only analyze client-to-server traffic (commands)
		if direction == "client_to_server" && strings.TrimSpace(line) != "" {
			s.analyzeAndRecordCommand(ctx, proxy, line, "ssh")
		}
	}
}

func (s *proxyService) monitorMySQLTraffic(ctx context.Context, proxy *ProxyConnection, src, dst net.Conn, direction string) {
	// MySQL protocol parsing is complex, but we can intercept text-based queries
	buffer := make([]byte, 4096)

	for {
		n, err := src.Read(buffer)
		if err != nil {
			break
		}

		// Write to destination
		dst.Write(buffer[:n])

		// Analyze SQL commands (simplified)
		if direction == "client_to_server" {
			data := string(buffer[:n])
			if s.containsSQLCommand(data) {
				s.analyzeAndRecordCommand(ctx, proxy, data, "mysql")
			}
		}
	}
}

func (s *proxyService) monitorPostgreSQLTraffic(ctx context.Context, proxy *ProxyConnection, src, dst net.Conn, direction string) {
	// PostgreSQL protocol parsing
	buffer := make([]byte, 4096)

	for {
		n, err := src.Read(buffer)
		if err != nil {
			break
		}

		// Write to destination
		dst.Write(buffer[:n])

		// Analyze SQL commands (simplified)
		if direction == "client_to_server" {
			data := string(buffer[:n])
			if s.containsSQLCommand(data) {
				s.analyzeAndRecordCommand(ctx, proxy, data, "postgresql")
			}
		}
	}
}

func (s *proxyService) analyzeAndRecordCommand(ctx context.Context, proxy *ProxyConnection, command, commandType string) {
	startTime := time.Now()

	// Analyze command for risk
	risk, shouldBlock, err := s.sessionCommandService.AnalyzeCommand(ctx, command, commandType)
	if err != nil {
		utils.Errorf("Failed to analyze command: %v", err)
		risk = "unknown"
	}

	// Create session command record
	sessionCommand := &domain.SessionCommand{
		ID:          uuid.New().String(),
		SessionID:   proxy.SessionID,
		UserID:      proxy.UserID,
		ResourceID:  proxy.ResourceID,
		Command:     command,
		CommandType: commandType,
		Status:      "executed",
		Risk:        risk,
		Timestamp:   startTime,
		Duration:    time.Since(startTime).Milliseconds(),
		CreatedAt:   time.Now(),
	}

	if shouldBlock {
		sessionCommand.Status = "blocked"

		// Create security alert
		alert := &domain.SecurityAlert{
			ID:          uuid.New().String(),
			SessionID:   proxy.SessionID,
			CommandID:   sessionCommand.ID,
			UserID:      proxy.UserID,
			ResourceID:  proxy.ResourceID,
			AlertType:   "blocked_command",
			Severity:    risk,
			Title:       "Blocked High-Risk Command",
			Description: fmt.Sprintf("Command blocked due to %s risk level", risk),
			RawData:     command,
			Action:      "blocked",
			CreatedAt:   time.Now(),
		}

		if err := s.securityAlertService.CreateAlert(ctx, alert); err != nil {
			utils.Errorf("Failed to create security alert: %v", err)
		}
	}

	// Record the command
	if err := s.sessionCommandService.RecordCommand(ctx, sessionCommand); err != nil {
		utils.Errorf("Failed to record command: %v", err)
	}

	utils.Infof("Recorded %s command in session %s: %s (risk: %s)",
		commandType, proxy.SessionID, command[:min(50, len(command))], risk)
}

func (s *proxyService) containsSQLCommand(data string) bool {
	// Simple heuristic to detect SQL commands
	sqlKeywords := []string{"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER"}
	upperData := strings.ToUpper(data)

	for _, keyword := range sqlKeywords {
		if strings.Contains(upperData, keyword) {
			return true
		}
	}
	return false
}

func (s *proxyService) findAvailablePort() (int, error) {
	// Find an available port in the range 10000-20000
	for port := 10000; port < 20000; port++ {
		if s.isPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available ports found")
}

func (s *proxyService) isPortAvailable(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	listener.Close()
	return true
}
