package models

import (
	"database/sql"
	"log"
	"strings"
)

// Entry is a definition for Blog entry objects
type Entry struct {
	ID                int
	Title, DatePosted string
	Tags              []string
}

// CreateEntriesTable creates an SQL table for storing blog entries
func CreateEntriesTable(db *sql.DB) error {
	sqlQuery := `
		CREATE TABLE IF NOT EXISTS entries (
      id INTEGER NOT NULL PRIMARY KEY,
      title TEXT NOT NULL,
			date_posted TEXT NOT NULL,
			tags TEXT DEFAULT ''
		);
	`

	_, err := db.Exec(sqlQuery)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlQuery)
		return err
	}

	return nil
}

// InsertEntryOne inserts a single Entry record into a SQL table
func InsertEntryOne(entry *Entry, db *sql.DB) {
	sqlQuery := "INSERT INTO entries (id, title, date_posted, tags) VALUES (?, ?, ?, ?)"

	// Create a "prepared" SQL statement context
	readyStatement, err := db.Prepare(sqlQuery)
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
func InsertEntryMany(entries map[string]Entry, db *sql.DB) {
	sqlQuery := "INSERT INTO entries (id, title, date_posted, tags) VALUES (?, ?, ?, ?)"

	// Create a "prepared" SQL statement context
	readyStatement, err := db.Prepare(sqlQuery)
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
func GetEntryOne(entryID int, db *sql.DB) *Entry {
	// Create a "prepared" SQL statement context
	sqlQuery := "SELECT title, date_posted, tags FROM entries WHERE id = ?"
	readyStatement, err := db.Prepare(sqlQuery)
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
func GetEntriesAll(db *sql.DB) map[int]*Entry {
	entryRecords := make(map[int]*Entry)

	// Generate a Rows iterator from a SQL query
	sqlQuery := "SELECT id, title, date_posted, tags FROM entries"
	rows, err := db.Query(sqlQuery)
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
