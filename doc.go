// Package goqueue provides a Redis-backed task queue SDK for Go services.
//
// The root package is the public SDK facade. Task model implementation lives
// in the task subpackage and is re-exported here for the convenient goqueue.X
// API. Redis producer APIs are available on top of this package boundary and a
// dedicated redisbackend package powers durable queue storage.
package goqueue
