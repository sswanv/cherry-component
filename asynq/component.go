package asynq

import (
	"context"
	"fmt"
	"reflect"

	clog "github.com/cherry-game/cherry/logger"

	cfacade "github.com/cherry-game/cherry/facade"
	"github.com/hibiken/asynq"
	_ "github.com/hibiken/asynq"
)

const (
	Name = "asynq_component"
)

const (
	PriorityCritical = "critical"
	PriorityDefault  = "default"
	PriorityLow      = "low"
)

type Mode byte

const (
	ModeNone   Mode = 1 << iota // None
	ModeServer                  // 服务端
	ModeClient                  // 客户端
)

type (
	Task interface {
		Name() string
		Opts() []asynq.Option
		Process(ctx context.Context, app cfacade.IApplication) error
	}
	TaskInfo struct {
		Name string       // 任务名
		Type reflect.Type // 类型信息
	}
)

type IController interface {
	PreInit(app cfacade.IApplication)
	Tasks() []Task
}

func ComponentName(mode Mode) string {
	switch mode {
	case ModeServer:
		return fmt.Sprintf("%s_%s", Name, "server")
	case ModeClient:
		return fmt.Sprintf("%s_%s", Name, "client")
	default:
		return Name
	}
}

func NewComponent(mode Mode) *Component {
	c := &Component{
		mode: mode,
	}
	if c.serverMode() {
		c.server = newServer()
	}
	if c.clientMode() {
		c.client = newClient()
	}
	return c
}

type config struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type Component struct {
	cfacade.Component
	config config

	mode   Mode
	server *server
	client *client
}

func (c *Component) Name() string {
	return ComponentName(c.mode)
}

func (c *Component) Init() {
	asynqConf := c.App().Settings().GetConfig("asynq")
	if err := asynqConf.LastError(); err != nil {
		clog.Fatalf("get asynq config err: %v", err)
	}
	if err := asynqConf.Unmarshal(&c.config); err != nil {
		clog.Fatalf("asynq config unmarshal err: %v", err)
	}

	if c.serverMode() {
		c.server.init(c.config, c.App())
	}
	if c.clientMode() {
		c.client.init(c.config, c.App())
	}
}

func (c *Component) OnAfterInit() {
	go c.Run()
}

func (c *Component) Run() {
	if c.serverMode() {
		c.server.run()
	}
}

func (c *Component) serverMode() bool {
	return c.mode&ModeServer == ModeServer
}

func (c *Component) clientMode() bool {
	return c.mode&ModeClient == ModeClient
}
