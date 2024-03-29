package rabbitmq

import (
	"encoding/json"
	"github.com/sswanv/cherry-component/rabbitmq/rabbitmq"
)

type RabbitMqSender struct {
	exchangeName string
	conf         rabbitmq.RabbitSenderConf
	sender       rabbitmq.Sender
}

func (r *RabbitMqSender) Send(routeKey string, msg any) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.sender.Send(r.exchangeName, routeKey, data)
}
