// suntime.go

// Sunrise calculates the time of sunrise for a given Julian day, longitude, and latitude.
// It returns the time in UTC.

package suntime

import (
	"fmt"
	"github.com/kelvins/sunrisesunset"
	"github.com/soniakeys/meeus/v3/julian"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	J1970              = 2440588.0 // Julian date of the Unix epoch (1970-01-01)
	J2000              = 2451545.0 // Julian date for the epoch (2000-01-01)
	DegreesToRadians   = math.Pi / 180.0
	RadiansToDegrees   = 180.0 / math.Pi
	MeanAnomalyCoeff   = 0.98560028
	MeanAnomalyBase    = 357.5291
	CenterCoeff1       = 1.9148
	CenterCoeff2       = 0.0200
	CenterCoeff3       = 0.0003
	EclipticLongBase   = 102.9372
	Obliquity          = 23.44
	SolarTransitCoeff1 = 0.0053
	SolarTransitCoeff2 = 0.0069
)

// Sunrise calculates the sunrise time for a given Julian day, longitude, and latitude.
func Sunrise(julianDay, longitude, latitude float64) time.Time {
	// Convert the input Julian day to UTC
	return calculateTime(JulianToUTC(julianDay), longitude, latitude, 90.833, true)
}

// Sunset calculates the sunset time for a given Julian day, longitude, and latitude.
func Sunset(julianDay, longitude, latitude float64) time.Time {
	return calculateTime(JulianToUTC(julianDay), longitude, latitude, 90.833, false)
}

// CivilTwilightSunrise calculates the civil twilight sunrise time.
func CivilTwilightSunrise(julianDay, longitude, latitude float64) time.Time {
	return calculateTime(JulianToUTC(julianDay), longitude, latitude, 96.0, true) // 90° + 6°
}

// CivilTwilightSunset calculates the civil twilight sunset time.
func CivilTwilightSunset(julianDay, longitude, latitude float64) time.Time {
	return calculateTime(JulianToUTC(julianDay), longitude, latitude, 96.0, false)
}

// NauticalTwilightSunrise calculates the nautical twilight sunrise time.
func NauticalTwilightSunrise(julianDay, longitude, latitude float64) time.Time {
	return calculateTime(JulianToUTC(julianDay), longitude, latitude, 102.0, true) // 90° + 12°
}

// NauticalTwilightSunset calculates the nautical twilight sunset time.
func NauticalTwilightSunset(julianDay, longitude, latitude float64) time.Time {
	return calculateTime(JulianToUTC(julianDay), longitude, latitude, 102.0, false)
}

// AstronomicalTwilightSunrise calculates the astronomical twilight sunrise time.
func AstronomicalTwilightSunrise(julianDay, longitude, latitude float64) time.Time {
	return calculateTime(JulianToUTC(julianDay), longitude, latitude, 108.0, true) // 90° + 18°
}

// AstronomicalTwilightSunset calculates the astronomical twilight sunset time.
func AstronomicalTwilightSunset(julianDay, longitude, latitude float64) time.Time {
	return calculateTime(JulianToUTC(julianDay), longitude, latitude, 108.0, false)
}

func JulianToUTC(julian float64) float64 {
	// Shift the Julian day to align with midnight UTC instead of noon UTC
	julianMidnight := julian + 0.5
	utcTime := FromJulianDay(julianMidnight).UTC()
	return ToJulianDay(utcTime)
}

// ToJulianDay converts a time.Time value to a Julian day.
func ToJulianDay(t time.Time) float64 {
	date := t.UTC()
	jd := julian.CalendarGregorianToJD(date.Year(), int(date.Month()), float64(date.Day()))
	return jd
}

// FromJulianDay converts a Julian day to a time.Time value.
func FromJulianDay(jd float64) time.Time {
	return julian.JDToTime(jd)
}

// Helper functions
func solarDeclination(d float64) float64 {
	return math.Asin(math.Sin(-23.44*math.Pi/180) * math.Cos(2*math.Pi*(d+10)/365.0))
}

func hourAngle(lat, decl, angle float64, isSunrise bool) float64 {
	latRad := lat * math.Pi / 180
	declRad := decl
	h := math.Acos(
		(math.Cos(angle*math.Pi/180) - math.Sin(latRad)*math.Sin(declRad)) /
			(math.Cos(latRad) * math.Cos(declRad)),
	)
	if isSunrise {
		return -h
	}
	return h
}

func solarTransit(d, lng, h float64) float64 {
	return J2000 + d + h/(2*math.Pi)
}

func calculateTime(julianDay, longitude, latitude, angle float64, isSunrise bool) time.Time {
	// Calculate the number of days since J2000.0
	n := julianDay - J2000

	// Calculate the mean solar noon
	Jstar := n - longitude/360.0

	// Calculate the solar mean anomaly
	M := (357.5291 + 0.98560028*Jstar) * DegreesToRadians

	// Calculate the equation of the center
	C := 1.9148*math.Sin(M) + 0.0200*math.Sin(2*M) + 0.0003*math.Sin(3*M)

	// Calculate the ecliptic longitude
	lambda := (M + C + 102.9372*DegreesToRadians + math.Pi) * RadiansToDegrees

	// Calculate the solar transit
	Jtransit := J2000 + Jstar + 0.0053*math.Sin(M) - 0.0069*math.Sin(2*lambda*DegreesToRadians)

	// Calculate the declination of the sun
	delta := math.Asin(math.Sin(lambda*DegreesToRadians) * math.Sin(23.44*DegreesToRadians))

	// Calculate the hour angle
	latRad := latitude * DegreesToRadians
	declRad := delta
	h := math.Acos(
		(math.Cos(angle*DegreesToRadians) - math.Sin(latRad)*math.Sin(declRad)) /
			(math.Cos(latRad) * math.Cos(declRad)),
	)
	if isSunrise {
		h = -h
	}

	// Calculate the sunrise or sunset time
	Jset := Jtransit + h/(2*math.Pi)

	// Correct for Julian day noon offset
	return FromJulianDay(Jset).Round(time.Second)
}

// Convert time from utc
func ConvertTimeFromUTC(t time.Time, offset int) time.Time {
	return t.Add(time.Duration(offset) * time.Hour)
}

// ParseDMS parses a DMS string into a DMS struct and direction
func ParseDMS(input string) (DMS, string, error) {
	// Regular expression to match DMS format
	re := regexp.MustCompile(`^(\d{1,2})°\s+(\d{1,2})'\s+(\d{1,2}(?:\.\d+)?)"\s+([NSEW])$`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return DMS{}, "", fmt.Errorf("invalid DMS format: %s", input)
	}

	// Extract components
	degrees, _ := strconv.Atoi(matches[1])
	minutes := 0
	if matches[2] != "" {
		minutes, _ = strconv.Atoi(matches[2])
	}
	seconds := 0.0
	if matches[3] != "" {
		seconds, _ = strconv.ParseFloat(matches[3], 64)
	}
	direction := strings.ToUpper(matches[4])

	return DMS{Degrees: degrees, Minutes: minutes, Seconds: seconds}, direction, nil
}

// Function: Convert DMS to Decimal Degrees
func DmsToDecimal(dms DMS, direction string) float64 {
	// Convert DMS to Decimal Degrees
	decimal := float64(dms.Degrees) + float64(dms.Minutes)/60 + dms.Seconds/3600

	// Adjust for direction (N/S/E/W)
	switch direction {
	case "S", "W":
		decimal = -decimal
	case "N", "E":
		decimal = decimal
	default:
		fmt.Println("Invalid direction. Use N, S, E, or W.")
	}

	return roundToPlaces(decimal, 7)
}

// Function: Convert Decimal Degrees to DMS
func DecimalToDMS(decimal float64, isLatitude bool) (DMS, string) {
	// Determine the direction
	direction := ""
	if isLatitude {
		if decimal < 0 {
			direction = "S"
			decimal = math.Abs(decimal)
		} else {
			direction = "N"
		}
	} else {
		if decimal < 0 {
			direction = "W"
			decimal = math.Abs(decimal)
		} else {
			direction = "E"
		}
	}

	// Extract degrees, minutes, and seconds
	degrees := int(decimal)
	minutes := int((decimal - float64(degrees)) * 60)
	seconds := roundToPlaces((decimal-float64(degrees))*60-float64(minutes), 4) * 60

	return DMS{Degrees: degrees, Minutes: minutes, Seconds: seconds}, direction
}

// roundToPlaces rounds a float64 to the specified number of decimal places
func roundToPlaces(value float64, places int) float64 {
	factor := math.Pow(10, float64(places))
	return math.Round(value*factor) / factor
}
