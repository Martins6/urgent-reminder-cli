package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"urgent-reminder/internal/models"
)

const (
	appName       = "urgent-reminder"
	remindersFile = "reminders.json"
)

type JSONStore struct {
	dataPath string
}

func NewJSONStore() (*JSONStore, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	dataHome := filepath.Join(homeDir, ".local", "share")
	appDataPath := filepath.Join(dataHome, appName)

	if err := os.MkdirAll(appDataPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	return &JSONStore{
		dataPath: filepath.Join(appDataPath, remindersFile),
	}, nil
}

func (s *JSONStore) GetDataPath() string {
	return s.dataPath
}

func (s *JSONStore) LoadReminders() ([]*models.Reminder, error) {
	data, err := os.ReadFile(s.dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*models.Reminder{}, nil
		}
		return nil, fmt.Errorf("failed to read reminders file: %w", err)
	}

	var reminders []*models.Reminder
	if err := json.Unmarshal(data, &reminders); err != nil {
		if strings.Contains(err.Error(), "id") {
			reminders, err = s.migrateOldFormat(data)
			if err != nil {
				return nil, fmt.Errorf("failed to migrate old format: %w", err)
			}
			if err := s.SaveReminders(reminders); err != nil {
				return nil, fmt.Errorf("failed to save migrated reminders: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to parse reminders: %w", err)
		}
	}

	return reminders, nil
}

type OldReminder struct {
	ID           string    `json:"id"`
	Description  string    `json:"description"`
	DueDate      time.Time `json:"due_date"`
	AlertEnabled bool      `json:"alert_enabled"`
	CreatedAt    time.Time `json:"created_at"`
}

func (s *JSONStore) migrateOldFormat(data []byte) ([]*models.Reminder, error) {
	var oldReminders []OldReminder
	if err := json.Unmarshal(data, &oldReminders); err != nil {
		return nil, err
	}

	var reminders []*models.Reminder
	nextID := 1
	for _, old := range oldReminders {
		reminder := &models.Reminder{
			ID:        nextID,
			Title:     old.Description,
			DueDate:   old.DueDate,
			CreatedAt: old.CreatedAt,
		}
		reminders = append(reminders, reminder)
		nextID++
	}

	return reminders, nil
}

func (s *JSONStore) SaveReminders(reminders []*models.Reminder) error {
	data, err := json.MarshalIndent(reminders, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal reminders: %w", err)
	}

	if err := os.WriteFile(s.dataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write reminders file: %w", err)
	}

	return nil
}

func (s *JSONStore) AddReminder(reminder *models.Reminder) error {
	reminders, err := s.LoadReminders()
	if err != nil {
		return err
	}

	reminders = append(reminders, reminder)
	return s.SaveReminders(reminders)
}

func (s *JSONStore) UpdateReminder(id int, updatedReminder *models.Reminder) error {
	reminders, err := s.LoadReminders()
	if err != nil {
		return err
	}

	found := false
	for i, r := range reminders {
		if r.ID == id {
			reminders[i] = updatedReminder
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("reminder with ID %d not found", id)
	}

	return s.SaveReminders(reminders)
}

func (s *JSONStore) DeleteReminder(id int) error {
	reminders, err := s.LoadReminders()
	if err != nil {
		return err
	}

	var updatedReminders []*models.Reminder
	found := false
	for _, r := range reminders {
		if r.ID == id {
			found = true
			continue
		}
		updatedReminders = append(updatedReminders, r)
	}

	if !found {
		return fmt.Errorf("reminder with ID %d not found", id)
	}

	return s.SaveReminders(updatedReminders)
}

func (s *JSONStore) GetNextID() (int, error) {
	reminders, err := s.LoadReminders()
	if err != nil {
		return 0, err
	}

	if len(reminders) == 0 {
		return 1, nil
	}

	maxID := 0
	for _, r := range reminders {
		if r.ID > maxID {
			maxID = r.ID
		}
	}

	return maxID + 1, nil
}
