package crony

import (
	"os/exec"
	"syscall"

	"github.com/mattn/go-shellwords"
)

type Process struct {
	command *exec.Cmd
	done    chan int
}

func (process *Process) Done() chan int {
	return process.done
}

type Executor interface {
	Execute(command string) (*Process, error)
}

type FakeExecutor struct {
	Commands []string
}

func (executor *FakeExecutor) Execute(command string) (*Process, error) {
	executor.Commands = append(executor.Commands, command)
	return &Process{}, nil
}

func CreateFakeExecutor() *FakeExecutor {
	return &FakeExecutor{make([]string, 0)}
}

type LocalExecutor struct {
	stdout ReaderConsumer
	stderr ReaderConsumer
}

func waitForExitCode(cmd *exec.Cmd, doneChan chan int) {
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				doneChan <- status.ExitStatus()
			} else {
				panic("Unable to convert to syscall.WaitStatus; this shouldn't happen")
			}
		} else {
			// we failed while attempting to launch process
			doneChan <- -1
		}
	} else {
		doneChan <- 0
	}
}

func (executor *LocalExecutor) Execute(command string) (*Process, error) {
	commandParts, err := shellwords.Parse(command)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(commandParts[0], commandParts[1:]...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	doneChan := make(chan int)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go executor.stderr.ConsumeReader(stderr)
	go executor.stdout.ConsumeReader(stdout)
	go waitForExitCode(cmd, doneChan)

	return &Process{
		command: cmd,
		done:    doneChan,
	}, nil
}

func NewLocalExecutor(stdout ReaderConsumer, stderr ReaderConsumer) *LocalExecutor {
	return &LocalExecutor{stdout, stderr}
}
