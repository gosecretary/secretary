package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"secretary/alpha/internal/domain"
	"secretary/alpha/pkg/utils"
)

type sessionRecordingService struct {
	recordings map[string]*domain.SessionRecording
	basePath   string
}

func NewSessionRecordingService() domain.SessionRecordingService {
	basePath := "./data/recordings"
	if err := os.MkdirAll(basePath, 0755); err != nil {
		utils.Errorf("Failed to create recordings directory: %v", err)
	}

	return &sessionRecordingService{
		recordings: make(map[string]*domain.SessionRecording),
		basePath:   basePath,
	}
}

func (s *sessionRecordingService) StartRecording(ctx context.Context, sessionID string) (*domain.SessionRecording, error) {
	recordingID := uuid.New().String()
	recordingPath := filepath.Join(s.basePath, fmt.Sprintf("session_%s_%s.txt", sessionID, recordingID))

	recording := &domain.SessionRecording{
		ID:            recordingID,
		SessionID:     sessionID,
		RecordingPath: recordingPath,
		Format:        "text",
		Size:          0,
		Duration:      0,
		CommandCount:  0,
		CreatedAt:     time.Now(),
	}

	s.recordings[recordingID] = recording

	// Create the recording file
	file, err := os.Create(recordingPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create recording file: %w", err)
	}
	defer file.Close()

	utils.Infof("Started recording for session %s: %s", sessionID, recordingPath)
	return recording, nil
}

func (s *sessionRecordingService) StopRecording(ctx context.Context, sessionID string) error {
	// Find recording by session ID
	var recording *domain.SessionRecording
	for _, r := range s.recordings {
		if r.SessionID == sessionID {
			recording = r
			break
		}
	}

	if recording == nil {
		return fmt.Errorf("no recording found for session %s", sessionID)
	}

	recording.Duration = int64(time.Since(recording.CreatedAt).Seconds())

	// Update file size
	if info, err := os.Stat(recording.RecordingPath); err == nil {
		recording.Size = info.Size()
	}

	utils.Infof("Stopped recording for session %s", sessionID)
	return nil
}

func (s *sessionRecordingService) GetRecording(ctx context.Context, sessionID string) (*domain.SessionRecording, error) {
	for _, recording := range s.recordings {
		if recording.SessionID == sessionID {
			return recording, nil
		}
	}
	return nil, fmt.Errorf("recording not found for session %s", sessionID)
}

func (s *sessionRecordingService) GetRecordingFile(ctx context.Context, recordingID string) ([]byte, error) {
	recording, exists := s.recordings[recordingID]
	if !exists {
		return nil, fmt.Errorf("recording %s not found", recordingID)
	}

	return os.ReadFile(recording.RecordingPath)
}

func (s *sessionRecordingService) DeleteRecording(ctx context.Context, recordingID string) error {
	recording, exists := s.recordings[recordingID]
	if !exists {
		return fmt.Errorf("recording %s not found", recordingID)
	}

	// Delete file
	if err := os.Remove(recording.RecordingPath); err != nil {
		return fmt.Errorf("failed to delete recording file: %w", err)
	}

	// Remove from memory
	delete(s.recordings, recordingID)

	return nil
}

func (s *sessionRecordingService) ListRecordings(ctx context.Context, userID string) ([]*domain.SessionRecording, error) {
	var recordings []*domain.SessionRecording
	for _, recording := range s.recordings {
		// In a real implementation, you'd filter by userID
		// For now, return all recordings
		recordings = append(recordings, recording)
	}
	return recordings, nil
}
