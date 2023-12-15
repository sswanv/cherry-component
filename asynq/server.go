package asynq

import (
	"context"
	"reflect"

	cfacade "github.com/cherry-game/cherry/facade"

	clog "github.com/cherry-game/cherry/logger"

	"github.com/hibiken/asynq"
)

func newServer() *server {
	return &server{}
}

type server struct {
	app         cfacade.IApplication
	server      *asynq.Server
	controllers []IController
}

func (s *server) init(conf config, app cfacade.IApplication) {
	s.server = asynq.NewServer(asynq.RedisClientOpt{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	}, asynq.Config{
		Concurrency: 20,
		Queues: map[string]int{
			PriorityCritical: 6,
			PriorityDefault:  3,
			PriorityLow:      1,
		},
	})
	s.app = app
}

func (s *server) handle(rt reflect.Type) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, task *asynq.Task) error {
		rv := reflect.New(rt)
		if rv.CanAddr() {
			rv = rv.Addr()
		}
		value := rv.Interface()

		err := unmarshal(task.Payload(), value)
		if err != nil {
			return err
		}

		if t, ok := value.(Task); ok {
			return t.Process(ctx, s.app)
		}

		return nil
	}
}

func (s *server) run() {
	mux := asynq.NewServeMux()
	for _, controller := range s.controllers {
		for _, t := range controller.Tasks() {
			task, ok := t.(Task)
			if !ok {
				continue
			}

			name := task.Name()
			rt := reflect.TypeOf(t)
			if rt.Kind() == reflect.Ptr {
				rt = rt.Elem()
			}
			mux.HandleFunc(name, s.handle(rt))
		}
	}

	err := s.server.Run(mux)
	if err != nil {
		clog.Fatalf("run job err: %v", err)
	}
}
