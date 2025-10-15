package models

import (
	"time"

	"example.com/rest-api/db"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"` // Gin uses "binding:" for validation
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int
}

// save event to database
func (e *Event) Save() error {
	query := `
		INSERT INTO events(name, description, location, dateTime, user_id)
		VALUES (?, ?, ?, ?, ?) 
	`
	// create reusable SQL query
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	// close when done
	defer stmt.Close()

	// run prepared statement with specific parameters, use .Exec for updating/changing data
	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	e.ID = id
	return err
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query) // use .Query for getting data
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event

	// for each row, scan the columns into event struct fields
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (event Event) Update() error {
	query := `
		UPDATE events 
		SET name = ?, description = ?, location = ?, dateTime = ?
		WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	
	defer stmt.Close()

	_, err = stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.ID)
	return err
}