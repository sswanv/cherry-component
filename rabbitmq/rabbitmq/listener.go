package rabbitmq

import (
	"log"

	clog "github.com/cherry-game/cherry/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	ConsumeHandle func(message string) error

	ConsumeHandler interface {
		Consume(message []byte) error
	}

	RabbitListener struct {
		conn    *amqp.Connection
		channel *amqp.Channel
		forever chan bool
		handler ConsumeHandler
		queues  RabbitListenerConf
	}
)

func MustNewListener(listenerConf RabbitListenerConf, handler ConsumeHandler) MessageQueue {
	listener := RabbitListener{queues: listenerConf, handler: handler, forever: make(chan bool)}
	conn, err := amqp.Dial(getRabbitURL(listenerConf.RabbitConf))
	if err != nil {
		log.Fatalf("failed to connect rabbitmq, error: %v", err)
	}

	listener.conn = conn
	channel, err := listener.conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}

	listener.channel = channel
	return listener
}

func (q RabbitListener) Start() {
	for _, que := range q.queues.ListenerQueues {
		msg, err := q.channel.Consume(
			que.Name,
			"",
			que.AutoAck,
			que.Exclusive,
			que.NoLocal,
			que.NoWait,
			nil,
		)
		if err != nil {
			clog.Fatalf("failed to listener, error: %v", err)
		}

		go func(que ConsumerConf) {
			for d := range msg {
				if err := q.handler.Consume(d.Body); err != nil {
					clog.Errorf("Error on consuming: %s, error: %v", string(d.Body), err)
				} else if !que.AutoAck {
					err = d.Ack(false)
					if err != nil {
						clog.Errorf("ack err: %v", err)
					}
				}
			}
		}(que)
	}

	<-q.forever
}

func (q RabbitListener) Stop() {
	q.channel.Close()
	q.conn.Close()
	close(q.forever)
}
