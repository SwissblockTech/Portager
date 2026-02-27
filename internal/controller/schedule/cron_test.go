package schedule

import (
	"testing"
	"time"
)

func TestNextSyncTime_StandardCron(t *testing.T) {
	s := NewScheduler()

	// Reference time: 2026-02-26 10:30:00 UTC
	ref := time.Date(2026, 2, 26, 10, 30, 0, 0, time.UTC)

	// "0 */6 * * *" = at minute 0, every 6 hours (00:00, 06:00, 12:00, 18:00).
	// After 10:30, the next hit is 12:00 the same day.
	next, err := s.NextSyncTime("0 */6 * * *", ref)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := time.Date(2026, 2, 26, 12, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestNextSyncTime_EveryInterval(t *testing.T) {
	s := NewScheduler()

	ref := time.Date(2026, 2, 26, 10, 30, 0, 0, time.UTC)

	// "@every 1h" returns exactly ref + 1 hour.
	next, err := s.NextSyncTime("@every 1h", ref)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := ref.Add(1 * time.Hour)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestNextSyncTime_Descriptors(t *testing.T) {
	s := NewScheduler()
	ref := time.Date(2026, 2, 26, 10, 30, 0, 0, time.UTC)

	// These should parse without error and return a future time.
	for _, expr := range []string{"@daily", "@hourly", "@weekly", "@monthly"} {
		next, err := s.NextSyncTime(expr, ref)
		if err != nil {
			t.Errorf("descriptor %q: unexpected error: %v", expr, err)
			continue
		}
		if !next.After(ref) {
			t.Errorf("descriptor %q: expected time after %v, got %v", expr, ref, next)
		}
	}
}

func TestNextSyncTime_InvalidExpression(t *testing.T) {
	s := NewScheduler()
	ref := time.Date(2026, 2, 26, 10, 30, 0, 0, time.UTC)

	_, err := s.NextSyncTime("not-a-cron", ref)
	if err == nil {
		t.Fatal("expected error for invalid expression, got nil")
	}
}

func TestValidate_Valid(t *testing.T) {
	s := NewScheduler()

	for _, expr := range []string{"0 */6 * * *", "@every 1h", "@daily", "*/5 * * * *"} {
		if err := s.Validate(expr); err != nil {
			t.Errorf("expression %q should be valid, got error: %v", expr, err)
		}
	}
}

func TestValidate_Invalid(t *testing.T) {
	s := NewScheduler()

	for _, expr := range []string{"not-valid", "* * *", "60 * * * *"} {
		if err := s.Validate(expr); err == nil {
			t.Errorf("expression %q should be invalid, got nil error", expr)
		}
	}
}
