package metric

type StopListener interface {
	AddListener(fn func())
}
