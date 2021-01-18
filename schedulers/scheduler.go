package schedulers

import (
	"errors"
	"fmt"
	"github.com/jake-hansen/jobs/consumers"
	"github.com/jake-hansen/jobs/jobs"
	"sync"
)

type Scheduler struct {
	dataChannel   chan interface{}
	errorChannel  chan error
	waitGroup     *sync.WaitGroup
	DataConsumer  consumers.DataConsumer
	ErrorConsumer consumers.ErrorConsumer
	Algorithm     SchedulerAlgorithm
	jobInProgress bool
}

func DefaultScheduler() *Scheduler {
	scheduler := &Scheduler{
		dataChannel:   nil,
		errorChannel:  nil,
		waitGroup:     nil,
		DataConsumer:  consumers.DataPrinterConsumer{},
		ErrorConsumer: consumers.ErrorPrinterConsumer{},
		Algorithm:     SequentialScheduler{},
	}

	return scheduler
}

func runWorker(worker jobs.Worker, dataChannel chan interface{}, errorChannel chan error, wg *sync.WaitGroup) {
	fmt.Printf("[DEBUG] starting worker %s\n", worker.WorkerName())
	defer wg.Done()
	val, err := worker.Run()
	dataChannel <- val
	errorChannel <- err
	fmt.Printf("[DEBUG] ended worker %s\n", worker.WorkerName())

}

type SchedulerAlgorithm interface {
	Schedule(workers *[]jobs.Worker) *[]jobs.Worker
}

type SequentialScheduler struct{}

func (s SequentialScheduler) Schedule(workers *[]jobs.Worker) *[]jobs.Worker {
	return workers
}

func (s *Scheduler) Schedule(job *jobs.Job) error {
	if job != nil {
		if !s.jobInProgress {
			s.waitGroup = new(sync.WaitGroup)
			s.jobInProgress = true
			s.dataChannel = make(chan interface{})
			s.errorChannel = make(chan error)

			for _, worker := range *s.Algorithm.Schedule(job.Workers) {
				s.waitGroup.Add(1)
				go runWorker(worker, s.dataChannel, s.errorChannel, s.waitGroup)
			}

			go s.consumeData()
			go s.consumeErrors()
			go s.cleanup()
		} else {
			return fmt.Errorf("scheduler: cannot schedule job [%s]. A job already in progress", job.Name)
		}
	} else {
		return errors.New("scheduler: job cannot be nil")
	}
	return nil
}

func (s *Scheduler) consumeData() {
	for val := range s.dataChannel {
		s.DataConsumer.Consume(val)
	}
}

func (s *Scheduler) consumeErrors() {
	for err := range s.errorChannel {
		s.ErrorConsumer.Consume(err)
	}
}

func (s *Scheduler) cleanup() {
	s.waitGroup.Wait()
	close(s.dataChannel)
	close(s.errorChannel)
	s.jobInProgress = false
}

func (s *Scheduler) WaitForWorkers() {
	if s.waitGroup != nil {
		s.waitGroup.Wait()
	}
}
