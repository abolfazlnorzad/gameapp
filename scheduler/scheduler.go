package scheduler

import (
	"fmt"
	"time"
)

type Scheduler struct {
}

func New() Scheduler {
	return Scheduler{}
}

func (s Scheduler) Start(done <-chan bool) {
	fmt.Println("start scheduler . . .")
	for {
		select {
		case d := <-done:
			fmt.Println("done case : ", d)
		default:
			now := time.Now()
			fmt.Println("scheduler now", now)
			time.Sleep(3 * time.Second)
		}
	}
}
