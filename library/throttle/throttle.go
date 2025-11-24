package throttle

import (
	"sync"
	"time"
)

type Caller interface {
	Call()
}

func New(limit time.Duration, call func()) Caller {
	return &executor{
		limit: limit,
		call:  call,
	}
}

type executor struct {
	limit time.Duration
	call  func()
	mutex sync.Mutex
	exist bool
	last  time.Time
}

func (exe *executor) Call() {
	exe.mutex.Lock()
	defer exe.mutex.Unlock()
	if exe.exist {
		return
	}
	interval := time.Since(exe.last)
	later := exe.limit - interval
	exe.exist = true
	time.AfterFunc(later, exe.execute)
}

func (exe *executor) execute() {
	exe.mutex.Lock()
	defer func() {
		exe.last = time.Now()
		exe.exist = false
		exe.mutex.Unlock()
	}()

	exe.call()
}
