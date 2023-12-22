package prometheus

import (
	"fmt"
	"net/http"

	cfacade "github.com/cherry-game/cherry/facade"
	clog "github.com/cherry-game/cherry/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	Name = "prometheus_component"
)

func NewComponent() *Component {
	return &Component{}
}

type Component struct {
	cfacade.Component
	config struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		Path string `json:"path"`
	}

	listeners []func()
}

func (c *Component) Init() {
	prometheusConfig := c.App().Settings().GetConfig("prometheus")
	if err := prometheusConfig.LastError(); err != nil {
		clog.Fatalf("get prometheus config err: %v", err)
	}
	if err := prometheusConfig.Unmarshal(&c.config); err != nil {
		clog.Fatalf("parse prometheus config err: %v", err)
	}
}

func (c *Component) OnAfterInit() {
	go c.run()
}

func (c *Component) Name() string {
	return Name
}

func (c *Component) run() {
	http.Handle(c.config.Path, promhttp.Handler())
	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	clog.Infof("Starting prometheus agent at %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		clog.Error(err)
	}
}

func (c *Component) OnStop() {
	for _, listener := range c.listeners {
		listener()
	}
}

func (c *Component) AddListener(fn func()) {
	c.listeners = append(c.listeners, fn)
}
