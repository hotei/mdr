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
	hrs := (dayofyear-1)*24 + float64(when.Hour())
	//fmt.Printf("hrs = %v \n",hrs)
	return yr + hrs/(daysinyear*24)
}
