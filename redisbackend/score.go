package redisbackend

import "time"

func unixMillis(timestamp time.Time) int64 {
	return timestamp.UnixNano() / int64(time.Millisecond)
}
