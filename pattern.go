package cron

import(
	"fmt"
	"math"
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

type Interval struct {
	field string
	numerator FieldMatcher
	denominator int
}

func (interval *Interval) Matches(timePart int) bool {
	return interval.numerator.Matches(timePart) && math.Mod(float64(timePart), float64(interval.denominator)) == 0
}

type List struct {
	field string
	values map[int]bool
}

func (list *List) Matches(timePart int) bool {
	return list.values[timePart]
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

func safeAtoi(str string) int {
	number, err := strconv.Atoi(str)
	if err != nil {
		panic(fmt.Sprintf("regexp matched integer, but couldn't convert %q to int", str))
	}
	return number
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
	intervalRegexp, err := regexp.Compile(`\A.*/\d+\z`)
	if err != nil {
		panic("invalid hardcoded regexp!")
	}
	listRegexp, err := regexp.Compile(`\A\d+(,\d+)+\z`)
	if err != nil {
		panic("invalid hardcoded regexp!")
	}

	if "*" == field {
		return &Wildcard{field}
	} else if constantRegexp.MatchString(field) {
		return &Constant{field, safeAtoi(field)}
	} else if rangeRegexp.MatchString(field) {
		rangeParts := strings.Split(field, "-")
		return &Range{
			field: field,
			min: safeAtoi(rangeParts[0]),
			max: safeAtoi(rangeParts[1]),
		}
	} else if intervalRegexp.MatchString(field) {
		intervalParts := strings.Split(field, "/")
		return &Interval{
			field: field,
			numerator: parseField(intervalParts[0]),
			denominator: safeAtoi(intervalParts[1]),
		}
	} else if listRegexp.MatchString(field) {
		parts := strings.Split(field, ",")
		values := make(map[int]bool)
		for _, part := range(parts) {
			values[safeAtoi(part)] = true
		}
		return &List{
			field: field,
			values: values,
		}
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
