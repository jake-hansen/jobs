package main

import (
	"fmt"
	"github.com/jake-hansen/jobs/jobs"
	"github.com/jake-hansen/jobs/schedulers"
	"math/rand"
	"time"
)

type MyWorker struct {
	duration time.Duration
	name     string
}

func main() {

	worker1 := NewWorker(5, "1")
	worker2 := NewWorker(2, "2")
	worker3 := NewWorker(1, "3")

	workers := make([]jobs.Worker, 0)
	workers = append(workers, worker1, worker2, worker3)

	job1 := jobs.NewJob("my job", &workers)

	scheduler := schedulers.DefaultScheduler()
	scheduler.Schedule(job1)
	scheduler.WaitForWorkers()
}

func NewWorker(seconds int, name string) *MyWorker {
	dur, _ := time.ParseDuration(fmt.Sprintf("%ds", seconds))
	return &MyWorker{duration: dur, name: name}
}

func (m *MyWorker) Run() (interface{}, error) {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(m.duration)
	return rand.Intn(100), nil
}

func (m *MyWorker) WorkerName() string {
	return m.name
}

func (m *MyWorker) GetPriority() interface{} {
	return 0
}

