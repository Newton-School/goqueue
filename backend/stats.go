package backend

import "github.com/Newton-School/goqueue/task"

// QueueStatsRequest asks a backend for queue storage counts.
type QueueStatsRequest struct {
	Queue task.QueueName
}

// QueueStats reports storage counts for a queue.
type QueueStats struct {
	Queue          task.QueueName
	ReadyCount     int64
	ScheduledCount int64
}

// Validate verifies the stats request.
func (r QueueStatsRequest) Validate() error {
	return task.ValidateQueueName(r.Queue.String())
}
