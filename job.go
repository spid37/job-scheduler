package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog/log"

	webhook "github.com/spid37/scheduler/webhook"
)

// JobData -
type JobData interface {
	LoadData(data []byte) error
	Send() error
}

// Job -
type Job struct {
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Version      string      `json:"version"`
	Author       string      `json:"author"`
	Schedule     Schedule    `json:"schedule"`
	AllowOverlap bool        `json:"allowOverlap"`
	Type         string      `json:"type"`
	Data         interface{} `json:"data"`
	FailCount    int
	LastError    error
	IsRunning    bool
}

// ScheduledJobs -
type ScheduledJobs struct {
	Date     time.Time
	Jobs     []*Job
	BatchRef string
}

// RunJob run a job
func (j *Job) RunJob(ctx context.Context) error {
	var err error
	log.Info().
		Str("jobName", j.Name).
		Msg("job starting")

	if j.IsRunning && !j.AllowOverlap {
		return errors.New("Job does not allow overlap")
	}
	j.IsRunning = true
	jobData := j.Data.(JobData)
	err = jobData.Send()
	if err != nil {
		j.FailCount++
		j.LastError = err
	} else {
		j.FailCount = 0
	}
	j.IsRunning = false

	return err
}

func loadJobs(jobsPath string) []*Job {
	var err error
	files, err := ioutil.ReadDir(jobsPath)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	var jobs []*Job

	for _, file := range files {
		jobPath := path.Join(jobsPath, file.Name())
		fmt.Println(jobPath)

		plan, _ := ioutil.ReadFile(jobPath)
		var job Job
		var data json.RawMessage
		job.Data = &data
		err = json.Unmarshal(plan, &job)
		if err != nil {
			panic(err)
		}

		switch job.Type {
		case "webhook":
			w := new(webhook.Data)
			if err = w.LoadData(data); err != nil {
				log.Fatal().Err(err)
			}
			job.Data = w
		default:
			log.Fatal().Msgf("unknown message type: %q", job.Type)
		}
		spew.Dump("HERE!!")
		nextRun := job.Schedule.findNextRun(time.Now())
		if nextRun.IsZero() {
			log.Fatal().Msgf("failed to find next run for: %s", job.Name)
		}

		jobs = append(jobs, &job)
	}

	return jobs
}

func getScheduledJobs(jobs []*Job) ScheduledJobs {
	date := time.Now()
	batchRef := date.Format("2006-01-02T15:04")
	scheduledJobs := ScheduledJobs{
		BatchRef: batchRef,
		Date:     date,
	}
	for _, job := range jobs {
		if job.Schedule.isNow(date) {
			log.Info().
				Str("jobName", job.Name).
				Msg("job should run now")
			scheduledJobs.Jobs = append(scheduledJobs.Jobs, job)
		} else {
			nextRun := job.Schedule.findNextRun(date)
			log.Debug().
				Str("jobName", job.Name).
				Str("runDate", nextRun.Format(time.RFC3339)).
				Msg("job next run")
		}
	}
	return scheduledJobs
}

func runJobs(
	ctx context.Context,
	workerChannel chan WorkerJob,
	scheduledJobs ScheduledJobs,
) {
	log.Info().
		Int("jobCount", len(scheduledJobs.Jobs)).
		Str("batchRef", scheduledJobs.BatchRef).
		Str("runDate", scheduledJobs.Date.Format(time.RFC3339)).
		Msg("running jobs")

	if len(scheduledJobs.Jobs) == 0 {
		// no jobs to run
		return
	}

	// send jobs to the workers
	for _, job := range scheduledJobs.Jobs {
		workerJob := WorkerJob{
			BatchRef: scheduledJobs.BatchRef,
			Job:      job,
		}
		select {
		// close context on context cancel
		case <-ctx.Done():
			log.Info().Msg("run jobs exited by context done")
			return
		// add the job to the queue
		case workerChannel <- workerJob:
			log.Info().Str("jobName", job.Name).Msg("job sent to channel")
			break
		}
	}
}
