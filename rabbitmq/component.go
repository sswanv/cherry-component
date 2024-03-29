package rabbitmq

import (
	cfacade "github.com/cherry-game/cherry/facade"
	clog "github.com/cherry-game/cherry/logger"
	cprofile "github.com/cherry-game/cherry/profile"
	"github.com/sswanv/cherry-component/rabbitmq/rabbitmq"
)

const (
	Name = "rabbitmq_component"
)

type config struct {
	Conf      rabbitmq.RabbitConf `json:"conf"`
	Exchanges map[string]struct {
		Exchange  rabbitmq.ExchangeConf              `json:"exchange"`
		Queues    map[string]rabbitmq.QueueConf      `json:"queues"`
		Bind      map[string][]string                `json:"bind"`
		Consumers map[string][]rabbitmq.ConsumerConf `json:"consumers"`
	}
}

func NewComponent(admin bool) *Component {
	return &Component{
		admin: admin,
	}
}

type Component struct {
	cfacade.Component

	admin  bool
	config config
}

func (c *Component) Name() string {
	return Name
}

func (c *Component) Init() {
	configData := cprofile.GetConfig("rabbitmq")
	if err := configData.LastError(); err != nil {
		clog.Fatalf("获取配置失败: %v", err)
	}
	if err := configData.Unmarshal(&c.config); err != nil {
		clog.Fatalf("获取配置失败: %v", err)
	}
}

func (c *Component) OnAfterInit() {
	if !c.admin {
		return
	}

	for name, conf := range c.config.Exchanges {
		exchangeConf := conf.Exchange
		admin := rabbitmq.MustNewAdmin(c.config.Conf)
		err := admin.DeclareExchange(exchangeConf, nil)
		if err != nil {
			clog.Fatalf("declare exchange err: %v", err)
		}

		for _, queue := range conf.Queues {
			err = admin.DeclareQueue(queue, nil)
			if err != nil {
				clog.Fatalf("declare queue err: %v", err)
			}
		}

		for routeKey, queues := range conf.Bind {
			for _, queue := range queues {
				err = admin.Bind(queue, routeKey, exchangeConf.ExchangeName, false, nil)
				if err != nil {
					clog.Fatalf("queue bind err: %v", err)
				}
			}
		}

		clog.Infof("exchange [%v] decleare successful", name)
	}
}
