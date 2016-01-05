package cron

import(
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type FieldMatcher interface {
	Matches(timePart int) bool
}

type Wildcard struct {
	field string
}

func (_ *Wildcard) Matches(timePart int) bool {
	return true
}

type Constant struct {
	field string
	value int
}

func (constant *Constant) Matches(timePart int) bool {
	return constant.value == timePart
}

type Range struct {
	field string
	min, max int
}

func (r *Range) Matches(timePart int) bool {
	return timePart >= r.min && timePart <= r.max
}

type CronSchedule struct {
	minutes FieldMatcher
	hours FieldMatcher
	dayOfMonth FieldMatcher
	monthOfYear FieldMatcher
	dayOfWeek FieldMatcher
	year FieldMatcher
}

type Schedule interface {
	ShouldRun(t time.Time) bool
}

func (schedule CronSchedule) ShouldRun(t time.Time) bool {
	return schedule.minutes.Matches(t.Minute()) && schedule.hours.Matches(t.Hour()) && schedule.dayOfMonth.Matches(t.Day()) && schedule.monthOfYear.Matches(int(t.Month())) && schedule.dayOfWeek.Matches(int(t.Weekday())) && schedule.year.Matches(t.Year())

}

func parseField(field string) FieldMatcher {
	constantRegexp, err := regexp.Compile(`\A\d+\z`)
	if err != nil {
		panic("invalid hardcoded regexp!")
	}
	rangeRegexp, err := regexp.Compile(`\A\d+\-\d+\z`)
	if err != nil {
		panic("invalid hardcoded regexp!")
	}

	if "*" == field {
		return &Wildcard{field}
	} else if constantRegexp.MatchString(field) {
		number, err := strconv.Atoi(field)
		if err != nil {
			panic(fmt.Sprintf("regexp matched integer, but couldn't convert %q to int", field))
		}
		return &Constant{field, number}
	} else if rangeRegexp.MatchString(field) {
		rangeParts := strings.Split(field, "-")
		min, err := strconv.Atoi(rangeParts[0])
		if err != nil {
			panic(fmt.Sprintf("regexp matched integer, but couldn't convert %q to int", rangeParts[0]))
		}
		max, err := strconv.Atoi(rangeParts[1])
		if err != nil {
			panic(fmt.Sprintf("regexp matched integer, but couldn't convert %q to int", rangeParts[1]))
		}
		return &Range{field, min, max}
	}
	return nil
}

func Parse(pattern string) Schedule {
	fields := strings.Split(pattern, " ")

	return &CronSchedule{
		minutes: parseField(fields[0]),
		hours: parseField(fields[1]),
		dayOfMonth: parseField(fields[2]),
		monthOfYear: parseField(fields[3]),
		dayOfWeek: parseField(fields[4]),
		year: parseField(fields[4]),
	}
}
