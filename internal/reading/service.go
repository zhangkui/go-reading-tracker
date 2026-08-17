package reading

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrNotFound      = errors.New("book not found")
	ErrDuplicateISBN = errors.New("isbn already exists")
	ErrInvalidInput  = errors.New("invalid input")
	ErrGoalNotFound  = errors.New("goal not found")
)

type Service struct {
	mu     sync.RWMutex
	books  map[int64]Book
	byISBN map[string]int64
	goals  map[string]Goal
	nextID int64
	now    func() time.Time
}

func NewService() *Service {
	return &Service{books: make(map[int64]Book), byISBN: make(map[string]int64), goals: make(map[string]Goal), nextID: 1, now: time.Now}
}

func (s *Service) setNow(now func() time.Time) { s.now = now }
func cloneBook(book Book) Book                 { book.Tags = append([]string(nil), book.Tags...); return book }
