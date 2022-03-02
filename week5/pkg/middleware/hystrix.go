package middleware

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

//定义一个窗口对象
type Bucket struct {
	sync.Mutex
	//当前窗口请求数
	TotalConn int64
	//当前窗口失败数
	FatalConn int64
	//时间戳
	TimeStamp time.Time
}

func NewBucket() *Bucket {
	return &Bucket{
		TimeStamp: time.Now(),
	}
}

// 记录当前连接数
func (b *Bucket) Record(result bool) {
	if result {
		atomic.AddInt64(&b.TotalConn, 1)
	} else {
		atomic.AddInt64(&b.FatalConn, 1)
	}
}

//SlidingWindow 定义一个滑动窗口 ..储存 单个窗口对象
type SlidingWindow struct {
	//当前状态
	broken bool
	//储存列表
	buckets []*Bucket
	//滑动器大小
	size int
	//触发熔断总连接数
	reqThreshold float64
	//触发熔断总失败数
	failedThreshold float64
	//是否需要熔断

	//最后一次熔断时间
	lastTime time.Time
	//熔断cd
	brokenTimeGap time.Duration
	sync.RWMutex
}

//新建一个窗口对象
func NewSlidingWindow(
	size int,
	reqThreshold, failedThreshold float64,
	brokenbrokenTimeGap time.Duration,
) *SlidingWindow {
	return &SlidingWindow{
		size:            size,
		buckets:         make([]*Bucket, 0, size),
		reqThreshold:    reqThreshold,
		failedThreshold: failedThreshold,
		brokenTimeGap:   brokenbrokenTimeGap,
	}
}

//添加桶...根据窗口大小对桶切片进行裁剪
func (s *SlidingWindow) AddBucket() {
	s.Lock()
	defer s.Unlock()
	s.buckets = append(s.buckets, NewBucket())
	if len(s.buckets) > s.size {
		s.buckets = s.buckets[1:]
	}
}

//在当前桶记录连接结果
func (s *SlidingWindow) RecordReqResult(result bool) {
	if len(s.buckets) == 0 {
		s.AddBucket()
	}
	s.buckets[len(s.buckets)-1].Record(result)
}

// 显示所有桶连接状态
func (s *SlidingWindow) ShowAllBucket() {
	for _, v := range s.buckets {
		fmt.Printf("id:[%v],tatol:[%d],failed:[%d]", v.TimeStamp, v.TotalConn, v.FatalConn)
	}
}

//启动滑动窗口,时间窗口100ms
func (s *SlidingWindow) Start() {
	go func() {
		for {
			s.AddBucket()
			time.Sleep(time.Millisecond * 100)
		}
	}()
}

//根据当前桶的状态该.判断是否触发熔断
func (s *SlidingWindow) BreakJudgement() bool {
	s.Lock()
	defer s.Unlock()
	total, failed := int64(0), int64(0)
	for _, v := range s.buckets {
		total += v.TotalConn
		failed += v.FatalConn
	}
	if float64(failed)/float64(total) > s.reqThreshold && float64(total) > s.failedThreshold {

		return true
	}
	return false
}

//监控桶内所有失败连接数及熔断状态判断是否触发熔断
func (s *SlidingWindow) Monitor() {
	//先判断是否在熔断状态
	go func() {
		for {
			if s.broken {
				if s.IsBrokenTimeGap() {
					s.Lock()
					s.broken = false
					s.Unlock()
				}
				continue
			}
			if s.BreakJudgement() {
				s.Lock()
				s.broken = true
				s.lastTime = time.Now()
				s.Unlock()
			}
		}
	}()

	//判断是否需要开启熔断

}

//每隔1s展示当前熔断状态
func (s *SlidingWindow) Stauts() {
	go func() {
		for {
			fmt.Println(s.broken)
			time.Sleep(time.Second)
		}
	}()
}

//判断是否在熔断cd中
func (s *SlidingWindow) IsBrokenTimeGap() bool {
	return time.Since(s.lastTime) > s.brokenTimeGap
}

//获取当前熔断状态
func (s *SlidingWindow) IsBroken() bool {
	return s.broken
}
