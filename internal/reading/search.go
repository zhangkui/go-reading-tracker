package reading

import "strings"

func matchesTags(bookTags, requested []string) bool {
	if len(requested) == 0 {
		return true
	}
	available := make(map[string]struct{}, len(bookTags))
	for _, tag := range bookTags {
		available[strings.ToLower(strings.TrimSpace(tag))] = struct{}{}
	}
	for _, tag := range requested {
		if _, exists := available[strings.ToLower(strings.TrimSpace(tag))]; exists {
			return true
		}
	}
	return false
}

func (s *Service) SearchBooks(options SearchOptions) SearchResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	page := options.Page
	if page < 1 {
		page = 1
	}
	pageSize := options.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	title := strings.ToLower(strings.TrimSpace(options.Title))
	author := strings.ToLower(strings.TrimSpace(options.Author))
	requestedTags := normalizeTags(options.Tags)
	filtered := make([]Book, 0)
	for _, book := range s.allBooks() {
		if title != "" && !strings.Contains(strings.ToLower(book.Title), title) {
			continue
		}
		if author != "" && !strings.Contains(strings.ToLower(book.Author), author) {
			continue
		}
		if options.Status != "" && book.Status != options.Status {
			continue
		}
		if !matchesTags(book.Tags, requestedTags) {
			continue
		}
		filtered = append(filtered, book)
	}
	start := (page - 1) * pageSize
	if start > len(filtered) {
		start = len(filtered)
	}
	end := start + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}
	return SearchResult{Items: filtered[start:end], Total: len(filtered), Page: page, PageSize: pageSize}
}
