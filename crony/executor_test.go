package crony

import (
	"testing"
)

func TestFakeExecutorRecordsCommands(t *testing.T) {
	executor := CreateFakeExecutor()

	_, err := executor.Execute("echo the command")

	if err != nil {
		t.Error("should not return an error")
	}

	if len(executor.Commands) != 1 {
		t.Error("should have recorded one command")
	}

	if executor.Commands[0] != "echo the command" {
		t.Error("should have recorded the command string")
	}
}

func TestLocalExecutorRunsTheCommandOnTheShell(t *testing.T) {
	stdoutSink := NewSink()
	stderrSink := NewSink()
	executor := NewLocalExecutor(stdoutSink, stderrSink)

	process, err := executor.Execute("/bin/echo hi")

	if err != nil {
		t.Error("should not return an error")
	}

	if exit := <-process.Done(); exit != 0 {
		t.Errorf("Should have exited with 0, got %q", exit)
	}

	if stderrSink.Content != "" {
		t.Error("should have gotten no stderr")
	}

	if stdoutSink.Content != "hi\n" {
		t.Errorf("should have gotten hi, got %q", stdoutSink.Content)
	}
}
