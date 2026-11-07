package provider

import (
	"context"
	"errors"
	"time"
)

type retryable struct{ error }

func (r retryable) Unwrap() error { return r.error }

func (c *Client) Complete(ctx context.Context, request Request) (Message, Usage, error) {
	for attempt := 0; ; attempt++ {
		if err := ctx.Err(); err != nil {
			return Message{}, Usage{}, err
		}
		message, usage, err := c.complete(ctx, request)
		var transient retryable
		if err == nil || !errors.As(err, &transient) || attempt == 3 {
			return message, usage, err
		}
		timer := time.NewTimer(time.Second << attempt)
		select {
		case <-ctx.Done():
			timer.Stop()
			return Message{}, Usage{}, ctx.Err()
		case <-timer.C:
		}
	}
}
