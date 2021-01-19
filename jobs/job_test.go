// Copyright © 2021 Jacob Hansen. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package jobs_test

import (
	"testing"

	"github.com/jake-hansen/jobs/jobs"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testJobName := "test job"

		worker := &testWorker{}

		var workers []jobs.Worker
		workers = append(workers, worker)

		job := jobs.NewJob(testJobName, &workers)

		assert.Equal(t, testJobName, job.Name)
		assert.Equal(t, &workers, job.Workers)
	})
}

type testWorker struct{}

func (t testWorker) Run() (interface{}, error) {
	return nil, nil
}

func (t testWorker) WorkerName() string {
	return "test worker"
}

func (t testWorker) GetPriority() interface{} {
	return 0
}
