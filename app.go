package goqueue

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
