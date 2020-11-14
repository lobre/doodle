package mysql

import (
	"database/sql"
	"errors"

	"github.com/lobre/doodle/pkg/models"
)

type EventModel struct {
	DB *sql.DB
}

func (m *EventModel) Insert(title, desc, time string) (int, error) {
	stmt := `INSERT INTO events (title, description, time)
	VALUES (?, ?, DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, desc, time)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *EventModel) Get(id int) (*models.Event, error) {
	stmt := `SELECT id, title, description, time FROM events
	WHERE time > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	evt := &models.Event{}

	err := row.Scan(&evt.ID, &evt.Title, &evt.Desc, &evt.Time)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return evt, nil
}

func (m *EventModel) Upcoming() ([]*models.Event, error) {
	stmt := `SELECT id, title, description, time FROM events
	WHERE time > UTC_TIMESTAMP() ORDER BY time DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*models.Event{}

	for rows.Next() {
		evt := &models.Event{}

		err = rows.Scan(&evt.ID, &evt.Title, &evt.Desc, &evt.Time)
		if err != nil {
			return nil, err
		}

		events = append(events, evt)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
