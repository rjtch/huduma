package handlers

import (
	"context"
	"net/http"

	"github.com/huduma/internal/book"
	"github.com/huduma/internal/mongo"
	"github.com/huduma/internal/router"
	"github.com/pkg/errors"
)

//checkError looks for errors type that could occur and transform them into
//router errors
func checkError(err error) error {

	switch errors.Cause(err) {
	case book.ErrNotFound:
		return router.ErrNotFound
	case book.ErrInvalidID:
		return router.ErrInvalidID
	}
	return err
}

//Books represents the Book API method handler set.
type Books struct {
	db *mongo.BooksDB
}

//Read lists all existing books
func (b *Books) Read(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	dbCon, err := b.db.Copy()
	if err != nil {
		return errors.Wrap(router.ErrDBNotConfigured, "")
	}
	defer dbCon.Close()

	bks, err := book.ListAllBooks(ctx, dbCon)
	if err = checkError(err); err != nil {
		return errors.Wrap(err, "")
	}
	router.Response(ctx, w, bks, http.StatusOK)
	return nil
}

//retreive returns a given specified book from internal
func (b *Books) retreive(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	dbCon, err := b.db.Copy()
	if err != nil {
		return errors.Wrap(router.ErrDBNotConfigured, "")
	}
	defer dbCon.Close()

	bk, err := book.FindBookByID(ctx, dbCon, params["id"])
	if err = checkError(err); err != nil {
		return errors.Wrap(err, "")
	}
	router.Response(ctx, w, bk, http.StatusOK)
	return nil
}

//Create adds a new book to the system
func (b *Books) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	dbCon, err := b.db.Copy()
	if err != nil {
		return errors.Wrap(router.ErrDBNotConfigured, "")
	}
	defer dbCon.Close()

	var newBook book.Book
	if err := router.Unmarshall(r.Body, &newBook); err != nil {
		return errors.Wrap(err, "")
	}

	nbk, err := book.AddNewBook(ctx, dbCon, &newBook)
	if err = checkError(err); err != nil {
		return errors.Wrapf(err, "Book: %+v", &newBook)
	}

	router.Response(ctx, w, nbk, http.StatusOK)
	return nil
}

//Delete deletes a specified book from the system
func (b *Books) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	dbCon, err := b.db.Copy()
	if err != nil {
		return errors.Wrap(router.ErrDBNotConfigured, "")
	}

	defer dbCon.Close()

	err = book.DeleteBookByID(ctx, dbCon, params["id"])
	if err = checkError(err); err != nil {
		return errors.Wrapf(err, "ID: %s", params["id"])
	}
	router.Response(ctx, w, nil, http.StatusNoContent)
	return nil
}

//Update updates an existing book in the system
func (b *Books) Update(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	dbCon, err := b.db.Copy()
	if err != nil {
		return errors.Wrap(router.ErrDBNotConfigured, "")
	}

	defer dbCon.Close()

	var newBook book.Book
	if err := router.Unmarshall(r.Body, &newBook); err != nil {
		return errors.Wrap(err, "")
	}

	err = book.UpdateBookbyID(ctx, dbCon, params["id"], &newBook)
	if err = checkError(err); err != nil {
		return errors.Wrapf(err, "Id: %s Book:%+v", params["id"], &newBook)
	}

	router.Response(ctx, w, nil, http.StatusNoContent)
	return nil
}
