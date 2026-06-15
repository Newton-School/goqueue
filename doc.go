// Package goqueue provides a Redis-backed task queue SDK for Go services.
//
// The root package is the public SDK facade. Task model implementation lives
// in the task subpackage and is re-exported here for the convenient goqueue.X
// API. Redis producer, worker, and scheduler APIs are available on top of this
// package boundary, while redisbackend powers durable queue storage, retries,
// dead-letter queues, pending recovery, and periodic dispatch coordination.
package goqueue
