// Copyright © 2021 Jacob Hansen. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package jobs_test

import (
	"fmt"
	"github.com/jake-hansen/jobs"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

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
		consumer := additionConsumer{Sum: 0}
		testJob.DataConsumer = &consumer
		err := testScheduler.SubmitJob(testJob)
		testScheduler.WaitForWorkers()

		assert.NoError(t, err)
		assert.Equal(t, len(*testWorkers)*5, consumer.safeSumRead())
	})

	t.Run("failure-multiple-job", func(t *testing.T) {
		testWorkers := spawnWorkers()

		testJob := jobs.NewJob("test job", testWorkers)
		testScheduler := jobs.DefaultScheduler()
		consumer := additionConsumer{Sum: 0}
		testJob.DataConsumer = &consumer
		err := testScheduler.SubmitJob(testJob)
		assert.NoError(t, err)
		err = testScheduler.SubmitJob(testJob)
		assert.Error(t, err)
	})

	t.Run("failure-nil-job", func(t *testing.T) {
		testScheduler := jobs.DefaultScheduler()
		err := testScheduler.SubmitJob(nil)
		testScheduler.WaitForWorkers()

		assert.Error(t, err)
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

// safeSumRead is a thread safe implementation to read the Sum
// variable.
func (a *additionConsumer) safeSumRead() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.Sum
}

// safeSumWrite() is a thread safe implementation to write the Sum
// variable.
func (a *additionConsumer) safeSumWrite(value int) {
	a.mu.Lock()
	a.Sum += value
	defer a.mu.Unlock()
}

func (a *additionConsumer) Consume(data interface{}) {
	integer, ok := data.(int)
	if ok {
		a.safeSumWrite(integer)
	}
}
