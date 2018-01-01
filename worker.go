package main

import (
	"context"
	"sync"

	logger "github.com/rs/zerolog/log"
)

// WorkerJob -
type WorkerJob struct {
	BatchRef string
	Job      *Job
}

// Workers -
type Workers struct {
	wg *sync.WaitGroup
	ch chan WorkerJob
}

func startWorkers(ctx context.Context, workers int) Workers {
	var workerChannel = make(chan WorkerJob)
	var wg sync.WaitGroup

	for worker := 0; worker < workers; worker++ {
		go func(workerID int) {
			wg.Add(1)
			defer wg.Done()
			startWorker(ctx, workerChannel, workerID)
		}(worker)
	}

	return Workers{&wg, workerChannel}
}

func startWorker(ctx context.Context, workerChannel <-chan WorkerJob, id int) {
	logger.Info().Int("workerID", id).Msg("Starting worker")
	for {
		select {
		case <-ctx.Done():
			logger.Info().Int("workerID", id).Msg("Closing worker by context done")
			return
		case workerJob := <-workerChannel:
			logger.Info().
				Str("batchRef", workerJob.BatchRef).
				Int("workerID", id).
				Str("jobName", workerJob.Job.Name).
				Msg("worker received a job")
			// run the job
			err := workerJob.Job.RunJob(ctx)
			if err != nil {
				logger.Warn().
					Str("batchRef", workerJob.BatchRef).
					Int("workerID", id).
					Str("jobName", workerJob.Job.Name).
					Str("error", err.Error()).
					Msg("job failed")
			} else {
				logger.Info().
					Str("batchRef", workerJob.BatchRef).
					Int("workerID", id).
					Str("jobName", workerJob.Job.Name).
					Msg("job successful")
			}

		}
	}
}
