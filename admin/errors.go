package admin

import "errors"

var (
	ErrNilAdmin = errors.New("admin: admin is nil")
	ErrAdminBackend = errors.New("admin: admin backend is nil or not support control operations")
	ErrNilTaskID = errors.New("admin: task id is required")
	ErrInvalidQueue = errors.New("admin: queue name is required")
	ErrInvalidControlOption = errors.New("admin: invalid control options")
)

