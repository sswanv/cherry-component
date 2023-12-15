package asynq

import (
	"context"

	cfacade "github.com/cherry-game/cherry/facade"
	"github.com/hibiken/asynq"
)

func newClient() *client {
	return &client{}
}

type client struct {
	app    cfacade.IApplication
	client *asynq.Client
}

func (c *client) init(conf config, app cfacade.IApplication) {
	c.client = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})
	c.app = app
}

func (c *client) enqueue(ctx context.Context, task Task, opts ...asynq.Option) error {
	payload, err := marshal(task)
	if err != nil {
		return err
	}
	if len(opts) == 0 {
		opts = task.Opts()
	}

	t := asynq.NewTask(task.Name(), payload, opts...)
	_, err = c.client.EnqueueContext(ctx, t)
	if err != nil {
		return err
	}
	return nil
}
