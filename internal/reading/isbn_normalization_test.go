package reading

import "testing"

// TestISBNNormalizationEquivalence verifies that equivalent ISBN spellings —
// differing only in hyphens, surrounding whitespace, or a combination of both —
// all collapse to a single shelf record, while unrelated behavior stays intact.
func TestISBNNormalizationEquivalence(t *testing.T) {
	cases := []struct {
		name string
		first string
		equiv string
	}{
		{"hyphens removed", "978-1-4919-5059-2", "9781491950592"},
		{"leading and trailing spaces", "9781491950592", "  9781491950592  "},
		{"hyphens and surrounding spaces combined", "978-1-4919-5059-2", "  978-1-4919-5059-2  "},
		{"spaces then hyphens removed", "  978-1-4919-5059-2  ", "9781491950592"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewService()
			first, err := service.AddBook(CreateBookInput{Title: "First", Author: "A", ISBN: tc.first, TotalPages: 100})
			if err != nil {
				t.Fatalf("add first: %v", err)
			}
			if first.ISBN != "9781491950592" {
				t.Fatalf("expected normalized isbn 9781491950592, got %q", first.ISBN)
			}
			if _, err := service.AddBook(CreateBookInput{Title: "Second", Author: "B", ISBN: tc.equiv, TotalPages: 120}); err != ErrDuplicateISBN {
				t.Fatalf("expected ErrDuplicateISBN for equivalent isbn, got %v", err)
			}
		})
	}
}

// TestISBNDistinctAndExactDuplicate ensures non-equivalent ISBNs still create
// separate records and that an exact (already-normalized) duplicate is rejected.
func TestISBNDistinctAndExactDuplicate(t *testing.T) {
	service := NewService()
	if _, err := service.AddBook(CreateBookInput{Title: "Go", Author: "A", ISBN: "9780134190440", TotalPages: 380}); err != nil {
		t.Fatal(err)
	}
	// A genuinely different ISBN must still be created.
	second, err := service.AddBook(CreateBookInput{Title: "Rust", Author: "B", ISBN: "9781593278281", TotalPages: 280})
	if err != nil {
		t.Fatalf("expected distinct isbn to be created, got %v", err)
	}
	// An exact duplicate of an existing (already normalized) ISBN must be rejected.
	if _, err := service.AddBook(CreateBookInput{Title: "Dup", Author: "C", ISBN: "9780134190440", TotalPages: 10}); err != ErrDuplicateISBN {
		t.Fatalf("expected ErrDuplicateISBN for exact duplicate, got %v", err)
	}
	// Query and modification behavior must be unaffected by the normalization fix.
	if got, err := service.GetBook(second.ID); err != nil || got.Title != "Rust" {
		t.Fatalf("unexpected get: %+v %v", got, err)
	}
	updated, err := service.UpdateBook(second.ID, UpdateBookInput{Title: strPtr("Rust in Action")})
	if err != nil || updated.Title != "Rust in Action" {
		t.Fatalf("unexpected update: %+v %v", updated, err)
	}
	// Updating the ISBN to an equivalent form of the other book's ISBN must conflict.
	if _, err := service.UpdateBook(second.ID, UpdateBookInput{ISBN: strPtr(" 978-0-1341-9044-0 ")}); err != ErrDuplicateISBN {
		t.Fatalf("expected ErrDuplicateISBN on equivalent isbn update, got %v", err)
	}
}

func strPtr(s string) *string { return &s }
