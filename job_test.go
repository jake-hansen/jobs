// Copyright © 2021 Jacob Hansen. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package jobs_test

import (
	"github.com/jake-hansen/jobs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testJobName := "test job"
		var task jobs.Task = new(myTask)

		worker := jobs.Worker{
			Task:     &task,
			Name:     "test worker",
			Priority: nil,
		}

		var workers []jobs.Worker
		workers = append(workers, worker)

		job := jobs.NewJob(testJobName, &workers)

		assert.Equal(t, testJobName, job.Name)
		assert.Equal(t, &workers, job.Workers)
	})
}

type myTask struct{}

func (t myTask) Run() (interface{}, error) {
	panic("implement me")
}

