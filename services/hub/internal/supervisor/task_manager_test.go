package supervisor

import (
	"math/rand"
	"testing"
	"time"

	"github.com/bfoody/Walmart-Scraper/domain"
	"github.com/bfoody/Walmart-Scraper/utils/uuid"
)

var fakeTasks []domain.ScrapeTask = []domain.ScrapeTask{
	domain.ScrapeTask{
		ID:                "1",
		Completed:         false,
		CreatedAt:         time.Now(),
		ScheduledFor:      time.Now().Add(1 * time.Second),
		ProductLocationID: "2",
		Repeat:            false,
		Interval:          1 * time.Second,
	},
	domain.ScrapeTask{
		ID:                "12",
		Completed:         false,
		CreatedAt:         time.Now(),
		ScheduledFor:      time.Now().Add(12 * time.Second),
		ProductLocationID: "2",
		Repeat:            false,
		Interval:          1 * time.Second,
	},
	domain.ScrapeTask{
		ID:                "3",
		Completed:         false,
		CreatedAt:         time.Now(),
		ScheduledFor:      time.Now().Add(3 * time.Second),
		ProductLocationID: "2",
		Repeat:            false,
		Interval:          1 * time.Second,
	},
	domain.ScrapeTask{
		ID:                "4",
		Completed:         false,
		CreatedAt:         time.Now(),
		ScheduledFor:      time.Now().Add(4 * time.Second),
		ProductLocationID: "2",
		Repeat:            false,
		Interval:          1 * time.Second,
	},
	domain.ScrapeTask{
		ID:                "8",
		Completed:         false,
		CreatedAt:         time.Now(),
		ScheduledFor:      time.Now().Add(8 * time.Second),
		ProductLocationID: "2",
		Repeat:            false,
		Interval:          1 * time.Second,
	},
	domain.ScrapeTask{
		ID:                "7",
		Completed:         false,
		CreatedAt:         time.Now(),
		ScheduledFor:      time.Now().Add(7 * time.Second),
		ProductLocationID: "2",
		Repeat:            false,
		Interval:          1 * time.Second,
	},
}

func fakeTaskGenerator(num int) []domain.ScrapeTask {
	tasks := []domain.ScrapeTask{}

	for i := 0; i < num; i++ {
		tasks = append(tasks, domain.ScrapeTask{
			ID:                uuid.Generate(),
			Completed:         false,
			CreatedAt:         time.Now(),
			ScheduledFor:      time.Now().Add(time.Duration(rand.Uint32()) * time.Second),
			ProductLocationID: "1",
			Repeat:            false,
			Interval:          1 * time.Second,
		})
	}

	return tasks
}

func TestQueueOrder(t *testing.T) {
	tm := NewTaskManager(nil)
	//	tasks := fakeTaskGenerator(5000)
	tasks := fakeTasks

	for i, task := range tasks {
		// Log every 500 indexes
		if i%500 == 0 {
			t.Logf("currently at index %d", i)
		}

		tm.pushTaskToQueue(task)
		tm.sortedEnqueue(task.ID)
	}

	tm.printListDebug()
	last := ""
	for e := tm.queue.Front(); e != nil; e = e.Next() {
		id := e.Value.(string)

		if last == "" {
			last = id
			continue
		}

		if !tm.tasks[id].ScheduledFor.After(tm.tasks[last].ScheduledFor) {
			t.Fatalf("task %s should be before last task %s", id, last)
		}

		last = id
	}

}
