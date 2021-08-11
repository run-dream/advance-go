package collector

import (
	"fmt"
	"metric/pkg/rolling"
	"sync"
	"time"
)

type RequestCollector struct {
	mutex *sync.RWMutex

	inflight    *rolling.Number
	totalReq    *rolling.Number
	passReq     *rolling.Number
	failReq     *rolling.Number
	runDuration *rolling.Timing
}

func newRequestMetricCollector(name string) MetricCollector {
	m := &RequestCollector{}
	m.mutex = &sync.RWMutex{}
	m.Reset()
	return m
}

func (d *RequestCollector) Inflight() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.inflight
}

func (d *RequestCollector) TotalReq() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.totalReq
}

func (d *RequestCollector) PassReq() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.passReq
}

func (d *RequestCollector) FailReq() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.failReq
}

func (d *RequestCollector) RunDuration() *rolling.Timing {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.runDuration
}

func (d *RequestCollector) Update(r MetricResult) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	d.inflight.Increment(r.Inflight)
	d.totalReq.Increment(r.TotalReq)
	d.passReq.Increment(r.PassReq)
	d.failReq.Increment(r.FailReq)

	d.runDuration.Add(r.RunDuration)
}

func (d *RequestCollector) Reset() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.inflight = rolling.NewNumber()
	d.totalReq = rolling.NewNumber()
	d.passReq = rolling.NewNumber()
	d.failReq = rolling.NewNumber()

	d.runDuration = rolling.NewTiming()
}

func (d *RequestCollector) ToString() string {
	return fmt.Sprintf("AvgInflight: %f, AvgTotalReq:%f, AvgPassReq:%f, AvgFailReq: %f, AvgRunDuration: %d",
		d.Inflight().Avg(time.Now()), d.TotalReq().Avg(time.Now()), d.PassReq().Avg(time.Now()),
		d.FailReq().Avg(time.Now()), d.RunDuration().Mean())
}
