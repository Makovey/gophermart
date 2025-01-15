package accrual

import (
	"fmt"
	"time"
)

type ManyRequestError struct {
	RetryAfter time.Duration
}

func (e *ManyRequestError) Error() string {
	return fmt.Sprintf("got too many request, retry after - %s", e.RetryAfter)
}
