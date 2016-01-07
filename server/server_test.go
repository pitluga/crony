package server

import (
	"github.com/pitluga/crony/crony"
	"testing"
	"time"
)

func TestOnceWithNoJobsDoesNothing(t *testing.T) {
	executor := crony.CreateFakeExecutor()
	server := &Server{
		executor: executor,
		jobs:     make([]crony.Job, 0),
	}

	server.Once(time.Now())

	if len(executor.Commands) > 0 {
		t.Error("Should not have run commands")
	}
}

func TestOnceWithAWildcardJobRunsIt(t *testing.T) {
	executor := crony.CreateFakeExecutor()
	server := &Server{
		executor: executor,
		jobs: []crony.Job{
			crony.Job{crony.Parse("* * * * * *"), "echo hi"},
		},
	}

	server.Once(time.Now())

	if len(executor.Commands) != 1 {
		t.Error("Should have run 1 command")
	}

	if executor.Commands[0] != "echo hi" {
		t.Error("Should have run 'echo hi'")
	}
}

func TestStartStopWithSingleJob(t *testing.T) {
	executor := crony.CreateFakeExecutor()
	server := Start(
		executor,
		[]crony.Job{crony.Job{crony.Parse("* * * * * *"), "echo hi"}},
		time.Millisecond*10,
	)

	timer := time.NewTimer(time.Millisecond * 12)
	<-timer.C

	server.Stop()

	if len(executor.Commands) != 1 {
		t.Error("Should have run 1 command")
	}

	if executor.Commands[0] != "echo hi" {
		t.Error("Should have run 'echo hi'")
	}

	timer = time.NewTimer(time.Millisecond * 12)
	<-timer.C
	if len(executor.Commands) != 1 {
		t.Error("Should have run 1 command")
	}
}
