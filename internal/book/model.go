package book

import (
	"time"
)

type (

	//Book represents the structure of the ressource
	Book struct {
		BookID       string     `json:"book_id" bson:"book_id"`
		BookType     int        `json:"type" bson:"type"`
		Title        string     `json:"title" bson:"title"`
		Authors      []Author   `json:"authors"  bson:"authors"`
		Price        string     `json:"price" bson:"price"`
		DateModified *time.Time `json:"date_modified" bson:"date_modified"`
		DateCreated  *time.Time `json:"date_created" bson:"date_created"`
	}
)

//Author is the structure of a Book's author
type Author struct {
	Firstname string `json:"firstname" bson:"firstname"`
	Lastname  string `json:"lastname" bson:"lastname"`
}
