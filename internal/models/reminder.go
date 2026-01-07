package models

import (
	"time"
)

type RecurrentType string

const (
	RecurrentNone     RecurrentType = "none"
	RecurrentWeekly   RecurrentType = "weekly"
	RecurrentBiWeekly RecurrentType = "bi-weekly"
	RecurrentMonthly  RecurrentType = "monthly"
)

type Reminder struct {
	ID                  int           `json:"id"`
	Title               string        `json:"title"`
	DueDate             time.Time     `json:"due_date"`
	Time                string        `json:"time,omitempty"`
	IsRecurrent         bool          `json:"is_recurrent"`
	RecurrentType       RecurrentType `json:"recurrent_type,omitempty"`
	RecurrentDays       []string      `json:"recurrent_days,omitempty"`
	RecurrentDayOfMonth int           `json:"recurrent_day_of_month,omitempty"`
	CreatedAt           time.Time     `json:"created_at"`
}

func NewReminder(id int, title string, dueDate time.Time) *Reminder {
	return &Reminder{
		ID:          id,
		Title:       title,
		DueDate:     dueDate,
		IsRecurrent: false,
		CreatedAt:   time.Now(),
	}
}

func NewRecurrentReminder(id int, title string, dueDate time.Time, recurrentType RecurrentType) *Reminder {
	return &Reminder{
		ID:            id,
		Title:         title,
		DueDate:       dueDate,
		IsRecurrent:   true,
		RecurrentType: recurrentType,
		CreatedAt:     time.Now(),
	}
}

func (r *Reminder) IsOverdue() bool {
	now := time.Now()
	dueDateTime := r.DueDate

	if r.Time != "" {
		parsedTime, _ := time.Parse("15:04", r.Time)
		dueDateTime = time.Date(r.DueDate.Year(), r.DueDate.Month(), r.DueDate.Day(),
			parsedTime.Hour(), parsedTime.Minute(), 0, 0, r.DueDate.Location())
	}

	return now.After(dueDateTime)
}

func (r *Reminder) IsDue() bool {
	now := time.Now()
	dueDateTime := r.DueDate

	if r.Time != "" {
		parsedTime, _ := time.Parse("15:04", r.Time)
		dueDateTime = time.Date(r.DueDate.Year(), r.DueDate.Month(), r.DueDate.Day(),
			parsedTime.Hour(), parsedTime.Minute(), 0, 0, r.DueDate.Location())
	}

	return now.After(dueDateTime) || now.Equal(dueDateTime)
}

func (r *Reminder) FormatDueDate() string {
	return r.DueDate.Format("2006-01-02")
}

func (r *Reminder) FormatTime() string {
	if r.Time == "" {
		return ""
	}
	return r.Time
}
