package extensions

// Logger is the logging interface used by all confignet components.
// Implement this interface to integrate your preferred logging library (zap, logrus, slog, etc.).
// Register your implementation via confignet.SetLogger.
type Logger interface {
	Printf(format string, args ...interface{})
}
