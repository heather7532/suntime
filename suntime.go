// suntime.go

// Sunrise calculates the time of sunrise for a given Julian day, longitude, and latitude.
// It returns the time in UTC.

package suntime

import (
	"math"
	"time"
)

const (
	J1970 = 2440588.0 // Julian date of the Unix epoch (1970-01-01)
	J2000 = 2451545.0 // Julian date for the epoch (2000-01-01)
)

// Sunrise calculates the sunrise time for a given Julian day, longitude, and latitude.
func Sunrise(julianDay, longitude, latitude float64) time.Time {
	return calculateTime(julianDay, longitude, latitude, 90.833, true)
}

// Sunset calculates the sunset time for a given Julian day, longitude, and latitude.
func Sunset(julianDay, longitude, latitude float64) time.Time {
	return calculateTime(julianDay, longitude, latitude, 90.833, false)
}

// CivilTwilightSunrise calculates the civil twilight sunrise time.
func CivilTwilightSunrise(julianDay, longitude, latitude float64) time.Time {
	return calculateTime(julianDay, longitude, latitude, 96.0, true) // 90° + 6°
}

// CivilTwilightSunset calculates the civil twilight sunset time.
func CivilTwilightSunset(julianDay, longitude, latitude float64) time.Time {
	return calculateTime(julianDay, longitude, latitude, 96.0, false)
}

// NauticalTwilightSunrise calculates the nautical twilight sunrise time.
func NauticalTwilightSunrise(julianDay, longitude, latitude float64) time.Time {
	return calculateTime(julianDay, longitude, latitude, 102.0, true) // 90° + 12°
}

// NauticalTwilightSunset calculates the nautical twilight sunset time.
func NauticalTwilightSunset(julianDay, longitude, latitude float64) time.Time {
	return calculateTime(julianDay, longitude, latitude, 102.0, false)
}

// AstronomicalTwilightSunrise calculates the astronomical twilight sunrise time.
func AstronomicalTwilightSunrise(julianDay, longitude, latitude float64) time.Time {
	return calculateTime(julianDay, longitude, latitude, 108.0, true) // 90° + 18°
}

// AstronomicalTwilightSunset calculates the astronomical twilight sunset time.
func AstronomicalTwilightSunset(julianDay, longitude, latitude float64) time.Time {
	return calculateTime(julianDay, longitude, latitude, 108.0, false)
}

// ToJulianDay converts a time.Time value to a Julian day.
func ToJulianDay(t time.Time) float64 {
	return float64(t.Unix())/86400.0 + J1970
}

// FromJulianDay converts a Julian day to a time.Time value.
func FromJulianDay(jd float64) time.Time {
	return time.Unix(int64((jd-J1970)*86400), 0).UTC()
}

// Helper functions
func solarDeclination(d float64) float64 {
	return math.Asin(math.Sin(-23.44*math.Pi/180) * math.Cos(2*math.Pi*(d+10)/365.0))
}

func hourAngle(lat, decl, angle float64, isSunrise bool) float64 {
	latRad := lat * math.Pi / 180
	declRad := decl
	h := math.Acos((math.Cos(angle*math.Pi/180) - math.Sin(latRad)*math.Sin(declRad)) /
		(math.Cos(latRad) * math.Cos(declRad)))
	if isSunrise {
		return -h
	}
	return h
}

func solarTransit(d, lng, h float64) float64 {
	return J2000 + d + h/(2*math.Pi)
}

func calculateTime(d, lng, lat, angle float64, isSunrise bool) time.Time {
	decl := solarDeclination(d)
	h := hourAngle(lat, decl, angle, isSunrise)
	transit := solarTransit(d, lng, h)
	return FromJulianDay(transit)
}
