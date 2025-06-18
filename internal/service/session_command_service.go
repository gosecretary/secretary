package service

import (
	"context"
	"regexp"
	"strings"
	"time"

	"secretary/alpha/internal/domain"
	"secretary/alpha/pkg/utils"

	"github.com/google/uuid"
)

type sessionCommandService struct {
	// This would typically have a repository for persistence
	commands map[string][]*domain.SessionCommand
}

func NewSessionCommandService() domain.SessionCommandService {
	return &sessionCommandService{
		commands: make(map[string][]*domain.SessionCommand),
	}
}

func (s *sessionCommandService) RecordCommand(ctx context.Context, command *domain.SessionCommand) error {
	if command.ID == "" {
		command.ID = uuid.New().String()
	}
	if command.CreatedAt.IsZero() {
		command.CreatedAt = time.Now()
	}
	if command.Timestamp.IsZero() {
		command.Timestamp = time.Now()
	}

	// Store in memory (in production, this would go to a database)
	if s.commands[command.SessionID] == nil {
		s.commands[command.SessionID] = make([]*domain.SessionCommand, 0)
	}
	s.commands[command.SessionID] = append(s.commands[command.SessionID], command)

	utils.Infof("Recorded command: %s (type: %s, risk: %s, status: %s)",
		command.Command[:min(100, len(command.Command))],
		command.CommandType, command.Risk, command.Status)

	return nil
}

func (s *sessionCommandService) GetSessionCommands(ctx context.Context, sessionID string) ([]*domain.SessionCommand, error) {
	commands, exists := s.commands[sessionID]
	if !exists {
		return []*domain.SessionCommand{}, nil
	}
	return commands, nil
}

func (s *sessionCommandService) GetCommandsByUser(ctx context.Context, userID string) ([]*domain.SessionCommand, error) {
	var userCommands []*domain.SessionCommand
	for _, sessionCommands := range s.commands {
		for _, cmd := range sessionCommands {
			if cmd.UserID == userID {
				userCommands = append(userCommands, cmd)
			}
		}
	}
	return userCommands, nil
}

func (s *sessionCommandService) GetCommandsByResource(ctx context.Context, resourceID string) ([]*domain.SessionCommand, error) {
	var resourceCommands []*domain.SessionCommand
	for _, sessionCommands := range s.commands {
		for _, cmd := range sessionCommands {
			if cmd.ResourceID == resourceID {
				resourceCommands = append(resourceCommands, cmd)
			}
		}
	}
	return resourceCommands, nil
}

func (s *sessionCommandService) GetHighRiskCommands(ctx context.Context) ([]*domain.SessionCommand, error) {
	var highRiskCommands []*domain.SessionCommand
	for _, sessionCommands := range s.commands {
		for _, cmd := range sessionCommands {
			if cmd.Risk == "high" || cmd.Risk == "critical" {
				highRiskCommands = append(highRiskCommands, cmd)
			}
		}
	}
	return highRiskCommands, nil
}

func (s *sessionCommandService) AnalyzeCommand(ctx context.Context, command string, commandType string) (risk string, shouldBlock bool, err error) {
	command = strings.TrimSpace(command)
	if command == "" {
		return "low", false, nil
	}

	risk = "low"
	shouldBlock = false

	switch strings.ToLower(commandType) {
	case "sql", "mysql", "postgresql", "postgres":
		risk, shouldBlock = s.analyzeSQLCommand(command)
	case "ssh", "shell", "bash":
		risk, shouldBlock = s.analyzeShellCommand(command)
	default:
		risk, shouldBlock = s.analyzeGenericCommand(command)
	}

	return risk, shouldBlock, nil
}

func (s *sessionCommandService) analyzeSQLCommand(command string) (string, bool) {
	upperCommand := strings.ToUpper(command)

	// Critical risk patterns
	criticalPatterns := []string{
		"DROP DATABASE",
		"DROP SCHEMA",
		"TRUNCATE",
		"DELETE FROM.*WHERE.*1=1",
		"UPDATE.*SET.*WHERE.*1=1",
		"SHUTDOWN",
	}

	for _, pattern := range criticalPatterns {
		if matched, _ := regexp.MatchString(pattern, upperCommand); matched {
			return "critical", true // Block critical operations
		}
	}

	// High risk patterns
	highRiskPatterns := []string{
		"DROP TABLE",
		"DROP VIEW",
		"ALTER TABLE.*DROP",
		"DELETE FROM",
		"UPDATE.*SET",
		"GRANT.*ALL",
		"REVOKE",
		"CREATE USER",
		"DROP USER",
		"--.*PASSWORD",
		"UNION.*SELECT",
		"LOAD_FILE",
		"INTO OUTFILE",
	}

	for _, pattern := range highRiskPatterns {
		if matched, _ := regexp.MatchString(pattern, upperCommand); matched {
			return "high", false
		}
	}

	// Medium risk patterns
	mediumRiskPatterns := []string{
		"INSERT INTO",
		"CREATE TABLE",
		"CREATE INDEX",
		"ALTER TABLE",
		"CREATE VIEW",
	}

	for _, pattern := range mediumRiskPatterns {
		if matched, _ := regexp.MatchString(pattern, upperCommand); matched {
			return "medium", false
		}
	}

	// SELECT statements are generally low risk
	if strings.HasPrefix(upperCommand, "SELECT") {
		// But check for potential data exfiltration
		if strings.Contains(upperCommand, "LIMIT") || strings.Contains(upperCommand, "ORDER BY") {
			return "low", false
		}
		// Large result sets might indicate data exfiltration
		if !strings.Contains(upperCommand, "LIMIT") && strings.Contains(upperCommand, "*") {
			return "medium", false
		}
		return "low", false
	}

	return "low", false
}

func (s *sessionCommandService) analyzeShellCommand(command string) (string, bool) {
	command = strings.TrimSpace(command)

	// Critical risk patterns - these should be blocked
	criticalPatterns := []string{
		`rm\s+-rf\s+/`,          // rm -rf /
		`:\(\)\{\s*:\|\:&\s*\}`, // fork bomb
		`mkfs\.`,                // format filesystem
		`dd\s+if=.*of=/dev/`,    // direct disk access
		`chmod\s+777\s+/`,       // dangerous permissions on root
	}

	for _, pattern := range criticalPatterns {
		if matched, _ := regexp.MatchString(pattern, command); matched {
			return "critical", true
		}
	}

	// High risk patterns
	highRiskPatterns := []string{
		`sudo\s+`,
		`su\s+`,
		`rm\s+-rf`,
		`chmod\s+[0-9]+`,
		`chown\s+`,
		`passwd\s+`,
		`useradd\s+`,
		`userdel\s+`,
		`crontab\s+`,
		`systemctl\s+`,
		`service\s+`,
		`iptables\s+`,
		`mount\s+`,
		`umount\s+`,
		`fdisk\s+`,
		`/etc/passwd`,
		`/etc/shadow`,
		`wget\s+.*\|\s*sh`,
		`curl\s+.*\|\s*sh`,
		`nc\s+-l`,            // netcat listener
		`python.*-c.*socket`, // python reverse shell
		`bash\s+-i\s+`,       // interactive bash
	}

	for _, pattern := range highRiskPatterns {
		if matched, _ := regexp.MatchString(pattern, command); matched {
			return "high", false
		}
	}

	// Medium risk patterns
	mediumRiskPatterns := []string{
		`ls\s+-la`,
		`find\s+/`,
		`grep\s+-r`,
		`cat\s+/etc/`,
		`ps\s+aux`,
		`netstat\s+`,
		`ss\s+`,
		`lsof\s+`,
		`history\s*$`,
		`env\s*$`,
		`id\s*$`,
		`whoami\s*$`,
	}

	for _, pattern := range mediumRiskPatterns {
		if matched, _ := regexp.MatchString(pattern, command); matched {
			return "medium", false
		}
	}

	return "low", false
}

func (s *sessionCommandService) analyzeGenericCommand(command string) (string, bool) {
	// Generic analysis for unknown command types
	command = strings.ToLower(command)

	// Look for suspicious patterns
	if strings.Contains(command, "password") ||
		strings.Contains(command, "secret") ||
		strings.Contains(command, "token") {
		return "medium", false
	}

	if strings.Contains(command, "delete") ||
		strings.Contains(command, "remove") ||
		strings.Contains(command, "drop") {
		return "medium", false
	}

	return "low", false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
