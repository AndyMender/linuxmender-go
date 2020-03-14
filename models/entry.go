package models

import (
	"log"
	"strings"
	"time"
	"database/sql"

	"linuxmender/utilities"
)

// Entry is a definition for blog Entry objects
type Entry struct {
	ID			int
	Title		string
	DatePosted	time.Time
	Tags		[]string
}

// EntryManager is a SQL-based manager for Entry records
type EntryManager struct {
	DB	*sql.DB
}

// CreateTable creates an SQL table for storing blog entries
func (mgr *EntryManager) CreateTable() error {
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

// InsertOne inserts a single Entry record into a SQL table
func (mgr *EntryManager) InsertOne(entry *Entry) {
	sqlQuery := "INSERT INTO entries (id, title, date_posted, tags) VALUES (?, ?, ?, ?)"

	// Create a "prepared" SQL statement context
	readyStatement, err := mgr.DB.Prepare(sqlQuery)
	if err != nil {
		log.Println(err)
		return
	}
	defer readyStatement.Close()

	// Execute statement
	_, err = readyStatement.Exec(
		entry.ID, 
		entry.Title, 
		entry.DatePosted.Format(utilities.TimeFormat), 
		strings.Join(entry.Tags, ","),
	)
	if err != nil {
		log.Println(err)
	}

}

// InsertMany is analogous to InsertOne, but accepts a map of Entry records
func (mgr *EntryManager) InsertMany(entries map[string]*Entry) {
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
		_, err = readyStatement.Exec(
			entry.ID, 
			entry.Title, 
			entry.DatePosted.Format(utilities.TimeFormat), 
			strings.Join(entry.Tags, ","),
		)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// GetOne returns a single Entry record from a SQL table
func (mgr *EntryManager) GetOne(entryID int) *Entry {
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
		DatePosted: utilities.IsoToTime(datePosted),
		Tags:       strings.Split(tagsText, ","),
	}
}

// GetAll returns all Entry records from a SQL table
func (mgr *EntryManager) GetAll() map[int]*Entry {
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
			DatePosted: utilities.IsoToTime(datePosted),
			Tags:       strings.Split(tagsText, ","),
		}
	}

	return entryRecords
}

// EntriesByYear re-organizes a map of Entries to group them by year
func EntriesByYear(entryRecords map[int]*Entry) map[int][]*Entry {
	var year int
	var recordsByYear map[int][]*Entry = make(map[int][]*Entry)

	for _, entry := range entryRecords {
		year = entry.DatePosted.Year()

		recordsByYear[year] = append(recordsByYear[year], entry)
	}

	return recordsByYear
}
