package book

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/huduma/internal/mongo"
	"gopkg.in/mgo.v2"

	"github.com/pkg/errors"
)

/*
From package github.com/pkg/errors

- errors.Wrap contructs a stack of erros, adding context to the preceding error. Depending
of the nature of the error it may be necessary to reverse the operation of errors.Wrap to retrieve
the original error for inspection.
*/
const booksCollection = "books"

var (
	//ErrNotFound is abstracting the mongodb not found error
	ErrNotFound = errors.New("Entry not found")

	//ErrInvalidID is used to check if ID is valid
	ErrInvalidID = errors.New("ID not valid")
)

//ListAllBooks is used to retreive all available books in the mongo collection
func ListAllBooks(ctx context.Context, dbCon *mongo.BooksDB) ([]Book, error) {
	b := []Book{}

	f := func(collection *mgo.Collection) error {
		return collection.Find(nil).All(&b)
	}

	if err := dbCon.Execute(ctx, booksCollection, f); err != nil {
		return nil, errors.Wrap(err, "mongo.Books.find()")
	}
	return b, nil
}

//FindBookByID retrieves a book by his ID in mongo
func FindBookByID(ctx context.Context, dbCon *mongo.BooksDB, bookID string) (*Book, error) {

	if !bson.IsObjectIdHex(bookID) {
		return nil, errors.Wrapf(ErrInvalidID, "bookID not valid", bookID)
	}

	queRy := bson.M{"book:id": bookID}

	var b *Book
	f := func(collection *mgo.Collection) error {
		return collection.Find(queRy).One(&b)
	}

	if err := dbCon.Execute(ctx, booksCollection, f); err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, fmt.Sprintf("db connection failed when finding book(%s) ", dbCon.Query(queRy)))
	}
	return b, nil
}

//AddNewBook creates a new book into mongo
func AddNewBook(ctx context.Context, dbCon *mongo.BooksDB, bk *Book) (*Book, error) {
	now := time.Now().UTC()

	b := Book{
		BookID:       bson.NewObjectId().Hex(),
		BookType:     bk.BookType,
		Title:        bk.Title,
		Authors:      make([]Author, len(bk.Authors)),
		Price:        bk.Price,
		DateModified: bk.DateModified,
		DateCreated:  bk.DateCreated,
	}

	bk.DateCreated = &now

	f := func(collection *mgo.Collection) error {
		return collection.Insert(&b)
	}

	if err := dbCon.Execute(ctx, booksCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db connection failed when insering book(%s)", dbCon.Query(&b)))
	}
	return &b, nil
}

//UpdateBookbyID updates an existing book in mongo
func UpdateBookbyID(ctx context.Context, dbCon *mongo.BooksDB, bookID string, cb *Book) error {

	if !bson.IsObjectIdHex(bookID) {
		return errors.Wrapf(ErrInvalidID, "bookID might be wrong", bookID)
	}

	now := time.Now().UTC()
	cb.DateModified = &now

	f := func(collection *mgo.Collection) error {
		return collection.Update(bson.M{"book_id": bookID}, &cb)
	}

	if err := dbCon.Execute(ctx, booksCollection, f); err != nil {
		if err == mgo.ErrNotFound {
			return mgo.ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("db connection failed when updating book(%s)", dbCon.Query(bson.M{"book_id": bookID})))
	}
	return nil
}

//DeleteBookByID delete an existing book from mongo using it ID
func DeleteBookByID(ctx context.Context, dbCon *mongo.BooksDB, bookID string) error {

	if !bson.IsObjectIdHex(bookID) {
		return errors.Wrapf(ErrInvalidID, "bookID might be wrong", bookID)
	}

	q := bson.M{"book_id": bookID}

	f := func(collection *mgo.Collection) error {
		return collection.Remove(q)
	}

	if err := dbCon.Execute(ctx, booksCollection, f); err != nil {
		if err == mgo.ErrNotFound {
			return mgo.ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("Cannot remove requested bookID(%s)", dbCon.Query(q)))
	}
	return nil
}
