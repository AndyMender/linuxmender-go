package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"strings"
)

const dateFormat = "06-01-02"

// Entry is a definition for blog Entry objects
type Entry struct {
	ID                int
	Title, DatePosted string
	Tags              []string
}

// EntryManager is a SQL-based manager for Entry records
type EntryManager struct {
	DB	*sql.DB
}

// CreateEntriesTable creates an SQL table for storing blog entries
func (mgr *EntryManager) CreateEntriesTable() error {
	sqlQuery := `
		CREATE TABLE IF NOT EXISTS entries (
      		id INTEGER NOT NULL PRIMARY KEY,
      		title TEXT NOT NULL,
			date_posted TEXT NOT NULL,
			tags TEXT DEFAULT ''
		);
	`

	_, err := mgr.DB.Exec(sqlQuery)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlQuery)
		return err
	}

	return nil
}

// InsertEntryOne inserts a single Entry record into a SQL table
func (mgr *EntryManager) InsertEntryOne(entry *Entry) {
	sqlQuery := "INSERT INTO entries (id, title, date_posted, tags) VALUES (?, ?, ?, ?)"

	// Create a "prepared" SQL statement context
	readyStatement, err := mgr.DB.Prepare(sqlQuery)
	if err != nil {
		log.Println(err)
		return
	}
	defer readyStatement.Close()

	// Execute statement
	_, err = readyStatement.Exec(entry.ID, entry.Title, entry.DatePosted, strings.Join(entry.Tags, ","))
	if err != nil {
		log.Println(err)
	}

}

// InsertEntryMany is analogous to InsertEntryOne, but accepts a slice of Entry records
func (mgr *EntryManager) InsertEntryMany(entries map[string]Entry) {
	sqlQuery := "INSERT INTO entries (id, title, date_posted, tags) VALUES (?, ?, ?, ?)"

	// Create a "prepared" SQL statement context
	readyStatement, err := mgr.DB.Prepare(sqlQuery)
	if err != nil {
		log.Println(err)
		return
	}
	defer readyStatement.Close()

	// Loop over Entry records and insert them one by one
	for _, entry := range entries {
		_, err = readyStatement.Exec(entry.ID, entry.Title, entry.DatePosted, strings.Join(entry.Tags, ","))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// GetEntryOne returns a single Entry record from a SQL table
func (mgr *EntryManager) GetEntryOne(entryID int) *Entry {
	// Create a "prepared" SQL statement context
	sqlQuery := "SELECT title, date_posted, tags FROM entries WHERE id = ?"
	readyStatement, err := mgr.DB.Prepare(sqlQuery)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer readyStatement.Close()

	// Fetch Entry record
	var title, datePosted, tagsText string
	err = readyStatement.QueryRow(entryID).Scan(&title, &datePosted, &tagsText)
	if err != nil {
		log.Println(err)
		return nil
	}

	// Populate Entry record
	return &Entry{
		ID:         entryID,
		Title:      title,
		DatePosted: datePosted,
		Tags:       strings.Split(tagsText, ","),
	}
}

// GetEntriesAll returns all Entry records from a SQL table
func (mgr *EntryManager) GetEntriesAll() map[int]*Entry {
	entryRecords := make(map[int]*Entry)

	// Generate a Rows iterator from a SQL query
	sqlQuery := "SELECT id, title, date_posted, tags FROM entries"
	rows, err := mgr.DB.Query(sqlQuery)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	// Iterate over rows and populate Entry records
	for rows.Next() {
		var (
			entryID                     int
			title, datePosted, tagsText string
		)

		err = rows.Scan(&entryID, &title, &datePosted, &tagsText)
		if err != nil {
			log.Println(err)
			return nil
		}

		entryRecords[entryID] = &Entry{
			ID:         entryID,
			Title:      title,
			DatePosted: datePosted,
			Tags:       strings.Split(tagsText, ","),
		}
	}

	return entryRecords
}

func parseDate(timeString string) string {
	// Load date string into Time object
	parsedTime, err := time.Parse("2006-01-02", timeString)
	if err != nil {
		log.Printf("Couldn't properly parse date string %v. Returning as is.\n", timeString)
		return timeString
	}

	// Decide on the day of month string ending 
	dayEnding := ""
	switch (parsedTime.Day()) {
	case 1:
		dayEnding = "st"
	case 2:
		dayEnding = "nd"
	case 3:
		dayEnding = "rd"
	default:
		dayEnding = "th"
	}

	return fmt.Sprintf(
		"%s %d%s, %d", 
		parsedTime.Month(), 
		parsedTime.Day(), 
		dayEnding, 
		parsedTime.Year(),
	)
}
