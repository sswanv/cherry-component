package asynq

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

type Invoker interface {
	Register(controllers ...IController)
	EnqueueCtx(ctx context.Context, task Task, opts ...asynq.Option) error
}

func marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (c *Component) Register(controllers ...IController) {
	if c.serverMode() {
		c.server.controllers = append(c.server.controllers, controllers...)
	}
}

func (c *Component) EnqueueCtx(ctx context.Context, task Task, opts ...asynq.Option) error {
	if c.clientMode() {
		return c.client.enqueue(ctx, task, opts...)
	}
	return nil
}
