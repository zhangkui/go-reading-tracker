package reading

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

func (s *Service) Import(ctx context.Context, reader io.Reader) (ImportResult, error) {
	if err := ctx.Err(); err != nil {
		return ImportResult{}, err
	}
	var records []ImportRecord
	if err := json.NewDecoder(reader).Decode(&records); err != nil {
		return ImportResult{}, fmt.Errorf("decode import records: %w", err)
	}
	result := ImportResult{}
	for index, record := range records {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		status := StatusPlanned
		if record.CurrentPage > 0 {
			status = StatusReading
		}
		if record.TotalPages > 0 && record.CurrentPage == record.TotalPages {
			status = StatusCompleted
		}
		_, err := s.AddBook(CreateBookInput{Title: record.Title, Author: record.Author, ISBN: record.ISBN, TotalPages: record.TotalPages, CurrentPage: record.CurrentPage, Status: status, Tags: record.Tags})
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("record %d: %v", index, err))
			continue
		}
		result.Imported++
	}
	return result, nil
}
