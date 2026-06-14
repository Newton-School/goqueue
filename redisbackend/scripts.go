package redisbackend

func readyEnqueueScript() string {
	return `
redis.call('SET', KEYS[1], ARGV[1], 'EX', ARGV[2])
return redis.call('XADD', KEYS[2], '*', 'id', ARGV[3], 'message', ARGV[1])
`
}

func scheduledEnqueueScript() string {
	return `
redis.call('SET', KEYS[1], ARGV[1], 'EX', ARGV[2])
redis.call('ZADD', KEYS[2], ARGV[3], ARGV[4])
return ARGV[4]
`
}
