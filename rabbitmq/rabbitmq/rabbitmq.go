package rabbitmq

type MessageQueue interface {
	Start()
	Stop()
}
