package collector

import (
	"sync"
	"time"
)

// Registry is the default metricCollectorRegistry that circuits will use to
// collect statistics about the health of the circuit.
var Registry = metricRegistry{
	lock: &sync.RWMutex{},
	registry: []func(name string) MetricCollector{
		newRequestMetricCollector,
	},
}

type metricRegistry struct {
	lock     *sync.RWMutex
	registry []func(name string) MetricCollector
}

func (m *metricRegistry) InitializeMetricCollectors(name string) []MetricCollector {
	m.lock.RLock()
	defer m.lock.RUnlock()

	metrics := make([]MetricCollector, len(m.registry))
	for i, metricCollectorInitializer := range m.registry {
		metrics[i] = metricCollectorInitializer(name)
	}
	return metrics
}

func (m *metricRegistry) Register(initMetricCollector func(string) MetricCollector) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.registry = append(m.registry, initMetricCollector)
}

type MetricResult struct {
	Inflight    float64 // 运行中请求数量
	TotalReq    float64 // 总请求数
	PassReq     float64 // 通过数量
	FailReq     float64 // 失败数量
	RunDuration time.Duration
}

type MetricCollector interface {
	Update(MetricResult)
	Reset()
	ToString() string
}
