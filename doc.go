// Package goqueue provides a Redis-backed task queue SDK for Go services.
//
// The Phase 1 surface establishes application construction, configuration
// validation, task envelopes, payload codecs, handler contracts, and task
// registration. Redis producer, worker, and scheduler APIs are added in later
// phases on top of this package boundary.
package goqueue
