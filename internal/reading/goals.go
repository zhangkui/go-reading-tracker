package reading

import (
	"sort"
	"time"
)

func goalKey(month string, goalType GoalType) string { return month + ":" + string(goalType) }
func parseMonth(month string) (time.Time, error)     { return time.Parse("2006-01", month) }

func (s *Service) SetGoal(goal Goal) (Goal, error) {
	if _, err := parseMonth(goal.Month); err != nil || goal.Target <= 0 || (goal.Type != GoalBooks && goal.Type != GoalPages) {
		return Goal{}, ErrInvalidInput
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.goals[goalKey(goal.Month, goal.Type)] = goal
	return goal, nil
}

func (s *Service) GetGoalProgress(month string, goalType GoalType) (GoalProgress, error) {
	requestedMonth, err := parseMonth(month)
	if err != nil {
		return GoalProgress{}, ErrInvalidInput
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	goal, exists := s.goals[goalKey(month, goalType)]
	if !exists {
		return GoalProgress{}, ErrGoalNotFound
	}
	current := 0
	for _, book := range s.books {
		switch goalType {
		case GoalBooks:
			if book.CompletedAt != nil && book.CompletedAt.Year() == requestedMonth.Year() && book.CompletedAt.Month() == requestedMonth.Month() {
				current++
			}
		case GoalPages:
			if book.UpdatedAt.Year() == requestedMonth.Year() && book.UpdatedAt.Month() == requestedMonth.Month() {
				current += book.CurrentPage
			}
		}
	}
	remaining := goal.Target - current
	if remaining < 0 {
		remaining = 0
	}
	percentage := current * 100 / goal.Target
	if percentage > 100 {
		percentage = 100
	}
	return GoalProgress{Goal: goal, Current: current, Remaining: remaining, Percentage: percentage}, nil
}

func (s *Service) Statistics(month string) (Statistics, error) {
	var requestedMonth time.Time
	var err error
	if month != "" {
		requestedMonth, err = parseMonth(month)
		if err != nil {
			return Statistics{}, ErrInvalidInput
		}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	authors := make(map[string]int)
	result := Statistics{Month: month}
	for _, book := range s.books {
		if month != "" && (book.UpdatedAt.Year() != requestedMonth.Year() || book.UpdatedAt.Month() != requestedMonth.Month()) {
			continue
		}
		result.PagesRead += book.CurrentPage
		if book.Status == StatusCompleted {
			result.CompletedBooks++
		}
		authors[book.Author]++
	}
	for author, count := range authors {
		result.Authors = append(result.Authors, AuthorCount{Author: author, Count: count})
	}
	sort.Slice(result.Authors, func(i, j int) bool { return result.Authors[i].Author < result.Authors[j].Author })
	return result, nil
}
