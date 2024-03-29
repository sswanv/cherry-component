package rabbitmq

import (
	clog "github.com/cherry-game/cherry/logger"
	"github.com/sswanv/cherry-component/rabbitmq/rabbitmq"
)

type (
	MessageQueue interface {
		rabbitmq.MessageQueue
	}
	Sender interface {
		Send(routeKey string, msg any) error
	}
)

type Invoker interface {
	MustNewListener(exchangeName string, consumerName string, handler rabbitmq.ConsumeHandler) MessageQueue
	MustNewSender(exchangeName string, contentTypes ...string) Sender
}

func (c *Component) MustNewListener(exchangeName string, consumerName string, handler rabbitmq.ConsumeHandler) MessageQueue {
	conf, ok := c.config.Exchanges[exchangeName]
	if !ok {
		clog.Fatalf("exchange conf [%v] not found", exchangeName)
	}

	consumers, ok := conf.Consumers[consumerName]
	if !ok {
		clog.Fatalf("consumer conf [%v] not found", consumerName)
	}

	var listenerQueues []rabbitmq.ConsumerConf
	for _, consumer := range consumers {
		if _, ok := conf.Queues[consumer.Name]; !ok {
			clog.Fatalf("queue [%v] not found", consumer.Name)
		}
		listenerQueues = append(listenerQueues, consumer)
	}

	return rabbitmq.MustNewListener(rabbitmq.RabbitListenerConf{
		RabbitConf:     c.config.Conf,
		ListenerQueues: listenerQueues,
	}, handler)
}

func (c *Component) MustNewSender(exchangeName string, contentTypes ...string) Sender {
	contentType := "application/json"
	if len(contentTypes) > 1 {
		contentType = contentTypes[0]
	}

	conf := rabbitmq.RabbitSenderConf{
		RabbitConf:  c.config.Conf,
		ContentType: contentType,
	}

	return &RabbitMqSender{
		exchangeName: exchangeName,
		conf:         conf,
		sender:       rabbitmq.MustNewSender(conf),
	}
}
