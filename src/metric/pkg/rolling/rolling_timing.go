package rolling

import (
	"math"
	"sort"
	"sync"
	"time"
)

// TIMING_SIZE 固定为 60s
const TIMING_SIZE = 60

// Timing 每个时间内，多个时间段
type Timing struct {
	Buckets map[int64]*timingBucket
	Mutex   *sync.RWMutex

	CachedSortedDurations []time.Duration
	LastCachedTime        int64
}

type timingBucket struct {
	Durations []time.Duration
}

// NewTiming creates a RollingTiming struct.
func NewTiming() *Timing {
	r := &Timing{
		Buckets: make(map[int64]*timingBucket),
		Mutex:   &sync.RWMutex{},
	}
	return r
}

type byDuration []time.Duration

func (c byDuration) Len() int           { return len(c) }
func (c byDuration) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c byDuration) Less(i, j int) bool { return c[i] < c[j] }

// SortedDurations 排序
func (r *Timing) SortedDurations() []time.Duration {
	r.Mutex.RLock()
	t := r.LastCachedTime
	r.Mutex.RUnlock()

	if t+time.Duration(1*time.Second).Nanoseconds() > time.Now().UnixNano() {
		return r.CachedSortedDurations
	}

	var durations byDuration
	now := time.Now()

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	for timestamp, b := range r.Buckets {
		if timestamp >= now.Unix()-TIMING_SIZE {
			for _, d := range b.Durations {
				durations = append(durations, d)
			}
		}
	}

	sort.Sort(durations)

	r.CachedSortedDurations = durations
	r.LastCachedTime = time.Now().UnixNano()

	return r.CachedSortedDurations
}

func (r *Timing) getCurrentBucket() *timingBucket {
	r.Mutex.RLock()
	now := time.Now()
	bucket, exists := r.Buckets[now.Unix()]
	r.Mutex.RUnlock()

	if !exists {
		r.Mutex.Lock()
		defer r.Mutex.Unlock()

		r.Buckets[now.Unix()] = &timingBucket{}
		bucket = r.Buckets[now.Unix()]
	}

	return bucket
}

func (r *Timing) removeOldBuckets() {
	now := time.Now()

	for timestamp := range r.Buckets {
		if timestamp <= now.Unix()-TIMING_SIZE {
			delete(r.Buckets, timestamp)
		}
	}
}

func (r *Timing) Add(duration time.Duration) {
	b := r.getCurrentBucket()

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	b.Durations = append(b.Durations, duration)
	r.removeOldBuckets()
}

// Percentile n 分位
func (r *Timing) Percentile(p float64) uint32 {
	sortedDurations := r.SortedDurations()
	length := len(sortedDurations)
	if length <= 0 {
		return 0
	}

	pos := r.ordinal(len(sortedDurations), p) - 1
	return uint32(sortedDurations[pos].Nanoseconds() / 1000000)
}

func (r *Timing) ordinal(length int, percentile float64) int64 {
	if percentile == 0 && length > 0 {
		return 1
	}

	return int64(math.Ceil((percentile / float64(100)) * float64(length)))
}

// Mean 平均值
func (r *Timing) Mean() uint32 {
	sortedDurations := r.SortedDurations()
	var sum time.Duration
	for _, d := range sortedDurations {
		sum += d
	}

	length := int64(len(sortedDurations))
	if length == 0 {
		return 0
	}

	return uint32(sum.Nanoseconds()/length) / 1000000
}
