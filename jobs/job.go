// Copyright © 2021 Jacob Hansen. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package jobs

// Worker represents an atomic task that needs to be executed.
type Worker interface {
	Run() (interface{}, error)
	WorkerName() string
	GetPriority() interface{}
}

// Job represents a collection of workers that need to be scheduled.
type Job struct {
	Name    string
	Workers *[]Worker
}

// NewJob creates a new job with the given name and given workers.
func NewJob(name string, workers *[]Worker) *Job {
	return &Job{
		Name:    name,
		Workers: workers,
	}
}
