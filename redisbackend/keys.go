package redisbackend

import "fmt"

type keyBuilder struct {
	namespace string
}

func newKeyBuilder(namespace string) keyBuilder {
	return keyBuilder{namespace: namespace}
}

func (b keyBuilder) readyStream(queue string) string {
	return fmt.Sprintf("%s:queue:%s:ready", b.namespace, queue)
}

func (b keyBuilder) scheduledSet(queue string) string {
	return fmt.Sprintf("%s:queue:%s:scheduled", b.namespace, queue)
}

func (b keyBuilder) deadLetterStream(queue string) string {
	return fmt.Sprintf("%s:queue:%s:dead", b.namespace, queue)
}

func (b keyBuilder) periodicDefinitionsHash() string {
	return fmt.Sprintf("%s:scheduler:periodic:definitions", b.namespace)
}

func (b keyBuilder) periodicDueSet() string {
	return fmt.Sprintf("%s:scheduler:periodic:due", b.namespace)
}

func (b keyBuilder) periodicLease(name string) string {
	return fmt.Sprintf("%s:scheduler:periodic:%s:lease", b.namespace, name)
}

func (b keyBuilder) message(taskID string) string {
	return fmt.Sprintf("%s:task:%s:message", b.namespace, taskID)
}

func (b keyBuilder) taskPrefix() string {
	return fmt.Sprintf("%s:task:", b.namespace)
}

func (b keyBuilder) state(taskID string) string {
	return fmt.Sprintf("%s:task:%s:state", b.namespace, taskID)
}

func (b keyBuilder) result(taskID string) string {
	return fmt.Sprintf("%s:task:%s:result", b.namespace, taskID)
}
