package rabbitmq

import "fmt"

type RabbitConf struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	VHost    string `json:"v_host"`
}

type RabbitListenerConf struct {
	RabbitConf
	ListenerQueues []ConsumerConf `json:"listener_queues"`
}

type ConsumerConf struct {
	Name      string `json:"name"`
	AutoAck   bool   `json:"auto_ack"`
	Exclusive bool   `json:"exclusive"`
	NoLocal   bool   `json:"no_local"`
	NoWait    bool   `json:"no_wait"`
}

type RabbitSenderConf struct {
	RabbitConf
	ContentType string `json:"content_type"`
}

type QueueConf struct {
	Name       string `json:"name"`
	Durable    bool   `json:"durable"`
	AutoDelete bool   `json:"auto_delete"`
	Exclusive  bool   `json:"exclusive"`
	NoWait     bool   `json:"no_wait"`
}

type ExchangeConf struct {
	ExchangeName string `json:"exchange_name"`
	Type         string `json:"type"`
	Durable      bool   `json:"durable"`
	AutoDelete   bool   `json:"auto_delete"`
	Internal     bool   `json:"internal"`
	NoWait       bool   `json:"no_wait"`
}

func getRabbitURL(rabbitConf RabbitConf) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", rabbitConf.Username, rabbitConf.Password,
		rabbitConf.Host, rabbitConf.Port, rabbitConf.VHost)
}
