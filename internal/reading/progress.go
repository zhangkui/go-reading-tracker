package reading

func (s *Service) UpdateProgress(id int64, currentPage int) (Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	book, exists := s.books[id]
	if !exists {
		return Book{}, ErrNotFound
	}
	if currentPage < 0 || currentPage > book.TotalPages {
		return Book{}, ErrInvalidInput
	}
	book.CurrentPage = currentPage
	book.UpdatedAt = s.now().UTC()
	if currentPage >= book.TotalPages {
		book.Status = StatusCompleted
		completedAt := book.UpdatedAt
		book.CompletedAt = &completedAt
	} else if currentPage > 0 {
		book.Status = StatusReading
		book.CompletedAt = nil
	} else {
		book.Status = StatusPlanned
		book.CompletedAt = nil
	}
	s.books[id] = book
	return cloneBook(book), nil
}
