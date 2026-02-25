package tools

import "time"

type Clock struct {
	startAt time.Time
	current time.Time
	step    time.Duration
}

func (c *Clock) GetStartAt() time.Time {
	return c.startAt
}
func NewClock(start time.Time, step time.Duration) *Clock {
	return &Clock{
		startAt: start,
		step:    step,
	}
}

func (c *Clock) Next() time.Time {
	if c.current.IsZero() {
		c.current = c.startAt.Add(c.step)
	} else {
		c.current = c.current.Add(c.step)
	}
	return c.current
}
