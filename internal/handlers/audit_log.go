package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"secretary/alpha/internal/domain"
	"secretary/alpha/pkg/utils"

	"github.com/gorilla/mux"
)

type AuditLogHandler struct {
	auditLogService domain.AuditLogService
}

func NewAuditLogHandler(auditLogService domain.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{
		auditLogService: auditLogService,
	}
}

func (h *AuditLogHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/audit-logs", h.List).Methods("GET")
	r.HandleFunc("/api/audit-logs/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/api/audit-logs/user/{userID}", h.GetByUserID).Methods("GET")
	r.HandleFunc("/api/audit-logs/resource/{resourceID}", h.GetByResourceID).Methods("GET")
	r.HandleFunc("/api/audit-logs/action/{action}", h.GetByAction).Methods("GET")
	r.HandleFunc("/api/audit-logs/date-range", h.GetByDateRange).Methods("GET")
}

type dateRangeRequest struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

func (h *AuditLogHandler) List(w http.ResponseWriter, r *http.Request) {
	logs, err := h.auditLogService.List(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to list audit logs", err.Error())
		return
	}

	utils.SuccessResponse(w, "Audit logs retrieved successfully", logs)
}

func (h *AuditLogHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log, err := h.auditLogService.GetByID(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "Audit log not found")
		return
	}

	utils.SuccessResponse(w, "Audit log retrieved successfully", log)
}

func (h *AuditLogHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	logs, err := h.auditLogService.GetByUserID(r.Context(), userID)
	if err != nil {
		utils.InternalError(w, "Failed to get audit logs by user ID", err.Error())
		return
	}

	utils.SuccessResponse(w, "Audit logs retrieved successfully", logs)
}

func (h *AuditLogHandler) GetByResourceID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["resourceID"]

	logs, err := h.auditLogService.GetByResourceID(r.Context(), resourceID)
	if err != nil {
		utils.InternalError(w, "Failed to get audit logs by resource ID", err.Error())
		return
	}

	utils.SuccessResponse(w, "Audit logs retrieved successfully", logs)
}

func (h *AuditLogHandler) GetByAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]

	logs, err := h.auditLogService.GetByAction(r.Context(), action)
	if err != nil {
		utils.InternalError(w, "Failed to get audit logs by action", err.Error())
		return
	}

	utils.SuccessResponse(w, "Audit logs retrieved successfully", logs)
}

func (h *AuditLogHandler) GetByDateRange(w http.ResponseWriter, r *http.Request) {
	var req dateRangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	logs, err := h.auditLogService.GetByDateRange(r.Context(), req.StartDate, req.EndDate)
	if err != nil {
		utils.InternalError(w, "Failed to get audit logs by date range", err.Error())
		return
	}

	utils.SuccessResponse(w, "Audit logs retrieved successfully", logs)
}
