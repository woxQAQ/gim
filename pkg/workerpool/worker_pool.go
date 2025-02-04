package workerpool

import (
	"context"
	"sync"
)

var (
	instance IWorkerPool
	once     sync.Once
)

func init() {
	once.Do(func() {
		instance = newWorkerPool(10, 1000)
	})
}

func GetInstance() IWorkerPool {
	return instance
}

// Task 表示一个异步任务
type Task func()

// IWorkerPool 定义工作协程池的接口
type IWorkerPool interface {
	// Start 启动工作协程池
	Start()
	// Stop 停止工作协程池
	Stop()
	// Submit 提交一个任务到工作协程池
	Submit(task Task)
}

var _ IWorkerPool = &workerPool{}

// workerPool 表示一个工作协程池
type workerPool struct {
	workers   int            // 工作协程数量
	taskQueue chan Task      // 任务队列
	wg        sync.WaitGroup // 用于等待所有工作协程完成
	ctx       context.Context
	cancel    context.CancelFunc
}

// newWorkerPool 创建一个新的工作协程池
func newWorkerPool(workers int, queueSize int) *workerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &workerPool{
		workers:   workers,
		taskQueue: make(chan Task, queueSize),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Start 启动工作协程池
func (p *workerPool) Start() {
	// 启动指定数量的工作协程
	once.Do(func() {
		for i := 0; i < p.workers; i++ {
			p.wg.Add(1)
			go p.worker()
		}
	})

}

// Stop 停止工作协程池
func (p *workerPool) Stop() {
	// 取消上下文
	p.cancel()
	// 关闭任务队列
	close(p.taskQueue)
	// 等待所有工作协程完成
	p.wg.Wait()
}

// Submit 提交一个任务到工作协程池
func (p *workerPool) Submit(task Task) {
	select {
	case <-p.ctx.Done():
		return
	case p.taskQueue <- task:
	}
}

// worker 工作协程的实现
func (p *workerPool) worker() {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case task, ok := <-p.taskQueue:
			if !ok {
				return
			}
			// 执行任务
			task()
		}
	}
}
