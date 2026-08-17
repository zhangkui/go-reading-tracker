package reading

import "time"

type Status string

const (
	StatusPlanned   Status = "planned"
	StatusReading   Status = "reading"
	StatusCompleted Status = "completed"
)

type Book struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Author      string     `json:"author"`
	ISBN        string     `json:"isbn"`
	TotalPages  int        `json:"total_pages"`
	CurrentPage int        `json:"current_page"`
	Status      Status     `json:"status"`
	Tags        []string   `json:"tags"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type CreateBookInput struct {
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	ISBN        string   `json:"isbn"`
	TotalPages  int      `json:"total_pages"`
	CurrentPage int      `json:"current_page"`
	Status      Status   `json:"status"`
	Tags        []string `json:"tags"`
}

type UpdateBookInput struct {
	Title      *string   `json:"title,omitempty"`
	Author     *string   `json:"author,omitempty"`
	ISBN       *string   `json:"isbn,omitempty"`
	TotalPages *int      `json:"total_pages,omitempty"`
	Status     *Status   `json:"status,omitempty"`
	Tags       *[]string `json:"tags,omitempty"`
}

type SearchOptions struct {
	Title    string
	Author   string
	Status   Status
	Tags     []string
	Page     int
	PageSize int
}
type SearchResult struct {
	Items    []Book `json:"items"`
	Total    int    `json:"total"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}
type GoalType string

const (
	GoalBooks GoalType = "books"
	GoalPages GoalType = "pages"
)

type Goal struct {
	Month  string   `json:"month"`
	Type   GoalType `json:"type"`
	Target int      `json:"target"`
}
type GoalProgress struct {
	Goal       Goal `json:"goal"`
	Current    int  `json:"current"`
	Remaining  int  `json:"remaining"`
	Percentage int  `json:"percentage"`
}
type ImportRecord struct {
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	ISBN        string   `json:"isbn"`
	TotalPages  int      `json:"total_pages"`
	CurrentPage int      `json:"current_page"`
	Tags        []string `json:"tags"`
}
type ImportResult struct {
	Imported int      `json:"imported"`
	Errors   []string `json:"errors,omitempty"`
}
type AuthorCount struct {
	Author string `json:"author"`
	Count  int    `json:"count"`
}
type Statistics struct {
	Month          string        `json:"month,omitempty"`
	CompletedBooks int           `json:"completed_books"`
	PagesRead      int           `json:"pages_read"`
	Authors        []AuthorCount `json:"authors"`
}
