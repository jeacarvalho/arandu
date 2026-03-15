package session

import (
	"errors"
	"time"
)

var (
	ErrInvalidDate    = errors.New("invalid session date")
	ErrSummaryTooLong = errors.New("summary is too long")
)

func (s *Session) Update(date time.Time, summary string) error {
	if date.IsZero() || date.After(time.Now()) {
		return ErrInvalidDate
	}

	if len(summary) > 10000 { // Example length limit
		return ErrSummaryTooLong
	}

	s.Date = date
	s.Summary = summary
	s.UpdatedAt = time.Now()

	return nil
}
