/*
author @chelizichen
description 2024.4.10 设计一个协程池
思路：协程轮询监听是否有新任务进来，协程池用满了则进入事件队列里面等待
datetime 2024.5.8 采用协程复用的方法，避免大量协程的创建与销毁
*/
package pool

import (
	"context"
	"sync"
	"time"
)

// 携程锁
var lock = sync.Mutex{}

type WithSgridRoutinePoltFuncfunc func(*RoutinePool)

func WithSgriddRoutineErrHand(errHand func(err interface{})) WithSgridRoutinePoltFuncfunc {
	return func(c *RoutinePool) {
		c.errHand = errHand
	}
}

func WithSgriddRoutineMaxSize(maxSize int) WithSgridRoutinePoltFuncfunc {
	return func(c *RoutinePool) {
		c.maxSize = maxSize
		c.Jobs = make(chan Job, maxSize)
	}
}

type Job func()

type RoutinePool struct {
	maxSize  int                // 最大容量
	Jobs     chan Job           // job
	ctx      context.Context    // 取消协程用
	cancel   context.CancelFunc // 取消函数
	wg       sync.WaitGroup
	taskList []Job
	errHand  func(err interface{})
}

func NewRoutinePool(opt ...WithSgridRoutinePoltFuncfunc) *RoutinePool {
	ctx, cancel := context.WithCancel(context.Background()) // 可取消的
	p := &RoutinePool{
		ctx:    ctx,
		cancel: cancel,
	}
	for _, v := range opt {
		v(p)
	}
	if p.maxSize == 0 {
		panic("painc/error : missing params maxSize")
	}
	return p
}

func (p *RoutinePool) addNeedRunTask(job Job) {
	p.taskList = append(p.taskList, job)
}

func (p *RoutinePool) Add(job Job) {
	var newJob = func() {
		job()
		if len(p.taskList) > 0 {
			lock.Lock()
			size := len(p.taskList)
			needRun := p.taskList[0]
			p.taskList = p.taskList[1:size]
			lock.Unlock()
			needRun()
		}
	}
	select {
	case p.Jobs <- newJob:
	default:
		p.addNeedRunTask(job)
	}
}

func (p *RoutinePool) Stop() {
	p.cancel()
	p.wg.Wait()
}

func (p *RoutinePool) Run() {
	for {
		select {
		case <-p.ctx.Done():
			return
		case job, ok := <-p.Jobs:
			if !ok {
				return
			}
			p.wg.Add(1)
			go func(job Job) {
				defer p.wg.Done()
				defer func() {
					if err := recover(); err != nil {
						p.errHand(err)
					}
				}()
				job()
			}(job)
		default:
			// 如果没有任务，等待一小段时间再检查
			time.Sleep(100 * time.Millisecond)
			var size int = 0
			lock.Lock()
			if len(p.taskList) != p.maxSize {
				size = len(p.taskList)
			} else {
				size = p.maxSize
			}
			newRunList := p.taskList[0:size] // 截取
			for _, v := range newRunList {
				f := v
				p.Add(f)
			}
			p.taskList = p.taskList[size:]
			lock.Unlock()
		}
	}
}

// func main() {
// 	pool := NewRoutinePool(10)
// 	for i := 0; i < 100; i++ {
// 		pool.Add(func() {
// 			fmt.Println("i", i)
// 			// 模拟任务执行
// 		})
// 	}
// 	go func() {
// 		time.Sleep(time.Second * 3)
// 		pool.Stop()
// 	}()
// 	pool.Run()
// 	fmt.Println("cancel")
// }
