package scheduler

import (
	"context"
	"fmt"
	"gameapp/dto"
	"gameapp/service/matchingservice"
	"github.com/go-co-op/gocron"
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	matchSvc matchingservice.Service
	sch      *gocron.Scheduler
}

func New(m matchingservice.Service) Scheduler {
	return Scheduler{
		matchSvc: m,
		sch:      gocron.NewScheduler(time.UTC),
	}
}

func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("start scheduler . . .")
	s.sch.Every(5).Second().Do(s.MatchWaitedUsers)
	s.sch.StartAsync()
	<-done

	// wait to finish job
	fmt.Println("stop scheduler..")
	s.sch.Stop()
}

func (s Scheduler) MatchWaitedUsers() {
	log.Println("MatchWaitedUsers started")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	users, err := s.matchSvc.MatchWaitedUsers(ctx, dto.MatchWaitedUsersRequest{})
	if err != nil {
		fmt.Println("matchSvc.MatchWaitedUsers error", err)
	}
	fmt.Println("users", users)
}
