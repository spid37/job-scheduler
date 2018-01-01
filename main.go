package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
)

const configJobsPath = "./jobs"

func init() {
	logger.Logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	// load the jobs from json files
	jobs := loadJobs(configJobsPath)

	// start the ticker and run the jobs
	start(jobs)
}

func start(jobs []*Job) {
	// create context to stop workers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start the workers
	workerCount := 10
	workers := startWorkers(ctx, workerCount)
	defer close(workers.ch)

	logger.Info().
		Int("count", workerCount).
		Msg("Workers started")

	// create a callback to send to the ticker to run on interval
	tick := func(ctx context.Context) {
		scheduledJobs := getScheduledJobs(jobs)
		logger.Info().
			Str("BatchRef", scheduledJobs.BatchRef).
			Str("RunDate", scheduledJobs.Date.Format(time.RFC3339)).
			Msg("Tick started")
		runJobs(ctx, workers.ch, scheduledJobs)
	}
	// start the ticker
	startTicker(ctx, tick)

	// catch ctrl-c and exit
	quiter := func() {
		logger.Info().Msg("Quitter has been called")
		cancel()
	}
	catchInterrupt(quiter)

	workers.wg.Wait() // wait for workers to quit
	logger.Info().Msg("Exiting..")
}

// catchInterrupt listen for an interrupt
// preform callback on action
func catchInterrupt(cb func()) chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		defer close(c)
		for sig := range c {
			logger.Info().Str("sig", sig.String()).Msg("Received interrupt")
			//fmt.Printf("Received ctrl-c: %s\n", sig)
			// sig is a ^C, handle it
			cb()
			return
		}
	}()

	return c
}
