package supervisor

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
	"time"

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
	t.queueMutex.Lock()
	defer t.queueMutex.Unlock()

	t.tasks[task.ID] = task

	t.sortedEnqueue(task.ID)
}

// sortedEnqueue adds a task ID to the task queue after it has been added to the map.
func (t *TaskManager) sortedEnqueue(id string) {
	// If the queue is empty, push immediately.
	if t.queue.Len() < 1 {
		t.queue.PushBack(id)
		return
	}

	// Try to push in before the first element scheduled after this one.
	for e := t.queue.Back(); e != nil; e = e.Prev() {
		currID := e.Value.(string)

		if currID == id {
			// Don't add duplicate values.
			return
		}

		if t.tasks[id].ScheduledFor.After(t.tasks[currID].ScheduledFor) {
			t.queue.InsertAfter(id, e)
			return
		}
	}

	// If nothing found, push in the front.
	t.queue.PushFront(id)
}

// printListDebug pretty-prints the queue to stdout.
func (t *TaskManager) printListDebug() {
	str := ""

	for e := t.queue.Front(); e != nil; e = e.Next() {
		str += e.Value.(string) + " "
	}

	fmt.Println(fmt.Sprintf("[ %s ]", str))
}

// resolveTask removes a task from the queue.
func (t *TaskManager) resolveTask(id string) {
	t.resolvedTasks[id] = true
}

// popTask locks the mutex, pops a task off the queue, and returns it.
func (t *TaskManager) popTask() (*domain.ScrapeTask, error) {
	if t.queue.Len() < 1 {
		return nil, errors.New("queue is empty")
	}

	t.queueMutex.Lock()
	defer t.queueMutex.Unlock()

	return t._popTask()
}

// _popTask pops a task and will recurse if the obtained task is already resolved.
func (t *TaskManager) _popTask() (*domain.ScrapeTask, error) {
	el := t.queue.Front()
	id := el.Value.(string)

	t.queue.Remove(el)

	// If the task is already resolved, recurse and pop another task.
	if _, ok := t.resolvedTasks[id]; ok {
		return t._popTask()
	}

	task := t.tasks[id]

	return &task, nil
}

// frontQueueID returns the ID of the frontmost (earliest timestamp) item on the queue.
func (t *TaskManager) frontQueueID() (string, error) {
	if t.queue.Len() < 1 {
		return "", errors.New("queue is empty")
	}

	t.queueMutex.RLock()
	defer t.queueMutex.RUnlock()

	return t._frontQueueID()
}

func (t *TaskManager) _frontQueueID() (string, error) {
	if t.queue.Len() < 1 {
		return "", errors.New("queue is empty")
	}

	el := t.queue.Front()
	id := el.Value.(string)

	// If the task is already resolved, recurse and pop another task.
	if _, ok := t.resolvedTasks[id]; ok {
		t.queue.Remove(el)

		return t._frontQueueID()
	}

	return id, nil
}

// timeUntilNextDueTask returns the amount of nanoseconds until the next task is due.
func (t *TaskManager) timeUntilNextDueTask() (time.Duration, error) {
	id, err := t.frontQueueID()
	if err != nil {
		return time.Duration(0), err
	}

	return time.Until(t.tasks[id].ScheduledFor), nil
}
