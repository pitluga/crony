package cron

import(
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type UnitMatcher interface {
	Matches(timePart int) bool
}

type Wildcard struct {
	subpattern string
}

func (_ *Wildcard) Matches(timePart int) bool {
	return true
}

type Constant struct {
	subpattern string
	value int
}

func (constant *Constant) Matches(timePart int) bool {
	return constant.value == timePart
}

type Range struct {
	subpattern string
	min, max int
}

func (r *Range) Matches(timePart int) bool {
	return timePart >= r.min && timePart <= r.max
}

type CronSchedule struct {
	minutes UnitMatcher
	hours UnitMatcher
	dayOfMonth UnitMatcher
	monthOfYear UnitMatcher
	dayOfWeek UnitMatcher
	year UnitMatcher
}

type Schedule interface {
	ShouldRun(t time.Time) bool
}

func (schedule CronSchedule) ShouldRun(t time.Time) bool {
	return schedule.minutes.Matches(t.Minute()) && schedule.hours.Matches(t.Hour()) && schedule.dayOfMonth.Matches(t.Day()) && schedule.monthOfYear.Matches(int(t.Month())) && schedule.dayOfWeek.Matches(int(t.Weekday())) && schedule.year.Matches(t.Year())

}

func parsePart(subpattern string) UnitMatcher {
	constantRegexp, err := regexp.Compile(`\A\d+\z`)
	if err != nil {
		panic("invalid hardcoded regexp!")
	}
	rangeRegexp, err := regexp.Compile(`\A\d+\-\d+\z`)
	if err != nil {
		panic("invalid hardcoded regexp!")
	}

	if "*" == subpattern {
		return &Wildcard{subpattern}
	} else if constantRegexp.MatchString(subpattern) {
		number, err := strconv.Atoi(subpattern)
		if err != nil {
			panic(fmt.Sprintf("regexp matched integer, but couldn't convert %q to int", subpattern))
		}
		return &Constant{subpattern, number}
	} else if rangeRegexp.MatchString(subpattern) {
		rangeParts := strings.Split(subpattern, "-")
		min, err := strconv.Atoi(rangeParts[0])
		if err != nil {
			panic(fmt.Sprintf("regexp matched integer, but couldn't convert %q to int", rangeParts[0]))
		}
		max, err := strconv.Atoi(rangeParts[1])
		if err != nil {
			panic(fmt.Sprintf("regexp matched integer, but couldn't convert %q to int", rangeParts[1]))
		}
		return &Range{subpattern, min, max}
	}
	return nil
}

func Parse(pattern string) Schedule {
	parts := strings.Split(pattern, " ")

	return &CronSchedule{
		minutes: parsePart(parts[0]),
		hours: parsePart(parts[1]),
		dayOfMonth: parsePart(parts[2]),
		monthOfYear: parsePart(parts[3]),
		dayOfWeek: parsePart(parts[4]),
		year: parsePart(parts[4]),
	}
}
