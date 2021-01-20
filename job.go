// Copyright © 2021 Jacob Hansen. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package jobs

import (
	"sync"

	"github.com/jake-hansen/jobs/consumers"
	"github.com/jake-hansen/jobs/utils"
)

// Job represents a collection of workers that need to be scheduled.
type Job struct {
	Name            string
	Workers         *[]Worker
	DataConsumer    consumers.DataConsumer
	ErrorConsumer   consumers.ErrorConsumer
	dataChannel     chan interface{}
	errorChannel    chan error
	workerWaitGroup *sync.WaitGroup
	jobWaitGroup    *sync.WaitGroup
	inProgress      utils.AtomicBool
}

// NewJob creates a new job with the given name and given workers.
func NewJob(name string, workers *[]Worker) *Job {
	job := &Job{
		Name:            name,
		Workers:         workers,
		DataConsumer:    consumers.DataPrinterConsumer{},
		ErrorConsumer:   consumers.ErrorPrinterConsumer{},
		dataChannel:     make(chan interface{}),
		errorChannel:    make(chan error),
		workerWaitGroup: new(sync.WaitGroup),
		jobWaitGroup:    new(sync.WaitGroup),
		inProgress:      utils.NewAtomicBool(false),
	}
	job.jobWaitGroup.Add(1)
	return job
}

// consumeData consumes the data channel for a Job.
func (j *Job) consumeData() {
	for val := range j.dataChannel {
		j.DataConsumer.Consume(val)
	}
	defer j.jobWaitGroup.Done()
}

// consumeErrors consumes the errors channel for a Job.
func (j *Job) consumeErrors() {
	for err := range j.errorChannel {
		j.ErrorConsumer.Consume(err)
	}
}

// waitForWorkers blocks until all Workers have finished executing.
func (j *Job) waitForWorkers() {
	if j.workerWaitGroup != nil {
		j.workerWaitGroup.Wait()
	}
}

// Wait blocks until the all the Job's Workers have finished executing and every data has been consumed by the
// DataConsumer.
func (j *Job) Wait() {
	if j.jobWaitGroup != nil {
		j.jobWaitGroup.Wait()
	}
}

// cleanup waits for all Workers to finish before closing the data and error channels.
func (j *Job) cleanup(job *Job) {
	j.waitForWorkers()
	close(job.dataChannel)
	close(job.errorChannel)
	job.inProgress.SafeSet(false)
}
