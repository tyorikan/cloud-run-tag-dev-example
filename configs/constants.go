package configs

type contextKey struct {
	name string
}

var (
	// LimitCtxKey is limit key of context value
	LimitCtxKey = &contextKey{"Limit"}
	// OffsetCtxKey is offset key of context value
	OffsetCtxKey = &contextKey{"Offset"}
)

const (
	// defaults
	DefaultEnv      = "dev"
	DefaultPort     = "8080"
	DefaultLogLevel = "debug"

	// DefaultLimit defines default number of items in the page
	DefaultLimit = 10
	// DefaultOffset defines default starting position.
	DefaultOffset = 0
	// DefaultTimeFormatLayout is default time format layout
	DefaultTimeFormatLayout = "2006-01-02T15:04:05Z"

	// UpperLimit defines maximum number of items in the page
	UpperLimit = 50

	// Env is set environment string
	Env = "ENV"
	// EnvPort is set server port
	EnvPort = "PORT"
	// EnvLogLevel is set logLevel
	EnvLogLevel = "LOG_LEVEL"
)
