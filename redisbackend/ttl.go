package redisbackend

import "time"

func ttlSeconds(ttl time.Duration) int64 {
	if ttl <= 0 {
		return 1
	}

	seconds := int64(ttl / time.Second)
	if ttl%time.Second != 0 {
		seconds++
	}
	if seconds < 1 {
		return 1
	}

	return seconds
}
