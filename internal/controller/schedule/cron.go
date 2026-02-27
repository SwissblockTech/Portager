package schedule

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// Scheduler parses cron expressions and computes the next sync time.
//
// It uses robfig/cron's Parser configured for standard 5-field cron expressions
// (minute, hour, day-of-month, month, day-of-week) plus descriptors like
// @every, @daily, and @hourly. The seconds field is NOT supported — this matches
// what users expect from crontab(5) syntax.
type Scheduler struct {
	parser cron.Parser
}

// NewScheduler creates a Scheduler with a pre-configured parser.
// The parser is created once and reused for all schedule evaluations.
func NewScheduler() *Scheduler {
	return &Scheduler{
		parser: cron.NewParser(
			cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
		),
	}
}

// NextSyncTime parses a cron expression and returns the next occurrence after
// the given time. The `after` parameter makes this function deterministic and
// testable — callers pass time.Now() in production or a fixed time in tests.
func (s *Scheduler) NextSyncTime(cronExpr string, after time.Time) (time.Time, error) {
	sched, err := s.parser.Parse(cronExpr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cron expression %q: %w", cronExpr, err)
	}
	return sched.Next(after), nil
}

// Validate checks whether a cron expression can be parsed without computing
// a time. Use this for fast validation in the reconcile loop before proceeding
// with the sync.
func (s *Scheduler) Validate(cronExpr string) error {
	_, err := s.parser.Parse(cronExpr)
	if err != nil {
		return fmt.Errorf("invalid cron expression %q: %w", cronExpr, err)
	}
	return nil
}
