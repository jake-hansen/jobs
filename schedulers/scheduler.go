package schedulers

import (
	"errors"
	"fmt"
	"sync"

	"github.com/jake-hansen/jobs/consumers"
	"github.com/jake-hansen/jobs/jobs"
)

// Scheduler manages scheduling and syncing Workers for a Job.
type Scheduler struct {
	dataChannel   chan interface{}
	errorChannel  chan error
	waitGroup     *sync.WaitGroup
	DataConsumer  consumers.DataConsumer
	ErrorConsumer consumers.ErrorConsumer
	Algorithm     SchedulerAlgorithm
	jobInProgress bool
	Debug         bool
}

// DefaultScheduler creates a Scheduler with a DataPrinterConsumer and ErrorPrinterConsumer. It also includes the
// SequentialScheduler as the scheduling algorithm.
func DefaultScheduler() *Scheduler {
	scheduler := &Scheduler{
		dataChannel:   nil,
		errorChannel:  nil,
		waitGroup:     nil,
		DataConsumer:  consumers.DataPrinterConsumer{},
		ErrorConsumer: consumers.ErrorPrinterConsumer{},
		Algorithm:     SequentialScheduler{},
		Debug:         false,
	}

	return scheduler
}

// runWorker manages running a Worker and passes the Worker's return value and error to the appropriate channel.
func runWorker(worker jobs.Worker, dataChannel chan interface{}, errorChannel chan error, wg *sync.WaitGroup, debug bool) {
	if debug {
		fmt.Printf("[DEBUG] starting worker %s\n", worker.WorkerName())
	}
	defer wg.Done()
	val, err := worker.Run()
	dataChannel <- val
	errorChannel <- err
	if debug {
		fmt.Printf("[DEBUG] ended worker %s\n", worker.WorkerName())
	}

}

// SchedulerAlgorithm defines an algorithm for scheduling Workers. An algorithm should return a pointer to a slice
// with the given Workers in the wanted order. Workers will be started sequentially using this returned slice.
//
// For example, if the given Workers slice looks like this
//
//					[ w1, w2, w3, w4, w5 ]
//
// a reordered slice could be returned that looks like this
//
//					[ w5, w3, w1, w2, w4 ]
//
// Then, when the job that contains these Workers is scheduled, the Workers will be scheduled as w5-> w3-> w1-> w2-> w4.
type SchedulerAlgorithm interface {
	Schedule(workers *[]jobs.Worker) *[]jobs.Worker
}

// SequentialScheduler is a SchedulerAlgorithm that schedules Workers in the order in which they appear.
type SequentialScheduler struct{}

// Schedule returns an unmodified slice of the given Workers.
func (s SequentialScheduler) Schedule(workers *[]jobs.Worker) *[]jobs.Worker {
	return workers
}

// Schedule manages running a Job. When executed, Schedule begins spawning Workers as picked by the SchedulerAlgorithm.
func (s *Scheduler) Schedule(job *jobs.Job) error {
	if job != nil {
		if !s.jobInProgress {
			s.waitGroup = new(sync.WaitGroup)
			s.jobInProgress = true
			s.dataChannel = make(chan interface{})
			s.errorChannel = make(chan error)

			for _, worker := range *s.Algorithm.Schedule(job.Workers) {
				s.waitGroup.Add(1)
				go runWorker(worker, s.dataChannel, s.errorChannel, s.waitGroup, s.Debug)
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

// consumeData consumes the data channel for a Scheduler.
func (s *Scheduler) consumeData() {
	for val := range s.dataChannel {
		s.DataConsumer.Consume(val)
	}
}

// consumeErrors consumes the errors channel for a Scheduler.
func (s *Scheduler) consumeErrors() {
	for err := range s.errorChannel {
		s.ErrorConsumer.Consume(err)
	}
}

// cleanup waits for all Workers to finish before closing the data and error channels.
func (s *Scheduler) cleanup() {
	s.waitGroup.Wait()
	close(s.dataChannel)
	close(s.errorChannel)
	s.jobInProgress = false
}

// WaitForWorkers blocks until all Workers have finished processing.
func (s *Scheduler) WaitForWorkers() {
	if s.waitGroup != nil {
		s.waitGroup.Wait()
	}
}
