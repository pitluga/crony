package cron

import(
	"time"
)

type Server struct {
	executor Executor
	ticker *time.Ticker
	jobs []Job
}

func Start(executor Executor, jobs []Job, interval time.Duration) *Server {
	ticker := time.NewTicker(interval)

	server := &Server{
		executor: executor,
		jobs: jobs,
		ticker: ticker,
	}

	go func () {
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
