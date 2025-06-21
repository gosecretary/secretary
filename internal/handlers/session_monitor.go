package handlers

import (
	"encoding/json"
	"net/http"

	"secretary/alpha/internal/domain"
	"secretary/alpha/pkg/utils"

	"github.com/gorilla/mux"
)

type SessionMonitorHandler struct {
	sessionCommandService   domain.SessionCommandService
	sessionRecordingService domain.SessionRecordingService
	proxyService            domain.ProxyService
	securityAlertService    domain.SecurityAlertService
}

func NewSessionMonitorHandler(
	sessionCommandService domain.SessionCommandService,
	sessionRecordingService domain.SessionRecordingService,
	proxyService domain.ProxyService,
	securityAlertService domain.SecurityAlertService,
) *SessionMonitorHandler {
	return &SessionMonitorHandler{
		sessionCommandService:   sessionCommandService,
		sessionRecordingService: sessionRecordingService,
		proxyService:            proxyService,
		securityAlertService:    securityAlertService,
	}
}

func (h *SessionMonitorHandler) RegisterRoutes(r *mux.Router) {
	// Session Command routes
	r.HandleFunc("/sessions/{session_id}/commands", h.GetSessionCommands).Methods("GET")
	r.HandleFunc("/users/{user_id}/commands", h.GetUserCommands).Methods("GET")
	r.HandleFunc("/resources/{resource_id}/commands", h.GetResourceCommands).Methods("GET")
	r.HandleFunc("/commands/high-risk", h.GetHighRiskCommands).Methods("GET")

	// Session Recording routes
	r.HandleFunc("/sessions/{session_id}/recording/start", h.StartRecording).Methods("POST")
	r.HandleFunc("/sessions/{session_id}/recording/stop", h.StopRecording).Methods("POST")
	r.HandleFunc("/sessions/{session_id}/recording", h.GetRecording).Methods("GET")
	r.HandleFunc("/recordings/{recording_id}/download", h.DownloadRecording).Methods("GET")
	r.HandleFunc("/users/{user_id}/recordings", h.GetUserRecordings).Methods("GET")

	// Proxy routes
	r.HandleFunc("/sessions/{session_id}/proxy", h.CreateProxy).Methods("POST")
	r.HandleFunc("/proxies/{proxy_id}/start", h.StartProxy).Methods("POST")
	r.HandleFunc("/proxies/{proxy_id}/stop", h.StopProxy).Methods("POST")
	r.HandleFunc("/proxies/active", h.GetActiveProxies).Methods("GET")
	r.HandleFunc("/sessions/{session_id}/proxy", h.GetSessionProxy).Methods("GET")

	// Security Alert routes
	r.HandleFunc("/sessions/{session_id}/alerts", h.GetSessionAlerts).Methods("GET")
	r.HandleFunc("/users/{user_id}/alerts", h.GetUserAlerts).Methods("GET")
	r.HandleFunc("/alerts/severity/{severity}", h.GetAlertsBySeverity).Methods("GET")
	r.HandleFunc("/alerts/{alert_id}/review", h.MarkAlertAsReviewed).Methods("POST")
}

// Session Command Handlers

func (h *SessionMonitorHandler) GetSessionCommands(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["session_id"]

	commands, err := h.sessionCommandService.GetSessionCommands(r.Context(), sessionID)
	if err != nil {
		utils.InternalError(w, "Failed to get session commands", err.Error())
		return
	}

	utils.SuccessResponse(w, "Session commands retrieved successfully", commands)
}

func (h *SessionMonitorHandler) GetUserCommands(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	commands, err := h.sessionCommandService.GetCommandsByUser(r.Context(), userID)
	if err != nil {
		utils.InternalError(w, "Failed to get user commands", err.Error())
		return
	}

	utils.SuccessResponse(w, "User commands retrieved successfully", commands)
}

func (h *SessionMonitorHandler) GetResourceCommands(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["resource_id"]

	commands, err := h.sessionCommandService.GetCommandsByResource(r.Context(), resourceID)
	if err != nil {
		utils.InternalError(w, "Failed to get resource commands", err.Error())
		return
	}

	utils.SuccessResponse(w, "Resource commands retrieved successfully", commands)
}

func (h *SessionMonitorHandler) GetHighRiskCommands(w http.ResponseWriter, r *http.Request) {
	commands, err := h.sessionCommandService.GetHighRiskCommands(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to get high risk commands", err.Error())
		return
	}

	utils.SuccessResponse(w, "High risk commands retrieved successfully", commands)
}

// Session Recording Handlers

func (h *SessionMonitorHandler) StartRecording(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["session_id"]

	recording, err := h.sessionRecordingService.StartRecording(r.Context(), sessionID)
	if err != nil {
		utils.InternalError(w, "Failed to start recording", err.Error())
		return
	}

	utils.SuccessResponse(w, "Recording started successfully", recording)
}

func (h *SessionMonitorHandler) StopRecording(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["session_id"]

	err := h.sessionRecordingService.StopRecording(r.Context(), sessionID)
	if err != nil {
		utils.InternalError(w, "Failed to stop recording", err.Error())
		return
	}

	utils.SuccessResponse(w, "Recording stopped successfully", nil)
}

func (h *SessionMonitorHandler) GetRecording(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["session_id"]

	recording, err := h.sessionRecordingService.GetRecording(r.Context(), sessionID)
	if err != nil {
		utils.NotFound(w, "Recording not found")
		return
	}

	utils.SuccessResponse(w, "Recording retrieved successfully", recording)
}

func (h *SessionMonitorHandler) DownloadRecording(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recordingID := vars["recording_id"]

	data, err := h.sessionRecordingService.GetRecordingFile(r.Context(), recordingID)
	if err != nil {
		utils.NotFound(w, "Recording file not found")
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=session_recording_"+recordingID+".txt")
	w.Write(data)
}

func (h *SessionMonitorHandler) GetUserRecordings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	recordings, err := h.sessionRecordingService.ListRecordings(r.Context(), userID)
	if err != nil {
		utils.InternalError(w, "Failed to get user recordings", err.Error())
		return
	}

	utils.SuccessResponse(w, "User recordings retrieved successfully", recordings)
}

// Proxy Handlers

type createProxyRequest struct {
	Protocol   string `json:"protocol"`
	RemoteHost string `json:"remote_host"`
	RemotePort int    `json:"remote_port"`
}

func (h *SessionMonitorHandler) CreateProxy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["session_id"]

	var req createProxyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	proxy, err := h.proxyService.CreateProxy(r.Context(), sessionID, req.Protocol, req.RemoteHost, req.RemotePort)
	if err != nil {
		utils.InternalError(w, "Failed to create proxy", err.Error())
		return
	}

	utils.SuccessResponse(w, "Proxy created successfully", proxy)
}

func (h *SessionMonitorHandler) StartProxy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proxyID := vars["proxy_id"]

	localPort, err := h.proxyService.StartProxy(r.Context(), proxyID)
	if err != nil {
		utils.InternalError(w, "Failed to start proxy", err.Error())
		return
	}

	response := map[string]interface{}{
		"proxy_id":   proxyID,
		"local_port": localPort,
		"status":     "started",
	}

	utils.SuccessResponse(w, "Proxy started successfully", response)
}

func (h *SessionMonitorHandler) StopProxy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proxyID := vars["proxy_id"]

	err := h.proxyService.StopProxy(r.Context(), proxyID)
	if err != nil {
		utils.InternalError(w, "Failed to stop proxy", err.Error())
		return
	}

	utils.SuccessResponse(w, "Proxy stopped successfully", nil)
}

func (h *SessionMonitorHandler) GetActiveProxies(w http.ResponseWriter, r *http.Request) {
	proxies, err := h.proxyService.GetActiveProxies(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to get active proxies", err.Error())
		return
	}

	utils.SuccessResponse(w, "Active proxies retrieved successfully", proxies)
}

func (h *SessionMonitorHandler) GetSessionProxy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["session_id"]

	proxy, err := h.proxyService.GetProxyBySession(r.Context(), sessionID)
	if err != nil {
		utils.NotFound(w, "Proxy not found for session")
		return
	}

	utils.SuccessResponse(w, "Session proxy retrieved successfully", proxy)
}

// Security Alert Handlers

func (h *SessionMonitorHandler) GetSessionAlerts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["session_id"]

	alerts, err := h.securityAlertService.GetAlerts(r.Context(), sessionID)
	if err != nil {
		utils.InternalError(w, "Failed to get session alerts", err.Error())
		return
	}

	utils.SuccessResponse(w, "Session alerts retrieved successfully", alerts)
}

func (h *SessionMonitorHandler) GetUserAlerts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	alerts, err := h.securityAlertService.GetAlertsByUser(r.Context(), userID)
	if err != nil {
		utils.InternalError(w, "Failed to get user alerts", err.Error())
		return
	}

	utils.SuccessResponse(w, "User alerts retrieved successfully", alerts)
}

func (h *SessionMonitorHandler) GetAlertsBySeverity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	severity := vars["severity"]

	alerts, err := h.securityAlertService.GetAlertsBySeverity(r.Context(), severity)
	if err != nil {
		utils.InternalError(w, "Failed to get alerts by severity", err.Error())
		return
	}

	utils.SuccessResponse(w, "Alerts retrieved successfully", alerts)
}

func (h *SessionMonitorHandler) MarkAlertAsReviewed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alertID := vars["alert_id"]

	err := h.securityAlertService.MarkAlertAsReviewed(r.Context(), alertID)
	if err != nil {
		utils.InternalError(w, "Failed to mark alert as reviewed", err.Error())
		return
	}

	utils.SuccessResponse(w, "Alert marked as reviewed successfully", nil)
}
