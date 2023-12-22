package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

type (
	CounterVec interface {
		Inc(labels ...string)
		Add(v float64, labels ...string)
		close() bool
	}

	CounterVecOpts struct {
		Namespace string
		Subsystem string
		Name      string
		Help      string
		Labels    []string
	}

	promCounterVec struct {
		counter *prom.CounterVec
	}
)

func NewCounterVec(cfg *CounterVecOpts, listener StopListener) CounterVec {
	if cfg == nil {
		return nil
	}

	vec := prom.NewCounterVec(prom.CounterOpts{
		Namespace: cfg.Namespace,
		Subsystem: cfg.Subsystem,
		Name:      cfg.Name,
		Help:      cfg.Help,
	}, cfg.Labels)
	prom.MustRegister(vec)
	cv := &promCounterVec{
		counter: vec,
	}
	listener.AddListener(func() {
		cv.close()
	})

	return cv
}

func (cv *promCounterVec) Inc(labels ...string) {
	cv.counter.WithLabelValues(labels...).Inc()
}

func (cv *promCounterVec) Add(v float64, labels ...string) {
	cv.counter.WithLabelValues(labels...).Add(v)
}

func (cv *promCounterVec) close() bool {
	return prom.Unregister(cv.counter)
}
