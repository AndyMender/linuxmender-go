package models

import (
	"log"
	"time"
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // provides the "sqlite3" driver in the background

	"linuxmender/utilities"
)

// Comment encapsulates a blog entry comment
type Comment struct {
	ID			int
	EntryID		int
	TimePosted	time.Time
	Name		string	`form:"name"`
	Email		string	`form:"email"`
	Comment 	string	`form:"comment"`
}

// CommentManager is a SQL-based manager for comment records
type CommentManager struct {
	DBName string
}

// CreateTable creates an SQL table for storing blog entries
func (mgr *CommentManager) CreateTable() error {
	db, err := sql.Open("sqlite3", mgr.DBName)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	sqlQuery := `
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER NOT NULL PRIMARY KEY,
			entry_id INTEGER NOT NULL,
			time_posted TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			comment TEXT NOT NULL,
			FOREIGN KEY(entry_id) REFERENCES entries(id)
		);
	`

	_, err = db.Exec(sqlQuery)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlQuery)
		return err
	}

	return nil
}

// InsertOne inserts a single Comment record into a SQL table
func (mgr *CommentManager) InsertOne(comment *Comment) {
	db, err := sql.Open("sqlite3", mgr.DBName)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	sqlQuery := `
		INSERT INTO comments (entry_id, time_posted, name, email, comment) 
		VALUES (?, ?, ?, ?, ?)
	`

	// Create a "prepared" SQL statement context
	readyStatement, err := db.Prepare(sqlQuery)
	if err != nil {
		log.Println(err)
		return
	}
	defer readyStatement.Close()

	// Execute statement
	// TODO: sanitize input fields or encode using BASE64?
	_, err = readyStatement.Exec(
		comment.EntryID,
		comment.TimePosted.Format(utilities.DatetimeFormat),
		comment.Name,
		comment.Email,
		comment.Comment,
	)
	if err != nil {
		log.Println(err)
	}
}

// InsertMany is analogous to InsertOne, but accepts a map of Comment records
func (mgr *CommentManager) InsertMany(comments map[string]*Comment) {
	db, err := sql.Open("sqlite3", mgr.DBName)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	sqlQuery := `
		INSERT INTO comments (entry_id, time_posted, name, email, comment) 
		VALUES (?, ?, ?, ?, ?)
	`

	// Create a "prepared" SQL statement context
	readyStatement, err := db.Prepare(sqlQuery)
	if err != nil {
		log.Println(err)
		return
	}
	defer readyStatement.Close()

	// Loop over Comment records and insert them one by one
	for _, comment := range comments {
		_, err = readyStatement.Exec(
			comment.EntryID,
			comment.TimePosted.Format(utilities.DatetimeFormat),
			comment.Name,
			comment.Email,
			comment.Comment,
		)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// GetOne returns a single Comment record from a SQL table
func (mgr *CommentManager) GetOne(commentID int) *Comment {
	db, err := sql.Open("sqlite3", mgr.DBName)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	sqlQuery := `
		SELECT entry_id, time_posted, name, email, comment
		FROM comments
		WHERE id = ?
	`

	// Create a "prepared" SQL statement context
	readyStatement, err := db.Prepare(sqlQuery)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer readyStatement.Close()

	// Fetch Comment record
	var (
		entryID int 
		timePosted, name, email, comment string
	)
	err = readyStatement.QueryRow(commentID).Scan(
		&entryID, 
		&timePosted, 
		&name,
		&email,
		&comment,
	)
	if err != nil {
		log.Println(err)
		return nil
	}

	// Populate Comment record
	return &Comment{
		ID:			commentID,
		EntryID:	entryID,
		TimePosted: utilities.IsoToTime2(timePosted),
		Name: 		name,
		Email: 		email,
		Comment: 	comment,
	}
}

// GetAll returns all Comment records from a SQL table
func (mgr *CommentManager) GetAll() map[int]*Comment {
	commentRecords := make(map[int]*Comment)

	db, err := sql.Open("sqlite3", mgr.DBName)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// Generate a Rows iterator from a SQL query
	sqlQuery := "SELECT id, entry_id, time_posted, name, email, comment FROM comments"
	rows, err := db.Query(sqlQuery)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	// Iterate over rows and populate Entry records
	for rows.Next() {
		var (
			commentID, entryID int
			timePosted, name, email, comment string
		)

		err = rows.Scan(
			&commentID,
			&entryID, 
			&timePosted, 
			&name,
			&email,
			&comment,
		)
		if err != nil {
			log.Println(err)
			return nil
		}

		commentRecords[commentID] = &Comment{
			ID:			commentID,
			EntryID:	entryID,
			TimePosted: utilities.IsoToTime2(timePosted),
			Name: 		name,
			Email: 		email,
			Comment: 	comment,
		}
	}

	return commentRecords
}

// GetByEntry returns all Comment records for a given Entry from a SQL table
func (mgr *CommentManager) GetByEntry(entryID int) map[int]*Comment {
	commentRecords := make(map[int]*Comment)

	db, err := sql.Open("sqlite3", mgr.DBName)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// Generate a Rows iterator from a SQL query
	sqlQuery := `
		SELECT id, entry_id, time_posted, name, email, comment 
		FROM comments
		WHERE entry_id = ?
	`
	// Create a "prepared" SQL statement context
	readyStatement, err := db.Prepare(sqlQuery)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer readyStatement.Close()

	rows, err := readyStatement.Query(entryID)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	// Iterate over rows and populate Entry records
	for rows.Next() {
		var (
			commentID, entryID int
			timePosted, name, email, comment string
		)

		err = rows.Scan(
			&commentID,
			&entryID, 
			&timePosted, 
			&name,
			&email,
			&comment,
		)
		if err != nil {
			log.Println(err)
			return nil
		}

		commentRecords[commentID] = &Comment{
			ID:			commentID,
			EntryID:	entryID,
			TimePosted: utilities.IsoToTime2(timePosted),
			Name: 		name,
			Email: 		email,
			Comment: 	comment,
		}
	}

	return commentRecords
}
