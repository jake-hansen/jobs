// Copyright © 2021 Jacob Hansen. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package schedulers_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/jake-hansen/jobs/consumers"
	"github.com/jake-hansen/jobs/jobs"
	"github.com/jake-hansen/jobs/schedulers"
	"github.com/stretchr/testify/assert"
)

func TestDefaultScheduler(t *testing.T) {
	testScheduler := schedulers.DefaultScheduler()

	assert.Equal(t, consumers.DataPrinterConsumer{}, testScheduler.DataConsumer)
	assert.Equal(t, consumers.ErrorPrinterConsumer{}, testScheduler.ErrorConsumer)
	assert.Equal(t, schedulers.SequentialScheduler{}, testScheduler.Algorithm)
}

func spawnWorkers() *[]jobs.Worker {
	var workers []jobs.Worker
	for i := 0; i < 100; i++ {
		workers = append(workers, integerWorker{Name: strconv.Itoa(i)})
	}
	return &workers
}

func TestScheduler_Schedule(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testWorkers := spawnWorkers()

		testJob := jobs.NewJob("test job", testWorkers)
		testScheduler := schedulers.DefaultScheduler()
		consumer := additionConsumer{Sum: 0}
		testScheduler.DataConsumer = &consumer
		err := testScheduler.Schedule(testJob)
		testScheduler.WaitForWorkers()

		assert.NoError(t, err)
		assert.Equal(t, len(*testWorkers)*5, consumer.safeSumRead())
	})

	t.Run("failure-multiple-job", func(t *testing.T) {
		testWorkers := spawnWorkers()

		testJob := jobs.NewJob("test job", testWorkers)
		testScheduler := schedulers.DefaultScheduler()
		consumer := additionConsumer{Sum: 0}
		testScheduler.DataConsumer = &consumer
		err := testScheduler.Schedule(testJob)
		assert.NoError(t, err)
		err = testScheduler.Schedule(testJob)
		assert.Error(t, err)
	})

	t.Run("failure-nil-job", func(t *testing.T) {
		testScheduler := schedulers.DefaultScheduler()
		consumer := additionConsumer{Sum: 0}
		testScheduler.DataConsumer = &consumer
		err := testScheduler.Schedule(nil)
		testScheduler.WaitForWorkers()

		assert.Error(t, err)
	})
}

type integerWorker struct {
	Name string
}

func (i integerWorker) Run() (interface{}, error) {
	rand.Seed(420)
	dur, _ := time.ParseDuration(fmt.Sprintf("%fms", rand.Float64()))
	time.Sleep(dur)
	num := 5
	return num, nil
}

func (i integerWorker) WorkerName() string {
	return i.Name
}

func (i integerWorker) GetPriority() interface{} {
	return 420
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
