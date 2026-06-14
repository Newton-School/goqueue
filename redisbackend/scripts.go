package redisbackend

func readyEnqueueScript() string {
	return `
redis.call('SET', KEYS[1], ARGV[1], 'EX', ARGV[2])
return redis.call('XADD', KEYS[2], '*', 'id', ARGV[3], 'message', ARGV[1])
`
}
