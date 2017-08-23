package server

import (
	"errors"
	"time"
)

type poolTaskWork func(interface{}) interface{}
type poolTaskCallback func(result interface{})

type poolTask struct {
	workFunc poolTaskWork
	data     interface{}

	result chan interface{}
}

type Pool struct {
	chs      chan *poolTask
	poolSize int
}

func NewPool(poolSize int) *Pool {
	pool := new(Pool)
	pool.chs = make(chan *poolTask, poolSize)

	for i := 0; i < poolSize; i++ {
		go func() {
			for task := range pool.chs {
				res := task.workFunc(task.data)
				if task.result != nil {
					task.result <- res
				}
			}
		}()
	}

	return pool
}

func (p *Pool) Size() int {
	return cap(p.chs)
}

func (p *Pool) AddTask(work poolTaskWork, data interface{}) (interface{}, error) {
	task := &poolTask{
		workFunc: work,
		data:     data,

		result: make(chan interface{}),
	}

	p.chs <- task
	return <-task.result, nil
}

func (p *Pool) AddTaskTimeout(work poolTaskWork, data interface{}, timeout time.Duration) (interface{}, error) {
	task := &poolTask{
		workFunc: work,
		data:     data,

		result: make(chan interface{}),
	}

	p.chs <- task

	select {
	case res := <-task.result:
		return res, nil
	case <-time.After(timeout):
		return nil, errors.New("cannot timeout")
	}
}

func (p *Pool) ThrowTask(simpleFunc func(interface{}), data interface{}) {
	task := &poolTask{
		workFunc: func(foo interface{}) interface{} {
			simpleFunc(foo)
			return nil
		},
		data: data,

		result: nil,
	}

	p.chs <- task
}
