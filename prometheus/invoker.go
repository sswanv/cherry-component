package prometheus

import (
	"github.com/sswanv/cherry-component/prometheus/metric"
)

type Invoker interface {
	NewCounterVec(cfg *metric.CounterVecOpts) metric.CounterVec
	NewGaugeVec(cfg *metric.GaugeVecOpts) metric.GaugeVec
}

func (c *Component) NewCounterVec(cfg *metric.CounterVecOpts) metric.CounterVec {
	return metric.NewCounterVec(cfg, c)
}

func (c *Component) NewGaugeVec(cfg *metric.GaugeVecOpts) metric.GaugeVec {
	return metric.NewGaugeVec(cfg, c)
}
