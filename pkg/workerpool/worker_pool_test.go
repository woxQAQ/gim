package workerpool

import (
	"sync"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("WorkerPool", func() {
	var (
		pool *workerPool
	)

	BeforeEach(func() {
		pool = newWorkerPool(3, 10)
		pool.Start()
	})

	AfterEach(func() {
		pool.Stop()
	})

	Context("基本功能测试", func() {
		It("应该正确执行所有提交的任务", func() {
			var counter int32

			for i := 0; i < 5; i++ {
				pool.Submit(func() {
					atomic.AddInt32(&counter, 1)
				})
			}

			Eventually(func() int32 {
				return atomic.LoadInt32(&counter)
			}).Should(Equal(int32(5)))
		})
	})

	Context("并发安全性测试", func() {
		It("应该能安全地处理并发任务", func() {
			pool = newWorkerPool(5, 100)
			pool.Start()

			var mu sync.Mutex
			results := make(map[int]bool)
			var wg sync.WaitGroup

			for i := 0; i < 50; i++ {
				wg.Add(1)
				num := i
				pool.Submit(func() {
					defer wg.Done()
					time.Sleep(10 * time.Millisecond)
					mu.Lock()
					results[num] = true
					mu.Unlock()
				})
			}

			wg.Wait()
			Eventually(func() int { return len(results) }).Should(Equal(50))
		})
	})

	Context("队列满时的行为测试", func() {
		It("应该能处理队列满的情况", func() {
			pool = newWorkerPool(1, 2)
			pool.Start()

			done := make(chan struct{})

			pool.Submit(func() {
				time.Sleep(100 * time.Millisecond)
				close(done)
			})

			for i := 0; i < 5; i++ {
				pool.Submit(func() {
					time.Sleep(10 * time.Millisecond)
				})
			}

			Eventually(done).Should(BeClosed())
		})
	})

	Context("停止功能测试", func() {
		It("停止后不应接受新任务", func() {
			var counter int32

			for i := 0; i < 5; i++ {
				pool.Submit(func() {
					time.Sleep(50 * time.Millisecond)
					atomic.AddInt32(&counter, 1)
				})
			}

			pool.Stop()

			pool.Submit(func() {
				atomic.AddInt32(&counter, 1)
			})

			Eventually(func() int32 {
				return atomic.LoadInt32(&counter)
			}).Should(BeNumerically("<=", 5))
		})
	})

	Measure("工作协程池性能", func(b Benchmarker) {
		pool = newWorkerPool(5, 1000)
		pool.Start()
		defer pool.Stop()

		b.Time("并发任务执行", func() {
			done := make(chan struct{})
			for i := 0; i < 1000; i++ {
				pool.Submit(func() {
					time.Sleep(time.Microsecond)
				})
			}
			close(done)
			Eventually(done).Should(BeClosed())
		})
	}, 10)
})
