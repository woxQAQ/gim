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
	once      sync.Once  // 实例级别的 once
	isClosed  bool       // 标记通道是否已关闭
	mu        sync.Mutex // 保护 isClosed
}

// newWorkerPool 创建一个新的工作协程池
func newWorkerPool(workers int, queueSize int) *workerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &workerPool{
		workers:   workers,
		taskQueue: make(chan Task, queueSize),
		ctx:       ctx,
		cancel:    cancel,
		isClosed:  false,
	}
}

// Start 启动工作协程池
func (p *workerPool) Start() {
	// 启动指定数量的工作协程
	p.once.Do(func() {
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
	// 等待所有工作协程完成
	p.wg.Wait()
	// 安全关闭任务队列
	p.mu.Lock()
	if !p.isClosed {
		close(p.taskQueue)
		p.isClosed = true
	}
	p.mu.Unlock()
}

// Submit 提交一个任务到工作协程池
func (p *workerPool) Submit(task Task) {
	p.mu.Lock()
	if p.isClosed {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	// 使用单个select语句避免嵌套select可能导致的死锁
	select {
	case <-p.ctx.Done():
		return
	case p.taskQueue <- task:
		// 任务已成功提交
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
