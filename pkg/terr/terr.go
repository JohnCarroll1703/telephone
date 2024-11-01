package terr

import "errors"

var (
	RecordNotFound          = errors.New("record not found")
	RecordExists            = errors.New("record already exists")
	DatabaseConnectionError = errors.New("database connection error")
	ErrDbUnexpected         = errors.New("db unexpected error")
)
