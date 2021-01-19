// Copyright © 2021 Jacob Hansen. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package jobs

import "github.com/jake-hansen/jobs/consumers"

// Job represents a collection of workers that need to be scheduled.
type Job struct {
	Name          string
	Workers       *[]Worker
	DataConsumer  consumers.DataConsumer
	ErrorConsumer consumers.ErrorConsumer
	dataChannel   chan interface{}
	errorChannel  chan error
}

// NewJob creates a new job with the given name and given workers.
func NewJob(name string, workers *[]Worker) *Job {
	return &Job{
		Name:          name,
		Workers:       workers,
		DataConsumer:  consumers.DataPrinterConsumer{},
		ErrorConsumer: consumers.ErrorPrinterConsumer{},
		dataChannel:   make(chan interface{}),
		errorChannel:  make(chan error),
	}
}

// consumeData consumes the data channel for a Job.
func (j *Job) consumeData() {
	for val := range j.dataChannel {
		j.DataConsumer.Consume(val)
	}
}

// consumeErrors consumes the errors channel for a Job.
func (j *Job) consumeErrors() {
	for err := range j.errorChannel {
		j.ErrorConsumer.Consume(err)
	}
}
