package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidDate  = errors.New("Invalid date")
	ErrInvalidYear  = errors.New("Invalid year")
	ErrInvalidMonth = errors.New("Invalid month")
	ErrInvalidDay   = errors.New("Invalid day")
)

type BusinessHours struct {
	Start int32
	End   int32
}

// Date stores year, month, & day of a single date.
type Date struct {
	Year  int
	Month int
	Day   int
}

func (d *Date) ToString() string {
	year := strconv.Itoa(d.Year)
	month := strconv.Itoa(d.Month)
	if d.Month < 10 {
		month = fmt.Sprintf("0%s", month)
	}
	day := strconv.Itoa(d.Day)
	if d.Day < 10 {
		day = fmt.Sprintf("0%s", day)
	}

	return fmt.Sprintf("%s-%s-%s", year, month, day)
}

func (d *Date) ToBusinessHours(startHour, endHour, timezone string) (BusinessHours, error) {
	var bh BusinessHours

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return bh, err
	}

	startStr := fmt.Sprintf("%sT%s", d.ToString(), startHour)
	endStr := fmt.Sprintf("%sT%s", d.ToString(), endHour)
	layout := "2006-01-02T15:04:05"

	startTime, err := time.ParseInLocation(layout, startStr, location)
	if err != nil {
		return bh, err
	}
	endTime, err := time.ParseInLocation(layout, endStr, location)
	if err != nil {
		return bh, err
	}

	bh = BusinessHours{Start: int32(startTime.Unix()), End: int32(endTime.Unix())}
	return bh, nil
}

func DateStringToTime(date, timezone string) (time.Time, error) {
	var t time.Time

	dateSlice := strings.Split(date, "-")
	if len(dateSlice) != 3 {
		return t, ErrInvalidDate
	}

	year, err := strconv.Atoi(dateSlice[0])
	if err != nil {
		return t, ErrInvalidYear
	}
	month, err := strconv.Atoi(dateSlice[1])
	if err != nil {
		return t, ErrInvalidMonth
	}
	day, err := strconv.Atoi(dateSlice[2])
	if err != nil {
		return t, ErrInvalidDay
	}
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return t, err
	}

	t = time.Date(year, time.Month(month), day, 0, 0, 0, 0, location)
	return t, nil
}

// Timespan stores dates, hours in the day, & timezone in RFC3339 format.
type Timespan struct {
	StartDate    string
	EndDate      string
	DayStartHour string
	DayEndHour   string
	Timezone     string
}

func (t *Timespan) startDateToTime() (time.Time, error) {
	return DateStringToTime(t.StartDate, t.Timezone)
}

func (t *Timespan) endDateToTime() (time.Time, error) {
	return DateStringToTime(t.EndDate, t.Timezone)
}

func (t *Timespan) DatesInBetween() ([]Date, error) {
	var dates []Date

	d1, err := t.startDateToTime()
	if err != nil {
		return dates, err
	}

	d2, err := t.endDateToTime()
	if err != nil {
		return dates, err
	}

	today := d1
	for today.Before(d2) || today.Equal(d2) {
		year, month, day := today.Date()
		dates = append(dates, Date{
			Year:  year,
			Month: int(month),
			Day:   day,
		})
		tomorrow := today.AddDate(0, 0, +1)
		today = tomorrow
	}

	return dates, nil
}

func (t *Timespan) BusinessHoursOfEachDate() ([]BusinessHours, error) {
	var bhs []BusinessHours

	dates, err := t.DatesInBetween()
	if err != nil {
		return bhs, err
	}

	for _, d := range dates {
		bh, err := d.ToBusinessHours(t.DayStartHour, t.DayEndHour, t.Timezone)
		if err != nil {
			return bhs, err
		}
		bhs = append(bhs, bh)
	}

	return bhs, nil
}

func main() {
	t1 := Timespan{
		StartDate:    "2020-07-10",
		EndDate:      "2020-07-14",
		DayStartHour: "08:00:00",
		DayEndHour:   "22:00:00",
		Timezone:     "Asia/Jakarta",
	}
	t1bhs, err := t1.BusinessHoursOfEachDate()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", t1bhs)
}
