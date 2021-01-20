// Copyright © 2021 Jacob Hansen. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package jobs

import (
	"errors"
	"fmt"
	"sync"
)

// Scheduler manages scheduling and syncing Workers for a Job.
type Scheduler struct {
	Algorithm SchedulerAlgorithm
	Debug     bool
}

// DefaultScheduler creates a Scheduler with a DataPrinterConsumer and ErrorPrinterConsumer. It also includes the
// SequentialScheduler as the scheduling algorithm.
func DefaultScheduler() *Scheduler {
	scheduler := &Scheduler{
		Algorithm: SequentialScheduler{},
		Debug:     false,
	}

	return scheduler
}

// spawnWorker manages running a Worker and passes the Worker's return value and error to the appropriate channel.
func spawnWorker(worker Worker, dataChannel chan interface{}, errorChannel chan error, wg *sync.WaitGroup, debug bool) {
	if debug {
		fmt.Printf("[DEBUG] starting worker %s\n", worker.Name)
	}
	defer wg.Done()
	val, err := (*worker.Task).Run()
	dataChannel <- val
	errorChannel <- err
	if debug {
		fmt.Printf("[DEBUG] ended worker %s\n", worker.Name)
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
	Schedule(workers *[]Worker) *[]Worker
}

// SequentialScheduler is a SchedulerAlgorithm that schedules Workers in the order in which they appear.
type SequentialScheduler struct{}

// Schedule returns an unmodified slice of the given Workers.
func (s SequentialScheduler) Schedule(workers *[]Worker) *[]Worker {
	return workers
}

// SubmitJob manages running a Job. When executed, SubmitJob begins spawning Workers as picked by the SchedulerAlgorithm.
func (s *Scheduler) SubmitJob(job *Job) error {
	if job != nil {
		job.inProgress.SafeSet(true)

		for _, worker := range *s.Algorithm.Schedule(job.Workers) {
			job.workerWaitGroup.Add(1)
			go spawnWorker(worker, job.dataChannel, job.errorChannel, job.workerWaitGroup, s.Debug)
		}

		go job.consumeData()
		go job.consumeErrors()
		go job.cleanup(job)
	} else {
		return errors.New("scheduler: job cannot be nil")
	}
	return nil
}
