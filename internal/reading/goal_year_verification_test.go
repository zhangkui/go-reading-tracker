package reading

import (
	"testing"
	"time"
)

func TestMonthlyBookGoalDoesNotCountSameMonthFromOtherYears(t *testing.T) {
	service := NewService()
	now := time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC)
	service.setNow(func() time.Time { return now })
	if _, err := service.AddBook(CreateBookInput{Title: "Old", Author: "A", ISBN: "goal-1", TotalPages: 10, Status: StatusCompleted}); err != nil {
		t.Fatal(err)
	}
	now = time.Date(2026, time.January, 20, 0, 0, 0, 0, time.UTC)
	if _, err := service.AddBook(CreateBookInput{Title: "Current", Author: "B", ISBN: "goal-2", TotalPages: 10, Status: StatusCompleted}); err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetGoal(Goal{Month: "2026-01", Type: GoalBooks, Target: 3}); err != nil {
		t.Fatal(err)
	}
	progress, err := service.GetGoalProgress("2026-01", GoalBooks)
	if err != nil {
		t.Fatal(err)
	}
	if progress.Current != 1 {
		t.Fatalf("expected one completion in 2026-01, got %+v", progress)
	}
}
