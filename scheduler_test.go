// Copyright © 2021 Jacob Hansen. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package jobs_test

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/jake-hansen/jobs"

	"github.com/stretchr/testify/assert"
)

func TestDefaultScheduler(t *testing.T) {
	testScheduler := jobs.DefaultScheduler()

	assert.Equal(t, jobs.SequentialScheduler{}, testScheduler.Algorithm)
}

func spawnWorkers() *[]jobs.Worker {
	var workers []jobs.Worker
	for i := 0; i < 100; i++ {
		var integerTask jobs.Task = new(integerTask)
		workers = append(workers, jobs.Worker{
			Task:     &integerTask,
			Name:     strconv.Itoa(i),
			Priority: nil,
		})
	}
	return &workers
}

func TestScheduler_Schedule(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testWorkers := spawnWorkers()

		testJob := jobs.NewJob("test job", testWorkers)
		testScheduler := jobs.DefaultScheduler()
		testScheduler.Debug = true
		consumer := additionConsumer{Sum: 0}
		testJob.DataConsumer = &consumer
		err := testScheduler.SubmitJob(testJob)
		testJob.Wait()

		assert.NoError(t, err)
		assert.Equal(t, len(*testWorkers)*5, consumer.Sum)
	})

	t.Run("failure-nil-job", func(t *testing.T) {
		testScheduler := jobs.DefaultScheduler()
		err := testScheduler.SubmitJob(nil)

		assert.Error(t, err)
	})

	t.Run("job-with-worker-errors", func(t *testing.T) {
		var task jobs.Task = &errTask{}
		errWorker := jobs.NewWorker(&task, "test worker", nil)
		errWorkers := make([]jobs.Worker, 0)
		errWorkers = append(errWorkers, *errWorker)
		job := jobs.NewJob("test job", &errWorkers)
		consumer := &errConsumer{}
		job.ErrorConsumer = consumer
		err := jobs.DefaultScheduler().SubmitJob(job)

		job.Wait()

		assert.NoError(t, err)
		assert.Error(t, consumer.err)
	})
}

type integerTask struct{}

func (i integerTask) Run() (interface{}, error) {
	rand.Seed(420)
	dur, _ := time.ParseDuration(fmt.Sprintf("%fms", rand.Float64()))
	time.Sleep(dur)
	num := 5
	return num, nil
}

type additionConsumer struct {
	mu  sync.Mutex
	Sum int
}

func (a *additionConsumer) Consume(data interface{}) {
	integer, ok := data.(int)
	if ok {
		a.Sum += integer
	}
}

type errTask struct {}
type errConsumer struct {
	err error
}

func (e errTask) Run() (interface{}, error) {
	return nil, errors.New("test error")
}

func (e *errConsumer) Consume(err error) {
	e.err = err
}
