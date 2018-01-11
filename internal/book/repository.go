package book

import (
	"context"

	"github.com/huduma/internal/mongo"
)

//Repository represents the mongo repository
type Repository interface {
	ListAllBooks(ctx context.Context, dbCon *mongo.BooksDB) ([]Book, error)
	FindBookByID(ctx context.Context, dbCon *mongo.BooksDB, bookID string) (*Book, error)
	DeleteBookByID(ctx context.Context, dbCon *mongo.BooksDB, bookID string) (*Book, error)
	AddNewBook(ctx context.Context, dbCon *mongo.BooksDB, bk *Book) (*Book, error)
	UpdateBookbyID(ctx context.Context, dbCon *mongo.BooksDB, bkID string, cb *Book) error
}
