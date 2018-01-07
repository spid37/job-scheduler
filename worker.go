package main

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
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
	log.Info().Int("workerID", id).Msg("starting worker")
	for {
		select {
		case <-ctx.Done():
			log.Info().Int("workerID", id).Msg("closing worker by context done")
			return
		case workerJob := <-workerChannel:
			log.Info().
				Str("batchRef", workerJob.BatchRef).
				Int("workerID", id).
				Str("jobName", workerJob.Job.Name).
				Msg("worker received a job")
			// run the job
			err := workerJob.Job.RunJob(ctx)
			if err != nil {
				log.Warn().
					Err(err).
					Str("batchRef", workerJob.BatchRef).
					Int("workerID", id).
					Str("jobName", workerJob.Job.Name).
					Msg("job failed")
			} else {
				log.Info().
					Str("batchRef", workerJob.BatchRef).
					Int("workerID", id).
					Str("jobName", workerJob.Job.Name).
					Msg("job successful")
			}

		}
	}
}
