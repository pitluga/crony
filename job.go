package cron

type Job struct {
	Schedule Schedule
	Command string
}
