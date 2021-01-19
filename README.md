# Jobs
Jobs is a simple asynchronous job scheduling library for Go. Jobs works by abstracting goroutines to a simple yet highly configurable job scheduling API. 

## Contents
- [Jobs](#jobs)
  - [Contents](#contents)
  - [Intro](#intro)
  - [Installation](#installation)
  - [Example](#example)

## Intro

There are three main components to the Jobs library: Workers, Jobs, and Schedulers.

Worker - A Worker is a single, atomic unit of work. This means that a Worker simply performs a task or operation.

Jobs - A Job is a collection of Workers. Multiple Workers work together within a Job to perform a larger operation or task, which is considered the Job.

Schedulers - A Scheduler manages the configuration for a Job and is reponsible for schedueling the inidvidual Workers within a Job.

## Installation

To use Jobs, you need to install Go and set your Go workspace.

1. To install Jobs, you'll need to run this command within your project.
   
   ```sh
   $ go get -u github.com/jake-hansen/jobs
   ```

2. Import jobs in your project

    ```go
    import "github.com/jake-hansen/jobs
    ```

## Example

The Jobs API is best described by working through an example on how break up a larger unit of work into smaller atomic units of work. Once you have a grasp on how the API works, it should become clear how Jobs can fit within your application. Keep in mind that Jobs performs best when managing tasks which are asynchronous by nature, meaning that each Worker does not rely upon the completion of another Worker before it can begin.

The example we'll work through is creating a simple Monte Carlo Pi approximator which has been adapted from https://golang.org/doc/play/pi.go.

It's not necessary to understand the mechanics and proof behind the Monte Carlo Pi approxmation to see how it is easily implemented using an asynchronous mechanism such as Jobs.

TL;DR - If you perform a certain calculation enough times and sum the result of each calculation, you will obtain a close approximation of Pi. The more times you perform the calculation, the more accurate the result becomes.

Let's get started.

We'll first start by defining our **Worker**.

Our PiWorker is an implementation of a Worker. You'll see some helper functions below such as `WorkerName()` and `GetPriority()`. These functions are useful in other contexts of the Jobs library and are not necessary to understand for this example.

Where the magic happens is in the `Run()` function. You'll see that the function takes no parameters, but performs the Monte Carlo calculation using a k value which is provided in the struct.

```go
type PiWorker struct {
	Name        string
	Priority    int
	KVal	    float64
}

func (p PiWorker) Run() (interface{}, error) {
	return 4 * math.Pow(-1, p.KVal) / (2 * p.KVal + 1), nil
}

func (p PiWorker) WorkerName() string {
	return p.Name
}

func (p PiWorker) GetPriority() interface{} {
	return p.Priority
}
```

PiWorker, by itself, is not very useful. Sure, we can create a new PiWorker and execute the function `Run()`. However, we wouldn't get a very accurate approximation of Pi. In order to get a more accurate approximation of Pi, as explained above, we need to perform the calculation that `Run()` provides multiple times using different k values, and then sum the results together. Jobs to the rescue!

The next thing we need to do is define a Job. Remember that a Job is a *collection* of workers. So we will first need a collection of our PiWorkers. To do this, we'll create a helper function. This helper function will take a paramter, `n`, and return a slice of that many PiWorkers. Note that each PiWorker takes on a unique `KVal`. 

```go
func createPiWorkers(n int) *[]jobs.Worker {
	var piSlice []jobs.Worker
	for i := 0; i <= n; i++ {
		piSlice = append(piSlice, piWorker{
			Name:     "piworker",
			Priority: 0,
			KVal:     float64(i),
		})
	}
	return &piSlice
}
```

Now that we have the ability to generate an arbitrary amount of PiWorkers, we can create our Job.

Here, we've created a variable and stored a new Job. The magic here is that we've stored *1000* different PiWorkers inside our Job.

```go
calculatePi := jobs.NewJob("monte carlo pi approximation", createPiWorkers(1000))
```

Now that we have our Job definition created, we need to schedule the workers within this Job to actually execute.

"But, wait!" you say, "I thought we needed to add up the individual calculations together." You are correct. Here is how we do that.

Jobs has the concept of a `DataConsumer`. A DataConsumer takes the result of a Woker and does *something* with that result. Let's go ahead and define a DataConsumer that will work for our use case.

PiAddition, implements the DataConsumer interface by defining a function `Consume(data interface{})`.

Our Consume function here takes in a paramter, `data`, checks to make sure the data is of type float64, and if it is, adds that value to our Pi sum. Note that the type checking here is important. 

```go
type PiAddition struct {
	Pi	float64
}

func (p *PiAddition) Consume(data interface{}) {
	if v, ok := data.(float64); ok {
		p.Pi += v
	}
}
```

Once we've created our Consumer, we can finally schedule our Job! Let's do that below.

Here, we create a DefaultScheduler. We then set the consumer of that scheduler to be the consumer we just defined above. Finally, we schedule our job with the line `piScheduler.Schedule(calculatePi)`.

In order to make sure our program doesn't exit until the last worker thread finishes, we need to make the call `piScheduler.WaitForWorkers()`, which blocks until the last Worker returns.

Finally, we print the Pi approximation stored in our consumer.

```go
piScheduler := schedulers.DefaultScheduler()
piConsumer := piAddition{Pi: 0}
piScheduler.DataConsumer = &piConsumer
_ = piScheduler.Schedule(calculatePi)

piScheduler.WaitForWorkers()

fmt.Println(piConsumer.Pi)
```

Now we have the result `3.1425916543395447`. If we run the scheduler again, this time with *1000000* Workers, we get a value of `3.1415936535887727`. As you can see as the number of Worker threads increase, so does the accuracy of our Pi approximation.

This example should have provided you a better understanding of the Jobs library. Clearly, the Monte Carlo Pi approximation as shown [here](https://golang.org/doc/play/pi.go) appears to be much simpler than the example that was just demonstrated. Keep in mind that this demo was just an example. The Jobs library has rich functionality built in such as defining priority for tasks and custom scheduling algorithms. The Monte Carlo Pi approximation algorithm did not take advantage of these features.
