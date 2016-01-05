package cron

import(
	"testing"
	"time"
)

type TimeTestCase struct {
	when time.Time
	runs bool
}

func IsScheduled(t *testing.T, pattern string, tests []TimeTestCase) {
	schedule := Parse(pattern)

	for _, test := range(tests) {
		if test.runs != schedule.ShouldRun(test.when) {
			should := ""
			if !test.runs {
				should = "not "
			}

			t.Errorf("pattern %q should %vrun at time %v", pattern, should, test.when)
		}
	}
}

func TestAllStars(t *testing.T) {
	IsScheduled(t, "* * * * * *", []TimeTestCase{{time.Now(), true}})
}

func TestFiveMinutesAfterTheHour(t *testing.T) {
	IsScheduled(t, "5 * * * * *", []TimeTestCase{
		{time.Date(2016, time.January, 5, 10, 15, 0, 0, time.UTC), false},
		{time.Date(2016, time.January, 5, 10, 5, 0, 0, time.UTC), true},
	})
}

func TestRanges(t *testing.T) {
	IsScheduled(t, "4-8 * * * * *", []TimeTestCase{
		{time.Date(2016, time.January, 5, 10, 15, 0, 0, time.UTC), false},
		{time.Date(2016, time.January, 5, 10, 5, 0, 0, time.UTC), true},
		{time.Date(2016, time.January, 5, 10, 4, 0, 0, time.UTC), true},
		{time.Date(2016, time.January, 5, 10, 8, 0, 0, time.UTC), true},
	})
}

func TestIntervals(t *testing.T) {
	IsScheduled(t, "*/3 * * * * *", []TimeTestCase{
		{time.Date(2016, time.January, 5, 10, 15, 0, 0, time.UTC), true},
		{time.Date(2016, time.January, 5, 10, 5, 0, 0, time.UTC), false},
		{time.Date(2016, time.January, 5, 10, 4, 0, 0, time.UTC), false},
		{time.Date(2016, time.January, 5, 10, 33, 0, 0, time.UTC), true},
	})

	IsScheduled(t, "1-30/3 * * * * *", []TimeTestCase{
		{time.Date(2016, time.January, 5, 10, 15, 0, 0, time.UTC), true},
		{time.Date(2016, time.January, 5, 10, 5, 0, 0, time.UTC), false},
		{time.Date(2016, time.January, 5, 10, 4, 0, 0, time.UTC), false},
		{time.Date(2016, time.January, 5, 10, 33, 0, 0, time.UTC), false},
	})
}
