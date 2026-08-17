package reading

import "testing"

func TestNormalizedISBNIsUnique(t *testing.T) {
	service := NewService()
	if _, err := service.AddBook(CreateBookInput{Title: "First", Author: "A", ISBN: "978-1-4919-5059-2", TotalPages: 100}); err != nil {
		t.Fatal(err)
	}
	if _, err := service.AddBook(CreateBookInput{Title: "Second", Author: "B", ISBN: " 9781491950592 ", TotalPages: 120}); err != ErrDuplicateISBN {
		t.Fatalf("expected ErrDuplicateISBN, got %v", err)
	}
}
