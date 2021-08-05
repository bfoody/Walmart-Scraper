package supervisor

import (
	"container/list"
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/bfoody/Walmart-Scraper/domain"
	"github.com/bfoody/Walmart-Scraper/services/hub"
)

const (
	defaultLimit = 512
)

// A TaskManager helps to manage an internal queue of tasks.
type TaskManager struct {
	service       hub.Service
	queueMutex    *sync.RWMutex
	tasks         map[string]domain.ScrapeTask
	_queue        []string
	queue         *list.List
	resolvedTasks map[string]bool
}

// NewTaskManager creates and returns a new TaskManager.
func NewTaskManager(service hub.Service) *TaskManager {
	return &TaskManager{
		service:       service,
		queueMutex:    &sync.RWMutex{},
		tasks:         map[string]domain.ScrapeTask{},
		_queue:        []string{},
		queue:         list.New(),
		resolvedTasks: map[string]bool{},
	}
}

// Initialize starts the task manager and fetches new tasks.
func (t *TaskManager) Initialize() error {
	return t.fetchTaskList()
}

// fetchTaskList pulls new tasks into the TaskManager's queue.
func (t *TaskManager) fetchTaskList() error {
	tasks, err := t.service.FetchUpcomingTasks(defaultLimit)
	if err != nil {
		return err
	}

	t.queueMutex.Lock()
	defer t.queueMutex.Unlock()

	for _, task := range tasks {
		t.pushTaskToQueue(task)
	}

	return nil
}

// pushTaskToQueue pushes a task into the internal task queue.
func (t *TaskManager) pushTaskToQueue(task domain.ScrapeTask) {
	t.tasks[task.ID] = task

	for _, id := range t._queue {
		if id == task.ID {
			return
		}
	}

	t._queue = append(t._queue, task.ID)
}

func (t *TaskManager) sortedEnqueue(id string) {
	// If the queue is empty, push immediately.
	if t.queue.Len() < 1 {
		t.queue.PushBack(id)
		return
	}

	// Try to push in before the first element scheduled after this one.
	for e := t.queue.Back(); e != nil; e = e.Prev() {
		currID := e.Value.(string)
		if t.tasks[id].ScheduledFor.After(t.tasks[currID].ScheduledFor) {
			t.queue.InsertAfter(id, e)
			return
		}
	}

	// If nothing found, push in the front.
	t.queue.PushFront(id)
}

func (t *TaskManager) printListDebug() {
	str := ""

	for e := t.queue.Front(); e != nil; e = e.Next() {
		str += e.Value.(string) + " "
	}

	fmt.Println(fmt.Sprintf("[ %s ]", str))
}

// sortQueue sorts the queue by upcoming tasks.
//
// Tasks will be sorted such that newer tasks will be at the end of the queue.
func (t *TaskManager) sortQueue() {
	sort.SliceStable(t._queue, func(i, j int) bool {
		return t.tasks[t._queue[i]].ScheduledFor.After(t.tasks[t._queue[j]].ScheduledFor)
	})
}

// resolveTask removes a task from the queue.
func (t *TaskManager) resolveTask(id string) {
	t.resolvedTasks[id] = true
}

// pushTask adds a task to the front of the internal queue.
func (t *TaskManager) pushTask(task domain.ScrapeTask) error {
	t.queueMutex.Lock()
	defer t.queueMutex.Unlock()

	t.tasks[task.ID] = task
	t._queue = append([]string{task.ID}, t._queue...)

	return nil
}

// popTask pops a task off the queue and returns it.
func (t *TaskManager) popTask() (*domain.ScrapeTask, error) {
	if len(t._queue) < 1 {
		return nil, errors.New("queue is empty")
	}

	t.queueMutex.Lock()
	defer t.queueMutex.Unlock()

	i := len(t._queue) - 1
	id := t._queue[i]

	// Take one element off the end of the array.
	t._queue = t._queue[:i]

	task := t.tasks[id]

	return &task, nil
}
