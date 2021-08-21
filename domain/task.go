package domain

import "time"

// A ScrapeTask represents a scraping job scheduled for some time in the future.
type ScrapeTask struct {
	ID                string        `db:"id"`                  // the entity's unique ID
	Completed         bool          `db:"completed"`           // whether or not the task was completed
	CreatedAt         time.Time     `db:"created_at"`          // when the task was created
	ScheduledFor      time.Time     `db:"scheduled_for"`       // when the task is to be completed
	ProductLocationID string        `db:"product_location_id"` // the ID of the product-location pair to be scraped
	Repeat            bool          `db:"repeat"`              // whether or not to schedule another task after completion
	Interval          time.Duration `db:"interval"`            // the duration between repetitions of the task, will be added to the current time when the schedule repeats
}
