package models

// Comment encapsulates a blog entry comment
type Comment struct {
	EntryID	int
	Name	string	`form:"name"`
	Email	string	`form:"email"`
	Comment string	`form:"comment"`
}
