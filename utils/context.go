package utils

import "context"

// ContextWait waits for context to be done.
func ContextWait(ctx context.Context) bool {
	for {
		select {
		case <-ctx.Done():
			return true
		default:
			return false
		}
	}
}
