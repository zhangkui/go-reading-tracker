package reading

import "testing"

func TestProgressAtLastPageCompletesBook(t *testing.T) {
	service := NewService()
	book, err := service.AddBook(CreateBookInput{Title: "Finish", Author: "A", ISBN: "progress-1", TotalPages: 240})
	if err != nil {
		t.Fatal(err)
	}
	book, err = service.UpdateProgress(book.ID, 240)
	if err != nil {
		t.Fatal(err)
	}
	if book.Status != StatusCompleted || book.CompletedAt == nil {
		t.Fatalf("expected completed book, got %+v", book)
	}
}
