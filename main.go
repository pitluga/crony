package main

import (
	"fmt"
	"github.com/pitluga/crony/crony"
	"github.com/pitluga/crony/server"
	"time"
)

func main() {
	done := make(chan string)
	server.Start(
		crony.NewLocalExecutor(
			crony.NewStdoutPipe("stdout"),
			crony.NewStdoutPipe("stderr"),
		),
		[]crony.Job{
			crony.Job{crony.Parse("* * * * * *"), "echo hi"},
		},
		time.Minute,
	)

	fmt.Print("Server Started...\n")

	fmt.Print(<-done)
}
