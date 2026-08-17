package reading

import "testing"

func TestSearchRequiresEveryRequestedTag(t *testing.T) {
	service := NewService()
	inputs := []CreateBookInput{
		{Title: "Both", Author: "A", ISBN: "tags-1", TotalPages: 10, Tags: []string{"go", "api"}},
		{Title: "Only Go", Author: "B", ISBN: "tags-2", TotalPages: 10, Tags: []string{"go"}},
		{Title: "Only API", Author: "C", ISBN: "tags-3", TotalPages: 10, Tags: []string{"api"}},
	}
	for _, input := range inputs {
		if _, err := service.AddBook(input); err != nil {
			t.Fatal(err)
		}
	}
	result := service.SearchBooks(SearchOptions{Tags: []string{"go", "api"}, Page: 1, PageSize: 20})
	if result.Total != 1 || len(result.Items) != 1 || result.Items[0].Title != "Both" {
		t.Fatalf("unexpected multi-tag result: %+v", result)
	}
}
