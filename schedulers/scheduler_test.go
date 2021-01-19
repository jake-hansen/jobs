package schedulers_test

import (
	"fmt"
	"math/rand"
	"strconv"
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
		assert.Equal(t, len(*testWorkers)*5, consumer.Sum)
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
	Sum int
}

func (a *additionConsumer) sum(value int) {
	a.Sum += value
}

func (a *additionConsumer) Consume(data interface{}) {
	integer, ok := data.(int)
	if ok {
		a.sum(integer)
	}
}
