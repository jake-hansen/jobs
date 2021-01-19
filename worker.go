// Copyright © 2021 Jacob Hansen. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package jobs

// Worker represents an atomic task that needs to be executed. It can be thought of as a
// shell for a Task that contains metadata about that Task.
type Worker struct {
	Task     *Task
	Name     string
	Priority interface{}
	complete	bool
}

func NewWorker(task *Task, name string, priority interface{}) *Worker {
	worker := &Worker{
		Task:     task,
		Name:     name,
		Priority: priority,
		complete: false,
	}
	return worker
}

// Task represents a function that a Worker will perform.
type Task interface {
	Run() (interface{}, error)
}
