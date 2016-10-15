package main

import (
	"fmt"
	"time"
)

var timeZone *time.Location

func init() {
	timeZone , _ = time.LoadLocation("America/New_York")
}

type Date struct {
	Day int
	Month time.Month
	Year int
}

func NewDate(t time.Time) *Date {
	d := new(Date)
	d.Year, d.Month, d.Day = t.Date()
	return d
}

func (d* Date) Set(s string) {
	t, err := time.ParseInLocation("Monday, January 2, 2006", s, timeZone)
	if err != nil {
		t, err = time.ParseInLocation("Mon 01/02/2006", s, timeZone)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	d.Year, d.Month, d.Day = t.Date()
}

func (d* Date) Today() {
	d.Year, d.Month, d.Day = time.Now().Date()
}

func (d* Date) Time() time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, timeZone)
}

func (d* Date) Format(layout string) string {
	return d.Time().Format(layout)
}

func (d* Date) AddDayTime(timeStr string) time.Time {
	t, err := time.Parse("3:04 PM", timeStr)
	if err != nil {
		t, err = time.Parse("3:04 PM MST", timeStr)
		if err != nil {
			fmt.Println(err.Error())
			return t
		}
	}

	midnight, _ := time.Parse("3:04 PM", "12:00 AM")
	dur := t.Sub(midnight)

	return d.Time().Add(dur)
}

func main() {
	var date Date

	/* Formatting the Date */
	date.Today()
	fmt.Println(date)
	fmt.Println(date.Format("Mon Jan 2 2006"))
	fmt.Println(date.Format("Monday Jan 2"))

	/* Create a Date from Time */
	d := NewDate(time.Now())
	fmt.Println(d.Format("2 Jan 2006"))

	/* Get Time from Date */
	t := d.Time()
	fmt.Println(t)

	/* Add time to date */
	fmt.Println(date.AddDayTime("9:04 AM"))

	fmt.Println(date.AddDayTime("7:30 PM EST"))

	date.Set("Sun 01/03/2016")
	fmt.Println(date)

	date.Set("Thursday, August 18, 2016")
	fmt.Println(date)
}
