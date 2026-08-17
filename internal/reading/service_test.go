package reading

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestBookLifecycle(t *testing.T) {
	service := NewService()
	book, err := service.AddBook(CreateBookInput{Title: "The Go Programming Language", Author: "Alan Donovan", ISBN: "9780134190440", TotalPages: 380, Tags: []string{"Go", "Programming"}})
	if err != nil {
		t.Fatal(err)
	}
	if book.Status != StatusPlanned || len(book.Tags) != 2 {
		t.Fatalf("unexpected book: %+v", book)
	}
	if _, err := service.AddBook(CreateBookInput{Title: "Duplicate", Author: "Someone", ISBN: "9780134190440", TotalPages: 10}); err != ErrDuplicateISBN {
		t.Fatalf("expected duplicate isbn, got %v", err)
	}
	updated, err := service.UpdateProgress(book.ID, 25)
	if err != nil || updated.Status != StatusReading {
		t.Fatalf("unexpected progress update: %+v %v", updated, err)
	}
	if err := service.DeleteBook(book.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := service.GetBook(book.ID); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestSearchAndPagination(t *testing.T) {
	service := NewService()
	inputs := []CreateBookInput{
		{Title: "Go in Action", Author: "William Kennedy", ISBN: "1", TotalPages: 250, Tags: []string{"go"}},
		{Title: "Distributed Services", Author: "Travis Jeffery", ISBN: "2", TotalPages: 300, Tags: []string{"go", "systems"}},
		{Title: "Clean Architecture", Author: "Robert Martin", ISBN: "3", TotalPages: 350, Tags: []string{"architecture"}},
	}
	for _, input := range inputs {
		if _, err := service.AddBook(input); err != nil {
			t.Fatal(err)
		}
	}
	result := service.SearchBooks(SearchOptions{Title: "go", Tags: []string{"go"}, Page: 1, PageSize: 1})
	if result.Total != 1 || len(result.Items) != 1 || result.Items[0].Title != "Go in Action" {
		t.Fatalf("unexpected search result: %+v", result)
	}
}

func TestGoalProgressAndStatistics(t *testing.T) {
	service := NewService()
	now := time.Date(2026, time.August, 12, 10, 0, 0, 0, time.UTC)
	service.setNow(func() time.Time { return now })
	if _, err := service.SetGoal(Goal{Month: "2026-08", Type: GoalPages, Target: 100}); err != nil {
		t.Fatal(err)
	}
	book, err := service.AddBook(CreateBookInput{Title: "Go", Author: "A", ISBN: "4", TotalPages: 120})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.UpdateProgress(book.ID, 40); err != nil {
		t.Fatal(err)
	}
	progress, err := service.GetGoalProgress("2026-08", GoalPages)
	if err != nil || progress.Current != 40 || progress.Percentage != 40 {
		t.Fatalf("unexpected goal progress: %+v %v", progress, err)
	}
	stats, err := service.Statistics("2026-08")
	if err != nil || stats.PagesRead != 40 || len(stats.Authors) != 1 {
		t.Fatalf("unexpected statistics: %+v %v", stats, err)
	}
}

func TestImportRecords(t *testing.T) {
	service := NewService()
	payload := `[{"title":"Imported","author":"Reader","isbn":"5","total_pages":90,"current_page":10,"tags":["memoir"]}]`
	result, err := service.Import(context.Background(), strings.NewReader(payload))
	if err != nil || result.Imported != 1 || len(result.Errors) != 0 {
		t.Fatalf("unexpected import result: %+v %v", result, err)
	}
}
