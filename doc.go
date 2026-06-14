// Package goqueue provides a Redis-backed task queue SDK for Go services.
//
// The root package is the public SDK facade. Task model implementation lives
// in the task subpackage and is re-exported here for the convenient goqueue.X
// API. Redis producer, worker, and scheduler APIs are added in later phases on
// top of this package boundary.
package goqueue
