package cron

import(
	"testing"
	"time"
)


func TestOnceWithNoJobsDoesNothing(t *testing.T) {
	executor := CreateFakeExecutor()
	server := &Server{
		executor: executor,
		jobs: make([]Job, 0),
	}

	server.Once(time.Now())

	if len(executor.Commands) > 0 {
		t.Error("Should not have run commands")
	}
}

func TestOnceWithAWildcardJobRunsIt(t *testing.T) {
	executor := CreateFakeExecutor()
	server := &Server{
		executor: executor,
		jobs: []Job{
			Job{Parse("* * * * * *"), "echo hi"},
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
	executor := CreateFakeExecutor()
	server := Start(
		executor,
		[]Job{Job{Parse("* * * * * *"), "echo hi"}},
		time.Millisecond * 10,
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
