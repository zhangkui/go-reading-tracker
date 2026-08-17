package reading

import (
	"sort"
	"strings"
)

func normalizeISBN(value string) string {
	value = strings.TrimSpace(value)
	return strings.ReplaceAll(value, "-", "")
}

func normalizeTags(tags []string) []string {
	seen := make(map[string]struct{}, len(tags))
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if tag == "" {
			continue
		}
		if _, exists := seen[tag]; exists {
			continue
		}
		seen[tag] = struct{}{}
		result = append(result, tag)
	}
	sort.Strings(result)
	return result
}

func validateBookInput(title, author, isbn string, totalPages, currentPage int) error {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(author) == "" || normalizeISBN(isbn) == "" {
		return ErrInvalidInput
	}
	if totalPages <= 0 || currentPage < 0 || currentPage > totalPages {
		return ErrInvalidInput
	}
	return nil
}

func (s *Service) AddBook(input CreateBookInput) (Book, error) {
	if err := validateBookInput(input.Title, input.Author, input.ISBN, input.TotalPages, input.CurrentPage); err != nil {
		return Book{}, err
	}
	isbn := normalizeISBN(input.ISBN)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.byISBN[isbn]; exists {
		return Book{}, ErrDuplicateISBN
	}
	now := s.now().UTC()
	status := input.Status
	if status == "" {
		status = StatusPlanned
	}
	if input.CurrentPage > 0 && status == StatusPlanned {
		status = StatusReading
	}
	book := Book{ID: s.nextID, Title: strings.TrimSpace(input.Title), Author: strings.TrimSpace(input.Author), ISBN: isbn, TotalPages: input.TotalPages, CurrentPage: input.CurrentPage, Status: status, Tags: normalizeTags(input.Tags), CreatedAt: now, UpdatedAt: now}
	if status == StatusCompleted {
		completedAt := now
		book.CompletedAt = &completedAt
		book.CurrentPage = book.TotalPages
	}
	s.books[book.ID] = book
	s.byISBN[book.ISBN] = book.ID
	s.nextID++
	return cloneBook(book), nil
}

func (s *Service) GetBook(id int64) (Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	book, exists := s.books[id]
	if !exists {
		return Book{}, ErrNotFound
	}
	return cloneBook(book), nil
}

func (s *Service) UpdateBook(id int64, input UpdateBookInput) (Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	book, exists := s.books[id]
	if !exists {
		return Book{}, ErrNotFound
	}
	if input.Title != nil {
		book.Title = strings.TrimSpace(*input.Title)
	}
	if input.Author != nil {
		book.Author = strings.TrimSpace(*input.Author)
	}
	if input.ISBN != nil {
		isbn := normalizeISBN(*input.ISBN)
		if owner, exists := s.byISBN[isbn]; exists && owner != id {
			return Book{}, ErrDuplicateISBN
		}
		delete(s.byISBN, book.ISBN)
		book.ISBN = isbn
		s.byISBN[isbn] = id
	}
	if input.TotalPages != nil {
		book.TotalPages = *input.TotalPages
	}
	if input.Status != nil {
		book.Status = *input.Status
		if book.Status == StatusCompleted {
			completedAt := s.now().UTC()
			book.CompletedAt = &completedAt
			book.CurrentPage = book.TotalPages
		}
	}
	if input.Tags != nil {
		book.Tags = normalizeTags(*input.Tags)
	}
	if err := validateBookInput(book.Title, book.Author, book.ISBN, book.TotalPages, book.CurrentPage); err != nil {
		return Book{}, err
	}
	book.UpdatedAt = s.now().UTC()
	s.books[id] = book
	return cloneBook(book), nil
}

func (s *Service) DeleteBook(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	book, exists := s.books[id]
	if !exists {
		return ErrNotFound
	}
	delete(s.byISBN, book.ISBN)
	delete(s.books, id)
	return nil
}

func (s *Service) allBooks() []Book {
	books := make([]Book, 0, len(s.books))
	for _, book := range s.books {
		books = append(books, cloneBook(book))
	}
	sort.Slice(books, func(i, j int) bool { return books[i].ID < books[j].ID })
	return books
}
