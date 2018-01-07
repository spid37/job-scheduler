package sleep

import (
	"encoding/json"
	"math/rand"
	"time"
)

// Data -
type Data struct {
	SleepSeconds int `json:"sleepSeconds"`
}

// LoadData load the job data fro webhook
func (d *Data) LoadData(data []byte) error {
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}
	return nil
}

// Send Send a Job request
// if sleep time not specified sleep for a random duration unser 30 sec
func (d *Data) Send() error {
	var err error
	var sleepDuration time.Duration
	if d.SleepSeconds == 0 {
		rand.Seed(time.Now().Unix())
		sleepDuration = time.Duration(rand.Intn(30)) * time.Second
	} else {
		sleepDuration = time.Duration(d.SleepSeconds) * time.Second
	}
	time.Sleep(sleepDuration)

	return err
}
