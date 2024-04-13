/*
author @chelizichen
description 2024.4.10 设计一个协程池
思路：协程轮询监听是否有新任务进来，协程池用满了则进入事件队列里面等待
*/
package pool

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Job func()

type RoutinePool struct {
	maxSize      int                // 最大容量
	Jobs         chan Job           // job
	ctx          context.Context    // 取消协程用
	cancel       context.CancelFunc // 取消函数
	wg           sync.WaitGroup
	runningCount int32 // 目前正在跑的
	taskList     []Job
}

func NewRoutinePool(maxSize int) *RoutinePool {
	ctx, cancel := context.WithCancel(context.Background()) // 可取消的
	return &RoutinePool{
		maxSize: maxSize,
		Jobs:    make(chan Job, maxSize),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (p *RoutinePool) addNeedRunTask(job Job) {
	p.taskList = append(p.taskList, job)
}

func (p *RoutinePool) Add(job Job) {
	select {
	case p.Jobs <- job:
		atomic.AddInt32(&p.runningCount, 1)
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
					atomic.AddInt32(&p.runningCount, -1)
				}()
				job()
			}(job)
		default:
			// 如果没有任务，等待一小段时间再检查
			time.Sleep(100 * time.Millisecond)
			var size int = 0
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
		}
	}
}

func (p *RoutinePool) RunningCount() int {
	return int(atomic.LoadInt32(&p.runningCount))
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
