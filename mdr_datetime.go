package mdr

import (
	"fmt"
	"time"
)

// turns duration into decimal minutes,hours,days as appropriate,
// maximum duration is about 290 years
// Test_002
func HumanTime(t time.Duration) (rs string) {
	sec := int64(t.Seconds()) // converts duration in nanosec to seconds
	if sec < 60 {
		rs = fmt.Sprintf("%d seconds", sec)
		return
	}
	if sec < 3600 {
		rs = fmt.Sprintf("%5.2f minutes", float64(sec)/60.0)
		return
	}
	if sec < 86400 {
		rs = fmt.Sprintf("%5.2f hours", float64(sec)/3600.0)
		return
	}
	if sec < int64(86400.0*30.4375) {
		rs = fmt.Sprintf("%5.2f days", float64(sec)/86400.0)
		return
	}
	if sec < int64(86400.0*365.25) {
		rs = fmt.Sprintf("%5.2f months", float64(sec)/(86400.0*30.4375))
		return
	}
	rs = fmt.Sprintf("%5.2f years", float64(sec)/(86400.0*365.25))
	return
}

// ValidDate tries to determine validity of given date.  It's not very
// much use in pre-Gregorian times and doesn't do BCE at all.
// https://en.wikipedia.org/wiki/Julian_calendar may help understanding,
// or may just confuse one further :-)
func ValidDate(year, month, day, hour, minute, second int) bool {
	// this isn't very kind to non-christians...
	if year < 0 {
		return false
	}
	// would upper limit on valid year make any sense?
	var monthdays []int = []int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if !InRangeI(1, month, 12) {
		return false
	}
	var days int
	if month == 2 {
		if year%400 == 0 {
			days = 29
		} else if year%100 == 0 {
			days = 28
		} else if year%4 == 0 {
			days = 29
		} else {
			days = 28
		}
	} else {
		days = monthdays[month]
	}
	if !InRangeI(1, day, days) {
		return false
	}
	if !InRangeI(0, hour, 23) {
		return false
	}
	if !InRangeI(0, minute, 59) {
		return false
	}
	if !InRangeI(0, second, 59) {
		return false
	}

	return true
}

// true if when is a leap year
func LeapYear(when time.Time) bool {
	year := when.Year()
	if year%400 == 0 {
		return true
	}
	if year%100 == 0 {
		return false
	}
	if year%4 == 0 {
		return true
	}
	return false
}

// returns decimal year to nearest 52 minutes - usually printed with %9.4f
func StarDate(when time.Time) float64 {
	yr := float64(when.Year())
	dayofyear := float64(when.YearDay())
	//fmt.Printf("%v %v\n", yr, dayofyear)
	var daysinyear float64
	if LeapYear(when) {
		daysinyear = 366
	} else {
		daysinyear = 365
	}
	// note Hour is rounded so should add fraction for Minutes*60+Seconds
	elapsedSeconds := when.Hour() * 3600
	elapsedSeconds += +when.Minute() * 60
	elapsedSeconds += when.Second()
	hrs := (dayofyear-1)*24 + (float64(elapsedSeconds) / 3600.0)
	//fmt.Printf("hrs = %v \n",hrs)
	// TODO should be rounded to .0000
	rv := yr + hrs/(daysinyear*24)
	//rv = math.Round(rv,4 places)
	return rv
}
