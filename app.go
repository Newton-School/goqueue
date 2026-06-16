package goqueue

import (
	"github.com/Newton-School/goqueue/inspect"
	"github.com/Newton-School/goqueue/producer"
	"github.com/Newton-School/goqueue/scheduler"
	"github.com/Newton-School/goqueue/worker"
	"github.com/Newton-School/goqueue/workflow"
)

// App is the root goqueue application instance.
type App struct {
	config   Config
	registry *TaskRegistry
}

// New creates a goqueue application with validated configuration.
func New(opts ...Option) (*App, error) {
	config := NewConfig(opts...)
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &App{
		config:   config,
		registry: NewTaskRegistry(),
	}, nil
}

// Config returns a copy of the application configuration.
func (a *App) Config() Config {
	return a.config
}

// RegisterTask registers a handler for name on this app.
func (a *App) RegisterTask(name TaskName, handler TaskHandler) error {
	return a.registry.Register(name, handler)
}

// LookupTask returns the registered handler for name.
func (a *App) LookupTask(name TaskName) (TaskHandler, error) {
	return a.registry.Lookup(name)
}

// TaskNames returns registered task names in sorted order.
func (a *App) TaskNames() []TaskName {
	return a.registry.Names()
}

// NewProducer creates a producer with app defaults.
func (a *App) NewProducer(opts ...producer.ProducerOption) (*producer.Producer, error) {
	backend, err := a.NewRedisBackend()
	if err != nil {
		return nil, err
	}

	allOpts := append(
		[]producer.ProducerOption{producer.WithProducerDefaultQueue(QueueName(a.config.DefaultQueue))},
		opts...,
	)

	return producer.NewProducer(backend, allOpts...)
}

// NewWorker creates a worker for queue-based task execution.
func (a *App) NewWorker(opts ...worker.WorkerOption) (*worker.Worker, error) {
	backend, err := a.NewRedisBackend()
	if err != nil {
		return nil, err
	}

	allOpts := append(
		[]worker.WorkerOption{worker.WithWorkerQueue(QueueName(a.config.DefaultQueue))},
		opts...,
	)

	return worker.NewWorker(backend, a.registry, allOpts...)
}

// NewScheduler creates a scheduler for periodic task dispatch.
func (a *App) NewScheduler(opts ...scheduler.SchedulerOption) (*scheduler.Scheduler, error) {
	backend, err := a.NewRedisBackend()
	if err != nil {
		return nil, err
	}

	allOpts := append(
		[]scheduler.SchedulerOption{scheduler.WithSchedulerDefaultQueue(QueueName(a.config.DefaultQueue))},
		opts...,
	)

	return scheduler.NewScheduler(backend, allOpts...)
}

// NewCanvas creates a workflow canvas for signatures, chains, groups, and chords.
func (a *App) NewCanvas(opts ...workflow.CanvasOption) (*workflow.Canvas, error) {
	backend, err := a.NewRedisBackend()
	if err != nil {
		return nil, err
	}

	allOpts := append(
		[]workflow.CanvasOption{workflow.WithCanvasDefaultQueue(QueueName(a.config.DefaultQueue))},
		opts...,
	)

	return workflow.NewCanvas(backend, allOpts...)
}

// NewInspector creates an inspection client bound to the configured backend.
func (a *App) NewInspector() (*inspect.Inspector, error) {
	backend, err := a.NewRedisBackend()
	if err != nil {
		return nil, err
	}

	return inspect.NewInspector(backend)
}
