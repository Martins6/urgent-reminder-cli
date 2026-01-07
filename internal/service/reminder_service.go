package service

import (
	"fmt"
	"time"

	"urgent-reminder/internal/models"
	"urgent-reminder/internal/storage"
)

type ReminderService struct {
	store *storage.JSONStore
}

func NewReminderService(store *storage.JSONStore) *ReminderService {
	return &ReminderService{store: store}
}

func (s *ReminderService) AddReminder(reminder *models.Reminder) error {
	return s.store.AddReminder(reminder)
}

func (s *ReminderService) ListReminders() ([]*models.Reminder, error) {
	return s.store.LoadReminders()
}

func (s *ReminderService) GetDueReminders() ([]*models.Reminder, error) {
	reminders, err := s.store.LoadReminders()
	if err != nil {
		return nil, err
	}

	var dueReminders []*models.Reminder
	for _, r := range reminders {
		if r.IsDue() {
			dueReminders = append(dueReminders, r)
		}
	}

	return dueReminders, nil
}

func (s *ReminderService) GetReminder(id int) (*models.Reminder, error) {
	reminders, err := s.store.LoadReminders()
	if err != nil {
		return nil, err
	}

	for _, r := range reminders {
		if r.ID == id {
			return r, nil
		}
	}

	return nil, fmt.Errorf("reminder with ID %d not found", id)
}

func (s *ReminderService) CheckReminder(id int) error {
	reminder, err := s.GetReminder(id)
	if err != nil {
		return err
	}

	if reminder.IsRecurrent {
		nextDueDate, err := s.calculateNextDueDate(reminder)
		if err != nil {
			return fmt.Errorf("failed to calculate next due date: %w", err)
		}
		reminder.DueDate = nextDueDate
		return s.store.UpdateReminder(id, reminder)
	}

	return s.store.DeleteReminder(id)
}

func (s *ReminderService) calculateNextDueDate(reminder *models.Reminder) (time.Time, error) {
	now := time.Now()

	switch reminder.RecurrentType {
	case models.RecurrentWeekly:
		return s.nextWeeklyOccurrence(reminder, now), nil
	case models.RecurrentBiWeekly:
		return s.nextBiWeeklyOccurrence(reminder, now), nil
	case models.RecurrentMonthly:
		return s.nextMonthlyOccurrence(reminder, now), nil
	default:
		return now.AddDate(1, 0, 0), nil
	}
}

func (s *ReminderService) nextWeeklyOccurrence(reminder *models.Reminder, now time.Time) time.Time {
	if len(reminder.RecurrentDays) == 0 {
		return now.AddDate(0, 0, 7)
	}

	dayMap := map[string]time.Weekday{
		"Mon": time.Monday,
		"Tue": time.Tuesday,
		"Wed": time.Wednesday,
		"Thu": time.Thursday,
		"Fri": time.Friday,
		"Sat": time.Saturday,
		"Sun": time.Sunday,
	}

	var weekdays []time.Weekday
	for _, day := range reminder.RecurrentDays {
		if wd, ok := dayMap[day]; ok {
			weekdays = append(weekdays, wd)
		}
	}

	currentWeekday := now.Weekday()
	minDays := 7

	for _, wd := range weekdays {
		daysUntil := int(wd-currentWeekday+7) % 7
		if daysUntil == 0 {
			daysUntil = 7
		}
		if daysUntil < minDays {
			minDays = daysUntil
		}
	}

	return now.AddDate(0, 0, minDays)
}

func (s *ReminderService) nextBiWeeklyOccurrence(reminder *models.Reminder, now time.Time) time.Time {
	weekly := s.nextWeeklyOccurrence(reminder, now)
	return weekly.AddDate(0, 0, 7)
}

func (s *ReminderService) nextMonthlyOccurrence(reminder *models.Reminder, now time.Time) time.Time {
	day := reminder.RecurrentDayOfMonth
	if day < 1 || day > 31 {
		day = 1
	}

	nextMonth := now.Month()
	nextYear := now.Year()

	for {
		daysInMonth := time.Date(nextYear, nextMonth+1, 0, 0, 0, 0, 0, time.Local).Day()
		if day > daysInMonth {
			day = daysInMonth
		}

		nextDate := time.Date(nextYear, nextMonth, day, 0, 0, 0, 0, time.Local)

		if nextDate.After(now) {
			return nextDate
		}

		nextMonth++
		if nextMonth > 12 {
			nextMonth = 1
			nextYear++
		}
	}
}

func (s *ReminderService) GetDataPath() string {
	return s.store.GetDataPath()
}

func (s *ReminderService) GetNextID() (int, error) {
	return s.store.GetNextID()
}
