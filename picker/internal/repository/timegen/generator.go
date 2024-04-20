package timegen

import "time"

// TimeGenerator is a generator of time.
// It is used to generate time for the system in a consistent way.
type TimeGenerator struct{}

func NewTimeGenerator() *TimeGenerator {
	return &TimeGenerator{}
}

// Now returns the current time in UTC.
func (TimeGenerator) Now() time.Time {
	return time.Now().UTC()
}
