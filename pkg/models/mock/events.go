package mock

import (
	"time"

	"github.com/lobre/doodle/pkg/models"
)

var mockEvent = &models.Event{
	ID:    1,
	Title: "Music festival",
	Desc:  "Happening every year, and always fun.",
	Time:  time.Now(),
}

type EventStore struct{}

func (m *EventStore) Insert(title, desc, time string) (int, error) {
	return 2, nil
}

func (m *EventStore) Get(id int) (*models.Event, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetStore) Upcoming() ([]*models.Event, error) {
	return []*models.Event{mockEvent}, nil
}
