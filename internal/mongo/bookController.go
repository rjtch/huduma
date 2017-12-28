package mongo

import (
	"gopkg.in/mgo.v2"
)

type (

	//BookController represents the controller for operating on the book ressources
	BookController struct {
		session *mgo.Session
	}
)

//NewBookController represents the constructor with provided mongo sessio
func NewBookController(s *mgo.Session) *BookController {
	return &BookController{s}
}

//ListAllBooks retrieves all availables books in memory
func (bk *BookController) ListAllBooks([]*Book, error) {
	//books := []*Book{}

}

//FindBookByID retrieves a book in Id in memory
func (bk *BookController) FindBookByID(*Book, error) {

}

//DeleteBookByID deletes a book in Id in memory
func (bk *BookController) DeleteBookByID(*Book, error) {

}

//AddNewBook add a book in  memory
func (bk *BookController) AddNewBook(b *Book) error {
	return nil
}
