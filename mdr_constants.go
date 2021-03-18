// constants.go

// constants and conversion functions

package mdr

import (
	"math"
)

const ()

var (
	PiX2             = math.Pi * 2.0
	RadiansPerDegree = math.Pi / 180.0
	Km2miles         = 0.621371
	M2ft             = 3.28084          // meters -> feet
	TwoTo31th        = float64(1 << 31) // 2 ^ 31
	GaiaKm           = 6371.0           //  Earth's radius in kilometers.
)

// MiFromCm turns centimeters into miles
func MiFromCm(cm uint32) float64 {
	return (float64(cm) * Km2miles) / 100000
}

// MphFromMMs converts millimeters per second to miles per hour
func MphFromMMs(mms uint32) float64 {
	return (float64(mms) * 3600 * Km2miles) / 1000000
}

// MphFromMps converts meters per second to miles per hour
func MphFromMps(mps float64) float64 {
	return (float64(mps) * 3600 * Km2miles) / 1000
}

func Centigrade(f int16) int16 {
	return int16((float32(f) - 32) / 1.8)
}

func Celcius(c int16) int16 {
	return Centigrade(c)
}

func Fahrenheit(c int16) int16 {
	var rv float32
	rv = float32(c) * 1.8
	rv += 32
	return int16(rv)
}
