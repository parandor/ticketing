package ticketing

import (
	"context"
	"errors"
//	"strconv"
	"sync"
	"net/http"

	connect "connectrpc.com/connect"
	v1 "github.com/parandor/ticketing/internal/gen/proto/library/v1"
	library_v1 "github.com/parandor/ticketing/internal/gen/proto/library/v1/libraryv1connect"
)

// BookRegistry represents the book registrar
type BookRegistry struct {
	mu     sync.RWMutex
	books  map[string]*v1.Book
	nextID int64 // Auto-incrementing ID
}

// NewBookRegistry creates a new instance of BookRegistry
func NewBookRegistry() *BookRegistry {
	return &BookRegistry{
		books:  make(map[string]*v1.Book),
		nextID: 1, // Initialize nextID to start from 1
	}
}

type LibraryService struct {
	bookRegistry *BookRegistry
}

func NewLibraryService() (string, http.Handler) {

	service := &LibraryService{
		bookRegistry: NewBookRegistry(),
	}

	service.bookRegistry.AddBook(&v1.Book{
		Id:     "123",
		Title:  "Sample Book",
		Author: "Sample Author",
	})

	path, handler := library_v1.NewLibraryServiceHandler(service)

	return path, handler
}

func (s *BookRegistry) GetAllBooks() ([]*v1.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	books := make([]*v1.Book, 0)
	for _, book := range s.books {
		books = append(books, book)
	}
	return books, nil
}

func (s *LibraryService) BorrowBook(ctx context.Context, req *connect.Request[v1.BorrowBookRequest]) (*connect.Response[v1.BorrowBookResponse], error) {
	bookID := req.Msg.BookId
	// Check if the book exists
	book, ok := s.bookRegistry.GetBook(bookID)
	if !ok {
		return nil, errors.New("book not found")
	}

	// Handle book borrowing logic here...
	// For example, mark the book as borrowed by the user identified by req.Msg.UserId
	// You may want to perform additional validation or business logic here

	// Construct the response
	response := &connect.Response[v1.BorrowBookResponse]{
		Msg: &v1.BorrowBookResponse{
			Success: true,
			Message: "Book checked out successfully",
			Book:    book,
		},
	}
	return response, nil
}

func (s *LibraryService) ReturnBook(ctx context.Context, req *connect.Request[v1.ReturnBookRequest]) (*connect.Response[v1.ReturnBookResponse], error) {
	bookID := req.Msg.BookId
	// Check if the book exists
	_, ok := s.bookRegistry.GetBook(bookID)
	if !ok {
		return nil, errors.New("book not found")
	}


	// Handle book returning logic here...
	// For example, mark the book as returned
	// You may want to perform additional validation or business logic here
	s.bookRegistry.RemoveBook(bookID)

	// Construct the response
	response := &connect.Response[v1.ReturnBookResponse]{
		Msg: &v1.ReturnBookResponse{
			Success: true,
			Message: "Book returned successfully",
		},
	}
	return response, nil
}

func (s *LibraryService) ListBooks(ctx context.Context, req *connect.Request[v1.ListBooksRequest]) (*connect.Response[v1.ListBooksResponse], error) {
	// Retrieve all books from the registry
	books, err := s.bookRegistry.GetAllBooks()
	if err != nil {
		return nil, errors.New("failed to get all books")
	}

	// Construct the response
	response := &connect.Response[v1.ListBooksResponse]{
		Msg: &v1.ListBooksResponse{
			Books: books,
		},
	}
	return response, nil
}

// AddBook adds a book to the registry
func (r *BookRegistry) AddBook(book *v1.Book) string {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.books[book.Id] = book
	return book.Id
}

// GetBook retrieves a book from the registry by its ID
func (r *BookRegistry) GetBook(bookID string) (*v1.Book, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	book, ok := r.books[bookID]
	return book, ok
}

// RemoveBook removes a book from the registry by its ID
func (r *BookRegistry) RemoveBook(bookID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.books, bookID)
}
