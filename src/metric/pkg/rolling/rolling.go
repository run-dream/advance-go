package rolling

import (
	"sync"
	"time"
)

// WINDOW_SIZE hytrix-go 的窗口大小固定为 10s
const WINDOW_SIZE = 10

// Number 滑动窗口
type Number struct {
	// 计数桶 key 为时间戳， value 为值
	Buckets map[int64]*numberBucket
	// 读写锁
	Mutex *sync.RWMutex
}

// numberBucket 值, 这里使用结构体的作用可能是方便扩展?
type numberBucket struct {
	Value float64
}

// NewNumber 构造器
func NewNumber() *Number {
	r := &Number{
		Buckets: make(map[int64]*numberBucket),
		Mutex:   &sync.RWMutex{},
	}
	return r
}

// getCurrentBucket 获取当前的 bucket 值， 若不存在则新建
func (r *Number) getCurrentBucket() *numberBucket {
	now := time.Now().Unix()
	var bucket *numberBucket
	var ok bool

	if bucket, ok = r.Buckets[now]; !ok {
		bucket = &numberBucket{}
		r.Buckets[now] = bucket
	}

	return bucket
}

// removeOldBuckets 删除过期的桶, 本质是用 map 模拟 FIFO 队列
func (r *Number) removeOldBuckets() {
	now := time.Now().Unix() - WINDOW_SIZE
	for timestamp := range r.Buckets {
		if timestamp <= now {
			delete(r.Buckets, timestamp)
		}
	}
}

// Increment 增加当前的 bucket 值
func (r *Number) Increment(i float64) {
	if i == 0 {
		return
	}

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	b := r.getCurrentBucket()
	b.Value += i
	r.removeOldBuckets()
}

// UpdateMax 更新当前的 bucket 值
func (r *Number) UpdateMax(n float64) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	b := r.getCurrentBucket()
	if n > b.Value {
		b.Value = n
	}
	r.removeOldBuckets()
}

// Sum 获取窗口内的 bucket 总和
func (r *Number) Sum(now time.Time) float64 {
	sum := float64(0)

	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	for timestamp, bucket := range r.Buckets {
		if timestamp >= now.Unix()-WINDOW_SIZE {
			sum += bucket.Value
		}
	}

	return sum
}

// Max 获取窗口内的 bucket 最大值
func (r *Number) Max(now time.Time) float64 {
	var max float64

	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	for timestamp, bucket := range r.Buckets {
		// TODO: configurable rolling window
		if timestamp >= now.Unix()-10 {
			if bucket.Value > max {
				max = bucket.Value
			}
		}
	}

	return max
}

// Avg 获取窗口内的 bucket 平均值
func (r *Number) Avg(now time.Time) float64 {
	return r.Sum(now) / WINDOW_SIZE
}
