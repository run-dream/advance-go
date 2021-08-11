package metric

import (
	"metric/pkg/collector"
	"sync"
	"time"
)

type Event struct {
	Type        string
	StartTime   time.Time
	RunDuration time.Duration
}

type MetricExchange struct {
	Name   string
	Events chan *Event
	Mutex  *sync.RWMutex

	collectors []collector.MetricCollector
}

func NewMetricExchange(name string) *MetricExchange {
	m := &MetricExchange{}
	m.Name = name

	m.Events = make(chan *Event, 2000)
	m.Mutex = &sync.RWMutex{}
	m.collectors = collector.Registry.InitializeMetricCollectors(name)
	m.Reset()

	go m.Monitor()

	return m
}

func (m *MetricExchange) Monitor() {
	for event := range m.Events {
		m.Mutex.RLock()

		wg := &sync.WaitGroup{}
		for _, collector := range m.collectors {
			wg.Add(1)
			go m.IncrementMetrics(wg, collector, event)
		}
		wg.Wait()

		m.Mutex.RUnlock()
	}
}

func (m *MetricExchange) IncrementMetrics(wg *sync.WaitGroup, collect collector.MetricCollector, event *Event) {
	r := collector.MetricResult{}
	switch event.Type {
	case "start":
		r.Inflight = 1
		r.TotalReq = 1
	case "success":
		r.Inflight = -1
		r.PassReq = 1
		r.RunDuration = event.RunDuration
	case "fail":
		r.Inflight = -1
		r.FailReq = 1
		r.RunDuration = event.RunDuration
	}

	collect.Update(r)

	wg.Done()
}

func (m *MetricExchange) Reset() {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	for _, collector := range m.collectors {
		collector.Reset()
	}
}

func (m *MetricExchange) Stats() string {
	result := ""
	for _, collector := range m.collectors {
		result += collector.ToString()
	}
	return result
}
