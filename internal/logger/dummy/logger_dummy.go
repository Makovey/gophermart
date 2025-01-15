package dummy

import "github.com/Makovey/gophermart/internal/logger"

type dummyLogger struct{}

func (d dummyLogger) Debug(format string, args ...any) {}
func (d dummyLogger) Info(format string, args ...any)  {}
func (d dummyLogger) Warn(format string, args ...any)  {}
func (d dummyLogger) Error(format string, args ...any) {}

func NewDummyLogger() logger.Logger {
	return &dummyLogger{}
}
