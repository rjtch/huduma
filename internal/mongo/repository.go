package mongo

//Repository represents the mongo repository
type Repository interface {
	ListAllBooks([]*Book, error)
	FindBookByID(*Book, error)
	DeleteBookByID(*Book, error)
	AddNewBook(b *Book) error
	UpdateBookbyID(*Book, error)
}
