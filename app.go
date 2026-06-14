package goqueue

// App is the root goqueue application instance.
type App struct {
	config Config
}

// New creates a goqueue application with validated configuration.
func New(opts ...Option) (*App, error) {
	config := NewConfig(opts...)
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &App{config: config}, nil
}

// Config returns a copy of the application configuration.
func (a *App) Config() Config {
	return a.config
}
