package ctx

import (
	"context"
	"time"
)

type ContextManager struct {
	timeout time.Duration
}

func NewContextManager(timeout time.Duration) *ContextManager {
	return &ContextManager{
		timeout: timeout,
	}
}

func (cm *ContextManager) WithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), cm.timeout)
}
