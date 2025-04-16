package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Session struct {
	CurrentUserID uint `json:"currentUserId"`
}

var currentSession *Session

func saveSession() error {
	if currentSession == nil {
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".config", "bodytracker")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	sessionFile := filepath.Join(configDir, "session.json")
	data, err := json.Marshal(currentSession)
	if err != nil {
		return err
	}

	return os.WriteFile(sessionFile, data, 0644)
}

func getCurrentUserID() uint {
	if currentSession == nil {
		return 0
	}
	return currentSession.CurrentUserID
}

func setCurrentUserID(userID uint) error {
	if currentSession == nil {
		currentSession = &Session{}
	}
	currentSession.CurrentUserID = userID
	return saveSession()
}
