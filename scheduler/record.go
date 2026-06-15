package scheduler

import (
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func (p PeriodicTask) toBackendRecord(defaultQueue task.QueueName, now time.Time) (backend.PeriodicTaskRecord, error) {
	normalized, err := p.Normalize(defaultQueue)
	if err != nil {
		return backend.PeriodicTaskRecord{}, err
	}

	return backend.PeriodicTaskRecord{
		Name:         normalized.Name.String(),
		TaskName:     normalized.TaskName,
		Queue:        normalized.Queue,
		Args:         copyAnySlice(normalized.Args),
		Kwargs:       copyAnyMap(normalized.Kwargs),
		Metadata:     copyStringMap(normalized.Metadata),
		ScheduleKind: backend.PeriodicScheduleInterval,
		Interval:     normalized.Schedule.Interval,
		StartAt:      normalized.StartAt,
		NextDueAt:    normalized.FirstDueAfter(now),
		Priority:     normalized.Priority,
		RetryPolicy:  normalized.RetryPolicy,
		UpdatedAt:    now.UTC(),
	}, nil
}

func periodicTaskFromBackendRecord(record backend.PeriodicTaskRecord) (PeriodicTask, error) {
	if err := record.Validate(); err != nil {
		return PeriodicTask{}, err
	}

	definition := PeriodicTask{
		Name:        PeriodicTaskName(record.Name),
		TaskName:    record.TaskName,
		Queue:       record.Queue,
		Args:        copyAnySlice(record.Args),
		Kwargs:      copyAnyMap(record.Kwargs),
		Metadata:    copyStringMap(record.Metadata),
		Schedule:    Every(record.Interval),
		StartAt:     record.StartAt,
		Priority:    record.Priority,
		RetryPolicy: record.RetryPolicy,
	}
	if err := definition.Validate(); err != nil {
		return PeriodicTask{}, err
	}

	return definition, nil
}
