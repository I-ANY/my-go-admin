package queue

import (
	"github.com/RichardKnop/machinery/v2"
	"github.com/pkg/errors"
	"sync"
)

type QueueCollection struct {
	queues map[string]*machinery.Server
	lock   *sync.RWMutex
}

var (
	Queue = &QueueCollection{
		queues: make(map[string]*machinery.Server),
		lock:   new(sync.RWMutex),
	}
)

func (q *QueueCollection) AddQueue(queueName string, queue *machinery.Server) error {
	if queue == nil {
		return errors.New("server is nil")
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	if _, ok := q.queues[queueName]; ok {
		return errors.Errorf("queue %s already exists", queueName)
	}
	q.queues[queueName] = queue
	return nil
}
func (q *QueueCollection) GetQueue(queueName string) *machinery.Server {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return q.queues[queueName]
}
