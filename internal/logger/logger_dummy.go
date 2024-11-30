package logger

type dummyLogger struct{}

func (d dummyLogger) Debug(format string, args ...any) {}
func (d dummyLogger) Info(format string, args ...any)  {}
func (d dummyLogger) Warn(format string, args ...any)  {}
func (d dummyLogger) Error(format string, args ...any) {}

func NewDummyLogger() Logger {
	return &dummyLogger{}
}
