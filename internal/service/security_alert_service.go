package service

import (
	"context"
	"fmt"
	"time"

	"secretary/alpha/internal/domain"
)

type securityAlertService struct {
	alerts map[string]*domain.SecurityAlert
}

func NewSecurityAlertService() domain.SecurityAlertService {
	return &securityAlertService{
		alerts: make(map[string]*domain.SecurityAlert),
	}
}

func (s *securityAlertService) CreateAlert(ctx context.Context, alert *domain.SecurityAlert) error {
	if alert.ID == "" {
		alert.ID = fmt.Sprintf("alert_%d", time.Now().UnixNano())
	}
	if alert.CreatedAt.IsZero() {
		alert.CreatedAt = time.Now()
	}

	s.alerts[alert.ID] = alert
	return nil
}

func (s *securityAlertService) GetAlerts(ctx context.Context, sessionID string) ([]*domain.SecurityAlert, error) {
	var alerts []*domain.SecurityAlert
	for _, alert := range s.alerts {
		if alert.SessionID == sessionID {
			alerts = append(alerts, alert)
		}
	}
	return alerts, nil
}

func (s *securityAlertService) GetAlertsByUser(ctx context.Context, userID string) ([]*domain.SecurityAlert, error) {
	var alerts []*domain.SecurityAlert
	for _, alert := range s.alerts {
		if alert.UserID == userID {
			alerts = append(alerts, alert)
		}
	}
	return alerts, nil
}

func (s *securityAlertService) GetAlertsBySeverity(ctx context.Context, severity string) ([]*domain.SecurityAlert, error) {
	var alerts []*domain.SecurityAlert
	for _, alert := range s.alerts {
		if alert.Severity == severity {
			alerts = append(alerts, alert)
		}
	}
	return alerts, nil
}

func (s *securityAlertService) MarkAlertAsReviewed(ctx context.Context, alertID string) error {
	alert, exists := s.alerts[alertID]
	if !exists {
		return fmt.Errorf("alert %s not found", alertID)
	}

	// In a real implementation, you'd update the alert status
	// For now, just return success
	_ = alert
	return nil
}
