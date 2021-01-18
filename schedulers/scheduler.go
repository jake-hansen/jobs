package schedulers

import (
	"fmt"
	"github.com/jake-hansen/jobs/consumers"
	"github.com/jake-hansen/jobs/jobs"
	"sync"
)

type Scheduler struct {
	dataChannel   chan interface{}
	errorChannel  chan error
	waitGroup     sync.WaitGroup
	DataConsumer  consumers.DataConsumer
	ErrorConsumer consumers.ErrorConsumer
}

func DefaultScheduler() *Scheduler {
	scheduler := &Scheduler{
		dataChannel:   nil,
		errorChannel:  nil,
		waitGroup:     sync.WaitGroup{},
		DataConsumer:  consumers.DataPrinterConsumer{},
		ErrorConsumer: consumers.ErrorPrinterConsumer{},
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

func (s *Scheduler) Schedule(job *jobs.Job) {
	if job != nil {
		s.dataChannel = make(chan interface{})
		s.errorChannel = make(chan error)

		for _, worker := range *job.Workers {
			s.waitGroup.Add(1)
			go runWorker(worker, s.dataChannel, s.errorChannel, &s.waitGroup)
		}

		go s.consumeData()
		go s.consumeErrors()
		go s.closeChannels()
	}
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

func (s *Scheduler) closeChannels() {
	s.waitGroup.Wait()
	close(s.dataChannel)
	close(s.errorChannel)
}

func (s *Scheduler) WaitForWorkers() {
	s.waitGroup.Wait()
}
