package server

import (
	"github.com/pitluga/crony/crony"
	"time"
)

type Server struct {
	executor crony.Executor
	ticker   *time.Ticker
	jobs     []crony.Job
}

func Start(executor crony.Executor, jobs []crony.Job, interval time.Duration) *Server {
	ticker := time.NewTicker(interval)

	server := &Server{
		executor: executor,
		jobs:     jobs,
		ticker:   ticker,
	}

	go func() {
		for now := range ticker.C {
			server.Once(now)
		}
	}()

	return server
}

func (server *Server) Once(now time.Time) {
	for _, job := range server.jobs {
		if job.Schedule.ShouldRun(now) {
			server.executor.Execute(job.Command)
		}
	}
}

func (server *Server) Stop() {
	server.ticker.Stop()
}
