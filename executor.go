package cron

type Executor interface {
	Execute(command string) error
}

type FakeExecutor struct {
	Commands []string
}

func (executor *FakeExecutor) Execute(command string) error {
	executor.Commands = append(executor.Commands, command)
	return nil
}

func CreateFakeExecutor() *FakeExecutor {
	return &FakeExecutor{make([]string, 0)}
}
