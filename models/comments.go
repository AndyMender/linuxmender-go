package models

import (
	"log"
	"time"
	"encoding/base64"
	"database/sql"

	_ "github.com/lib/pq" // provides the "pq" driver in the background
)

// Comment encapsulates a blog entry comment
type Comment struct {
	ID			int
	EntryID		int
	TimePosted	time.Time
	Name		string	`form:"name"`
	Email		string	`form:"email"`
	Text 		string	`form:"comment"`
}

// CommentManager is a SQL-based manager for comment records
type CommentManager struct {
	ConnStr string
}

// CreateTable creates an SQL table for storing blog entries
func (mgr *CommentManager) CreateTable() error {
	db, err := sql.Open("postgres", mgr.ConnStr)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	queryStr := `
		CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			entry_id INTEGER NOT NULL,
			time_posted TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			comment TEXT NOT NULL,
			FOREIGN KEY(entry_id) REFERENCES entries(id)
		);
	`

	_, err = db.Exec(queryStr)
	if err != nil {
		log.Printf("%q: %s\n", err, queryStr)
		return err
	}

	return nil
}

// InsertOne inserts a single Comment record into a SQL table
func (mgr *CommentManager) InsertOne(comment *Comment) {
	db, err := sql.Open("postgres", mgr.ConnStr)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	queryStr := `
		INSERT INTO comments (entry_id, time_posted, name, email, comment) 
		VALUES ($1, $2, $3, $4, $5)
	`

	// Create a "prepared" SQL statement context
	stmt, err := db.Prepare(queryStr)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	// Execute statement
	// TODO: sanitize input fields or encode using BASE64?
	_, err = stmt.Exec(
		comment.EntryID,
		comment.TimePosted,
		comment.Name,
		comment.Email,
		base64.StdEncoding.EncodeToString([]byte(comment.Text)),
	)
	if err != nil {
		log.Println(err)
	}
}

// InsertMany is analogous to InsertOne, but accepts a map of Comment records
func (mgr *CommentManager) InsertMany(comments map[int]*Comment) {
	db, err := sql.Open("postgres", mgr.ConnStr)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	queryStr := `
		INSERT INTO comments (entry_id, time_posted, name, email, comment) 
		VALUES ($1, $2, $3, $4, $5)
	`

	// Create a "prepared" SQL statement context
	stmt, err := db.Prepare(queryStr)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	// Loop over Comment records and insert them one by one
	for _, comment := range comments {
		_, err = stmt.Exec(
			comment.EntryID,
			comment.TimePosted,
			comment.Name,
			comment.Email,
			base64.StdEncoding.EncodeToString([]byte(comment.Text)),
		)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// GetOne returns a single Comment record from a SQL table
func (mgr *CommentManager) GetOne(commentID int) *Comment {
	db, err := sql.Open("postgres", mgr.ConnStr)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	queryStr := `
		SELECT entry_id, time_posted, name, email, comment
		FROM comments
		WHERE id = $1
	`

	// Create a "prepared" SQL statement context
	stmt, err := db.Prepare(queryStr)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer stmt.Close()

	// Fetch Comment record
	var (
		entryID 			 	int 
		name, email, comment 	string
		timePosted			 	time.Time
	)
	err = stmt.QueryRow(commentID).Scan(
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

	// Decode Comment text
	commentBytes, err := base64.StdEncoding.DecodeString(comment)
	if err != nil {
		log.Println(err)
		return nil
	}

	// Populate Comment record
	return &Comment{
		ID:			commentID,
		EntryID:	entryID,
		TimePosted: timePosted,
		Name: 		name,
		Email: 		email,
		Text: 		string(commentBytes),
	}
}

// GetAll returns all Comment records from a SQL table
func (mgr *CommentManager) GetAll() map[int]*Comment {
	comments := make(map[int]*Comment)

	db, err := sql.Open("postgres", mgr.ConnStr)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// Generate a Rows iterator from a SQL query
	queryStr := "SELECT id, entry_id, time_posted, name, email, comment FROM comments"
	rows, err := db.Query(queryStr)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	// Iterate over rows and populate Entry records
	for rows.Next() {
		var (
			commentID, entryID 		int
			name, email, comment 	string
			timePosted 				time.Time
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

		// Decode Comment text
		commentBytes, err := base64.StdEncoding.DecodeString(comment)
		if err != nil {
			log.Println(err)
			return nil
		}

		comments[commentID] = &Comment{
			ID:			commentID,
			EntryID:	entryID,
			TimePosted: timePosted,
			Name: 		name,
			Email: 		email,
			Text: 		string(commentBytes),
		}
	}

	return comments
}

// GetByEntry returns all Comment records for a given Entry from a SQL table
func (mgr *CommentManager) GetByEntry(entryID int) map[int]*Comment {
	comments := make(map[int]*Comment)

	db, err := sql.Open("postgres", mgr.ConnStr)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// Generate a Rows iterator from a SQL query
	queryStr := `
		SELECT id, entry_id, time_posted, name, email, comment 
		FROM comments
		WHERE entry_id = $1
	`
	// Create a "prepared" SQL statement context
	stmt, err := db.Prepare(queryStr)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer stmt.Close()

	rows, err := stmt.Query(entryID)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	// Iterate over rows and populate Entry records
	for rows.Next() {
		var (
			commentID, entryID 		int
			name, email, comment 	string
			timePosted 				time.Time
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

		// Decode Comment text
		commentBytes, err := base64.StdEncoding.DecodeString(comment)
		if err != nil {
			log.Println(err)
			return nil
		}

		// TODO: use `utilities.HumanizeTime` to improve comment timestamp?
		comments[commentID] = &Comment{
			ID:			commentID,
			EntryID:	entryID,
			TimePosted: timePosted,
			Name: 		name,
			Email: 		email,
			Text: 		string(commentBytes),
		}
	}

	return comments
}
