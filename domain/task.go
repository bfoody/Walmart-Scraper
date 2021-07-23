package domain

import "time"

// A ScrapeTask represents a scraping job scheduled for some time in the future.
type ScrapeTask struct {
	ID                string        // the entity's unique ID
	Completed         bool          // whether or not the task was completed
	CreatedAt         time.Time     // when the task was created
	ScheduledFor      time.Time     // when the task is to be completed
	ProductLocationID string        // the ID of the product-location pair to be scraped
	Repeat            bool          // whether or not to schedule another task after completion
	Interval          time.Duration // the duration between repetitions of the task, will be added to the current time when the schedule repeats
}
