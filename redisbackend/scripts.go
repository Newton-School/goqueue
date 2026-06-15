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

func moveDueScheduledScript() string {
	return `
local ids = redis.call('ZRANGEBYSCORE', KEYS[1], '-inf', ARGV[1], 'LIMIT', 0, ARGV[2])
local moved = {}
for _, id in ipairs(ids) do
  local messageKey = ARGV[3] .. id .. ':message'
  local message = redis.call('GET', messageKey)
  if message then
    local streamID = redis.call('XADD', KEYS[2], '*', 'id', id, 'message', message)
    redis.call('ZREM', KEYS[1], id)
    table.insert(moved, streamID)
    table.insert(moved, message)
  else
    redis.call('ZREM', KEYS[1], id)
  end
end
return moved
`
}

func markPeriodicDispatchedScript() string {
	return `
local token = redis.call('GET', KEYS[3])
if token ~= ARGV[1] then
  return 0
end
redis.call('HSET', KEYS[1], ARGV[2], ARGV[3])
redis.call('ZADD', KEYS[2], ARGV[4], ARGV[2])
redis.call('DEL', KEYS[3])
return 1
`
}

func advanceWorkflowChainScript() string {
	return `
local total = tonumber(redis.call('HGET', KEYS[1], 'total') or '-1')
if total < 1 then
  return {0, 0, ''}
end
local completedIndex = tonumber(redis.call('HGET', KEYS[1], 'completed_index') or '-1')
local dispatchedIndex = tonumber(redis.call('HGET', KEYS[1], 'dispatched_index') or '-1')
local step = tonumber(ARGV[1])
local nextStep = step + 1
if completedIndex >= step then
  return {0, 0, ''}
end
redis.call('HSET', KEYS[1], 'completed_index', step, 'completed_task_id', ARGV[2], 'completed_at', ARGV[3])
if nextStep >= total then
  return {1, 1, ''}
end
if dispatchedIndex >= nextStep then
  return {1, 0, ''}
end
local nextSignature = redis.call('HGET', KEYS[2], tostring(nextStep))
if not nextSignature then
  return {1, 0, ''}
end
redis.call('HSET', KEYS[1], 'dispatched_index', nextStep)
return {1, 0, nextSignature}
`
}
