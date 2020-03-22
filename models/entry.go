package models

import (
	"log"
	"strings"
	"time"
	"database/sql"
	_ "github.com/lib/pq" // provides the "pq" driver in the background
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
	ConnStr string
}

// CreateTable creates an SQL table for storing blog entries
func (mgr *EntryManager) CreateTable() error {
	db, err := sql.Open("postgres", mgr.ConnStr)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	queryStr := `
		CREATE TABLE IF NOT EXISTS entries (
      		id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			date_posted DATE NOT NULL,
			tags TEXT DEFAULT ''
		);
	`

	_, err = db.Exec(queryStr)
	if err != nil {
		log.Printf("%q: %s\n", err, queryStr)
		return err
	}

	return nil
}

// InsertOne inserts a single Entry record into a SQL table
func (mgr *EntryManager) InsertOne(entry *Entry) {
	db, err := sql.Open("postgres", mgr.ConnStr)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	queryStr := `
		INSERT INTO entries (id, title, date_posted, tags) 
		VALUES (?, ?, ?, ?)
	`

	// Create a "prepared" SQL statement context
	stmt, err := db.Prepare(queryStr)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	// Execute statement
	_, err = stmt.Exec(
		entry.ID, 
		entry.Title, 
		entry.DatePosted, 
		strings.Join(entry.Tags, ","),
	)
	if err != nil {
		log.Println(err)
	}
}

// InsertMany is analogous to InsertOne, but accepts a map of Entry records
func (mgr *EntryManager) InsertMany(entries map[int]*Entry) {
	db, err := sql.Open("postgres", mgr.ConnStr)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	queryStr := `
		INSERT INTO entries (id, title, date_posted, tags) 
		VALUES (?, ?, ?, ?)
	`

	// Create a "prepared" SQL statement context
	stmt, err := db.Prepare(queryStr)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	// Loop over Entry records and insert them one by one
	for _, entry := range entries {
		_, err = stmt.Exec(
			entry.ID, 
			entry.Title, 
			entry.DatePosted, 
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
	db, err := sql.Open("postgres", mgr.ConnStr)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	queryStr := `
		SELECT title, date_posted, tags 
		FROM entries 
		WHERE id = ?
	`

	// Create a "prepared" SQL statement context
	stmt, err := db.Prepare(queryStr)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer stmt.Close()

	// Fetch Entry record
	var (
		title, tagsText string
		datePosted 		time.Time
	)
	err = stmt.QueryRow(entryID).Scan(&title, &datePosted, &tagsText)
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

// GetAll returns all Entry records from a SQL table
func (mgr *EntryManager) GetAll() map[int]*Entry {
	entries := make(map[int]*Entry)

	db, err := sql.Open("postgres", mgr.ConnStr)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// Generate a Rows iterator from a SQL query
	queryStr := "SELECT id, title, date_posted, tags FROM entries"
	rows, err := db.Query(queryStr)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	// Iterate over rows and populate Entry records
	for rows.Next() {
		var (
			entryID         int
			title, tagsText string
			datePosted 		time.Time
		)

		err = rows.Scan(&entryID, &title, &datePosted, &tagsText)
		if err != nil {
			log.Println(err)
			return nil
		}

		entries[entryID] = &Entry{
			ID:         entryID,
			Title:      title,
			DatePosted: datePosted,
			Tags:       strings.Split(tagsText, ","),
		}
	}

	return entries
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
