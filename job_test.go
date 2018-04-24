package main

import (
	"fmt"
	"testing"
)

func checkExpectedError(expectedErrMsg string, err error) error {
	if err == nil {
		return fmt.Errorf("Expected %q, got no error", expectedErrMsg)
	}
	if err.Error() != expectedErrMsg {
		return fmt.Errorf("Expected %q, got %q", expectedErrMsg, err.Error())
	}
	return nil
}

func checkJobCount(expectedJobCount int, jobCount int) error {
	if jobCount != expectedJobCount {
		return fmt.Errorf("Expected %d jobs, got %d", expectedJobCount, jobCount)
	}
	return nil
}

func TestLoadJobs(t *testing.T) {
	jobs, _ := loadJobs("./jobs")
	err := checkJobCount(2, len(jobs))
	if err != nil {
		t.Error(err)
	}
}

func TestLoadJobsFalsePath(t *testing.T) {
	_, err := loadJobs("./dsfsdfsdf")
	expectedErrorMsg := "open ./dsfsdfsdf: no such file or directory"
	err = checkExpectedError(expectedErrorMsg, err)
	if err != nil {
		t.Error(err)
	}
}

func TestUnknownJobLoadJobs(t *testing.T) {
	_, err := loadJobs("./test-jobs/unknown")
	expectedErrorMsg := "unknown message type: \"fake\""
	err = checkExpectedError(expectedErrorMsg, err)
	if err != nil {
		t.Error(err)
	}
}

func TestBrokenJobJsonLoadJobs(t *testing.T) {
	_, err := loadJobs("./test-jobs/broken-json")
	expectedErrorMsg := "unexpected end of JSON input"
	err = checkExpectedError(expectedErrorMsg, err)
	if err != nil {
		t.Error(err)
	}
}
