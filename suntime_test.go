// suntime_test.go

package suntime

import (
	"testing"
	"time"
)

/*
Day	Date	Astro Dawn	Naut Dawn	Dawn	    Sunrise	    Solar Noon	Sunset		Dusk		Naut Dusk	Astro Dusk
Tue	1/7/25	17:40:19	18:12:32	18:45:40	19:15:07	0:02:50		4:50:32		5:20:00		5:53:08		6:25:21
*/

// Flint Hill, MO
var testLongitude float64 = 90.85866
var testLatitude float64 = 38.85563244
var testDate time.Time = time.Date(2025, 1, 7, 0, 0, 0, 0, time.Now().Location())

// test converting testDate to Julian day
func TestToJulianDay(t *testing.T) {
	expected := 2460682.5000000
	result := ToJulianDay(testDate)
	if result != expected {
		t.Errorf("ToJulianDay() = %v, want %v", result, expected)
	}
}

func TestSunrise(t *testing.T) {
	julianDay := ToJulianDay(testDate)
	expected := time.Date(2025, 1, 7, 13, 19, 48, 0, time.UTC) // Example expected time

	result := Sunrise(julianDay, testLongitude, testLatitude)
	if !result.Equal(expected) {
		t.Errorf("Sunrise() = %v, want %v", result, expected)
	}
}

func TestSunset(t *testing.T) {
	julianDay := ToJulianDay(testDate)
	expected := time.Date(2025, 1, 7, 22, 56, 0, 0, time.UTC) // Example expected time

	result := Sunset(julianDay, testLongitude, testLatitude)
	if !result.Equal(expected) {
		t.Errorf("Sunset() = %v, want %v", result, expected)
	}
}

func TestCivilTwilightSunrise(t *testing.T) {
	julianDay := ToJulianDay(testDate)
	expected := time.Date(2025, 1, 7, 12, 50, 0, 0, time.UTC) // Example expected time

	result := CivilTwilightSunrise(julianDay, testLongitude, testLatitude)
	if !result.Equal(expected) {
		t.Errorf("CivilTwilightSunrise() = %v, want %v", result, expected)
	}
}

func TestCivilTwilightSunset(t *testing.T) {
	julianDay := ToJulianDay(testDate)
	expected := time.Date(2025, 1, 7, 23, 30, 0, 0, time.UTC) // Example expected time

	result := CivilTwilightSunset(julianDay, testLongitude, testLatitude)
	if !result.Equal(expected) {
		t.Errorf("CivilTwilightSunset() = %v, want %v", result, expected)
	}
}

func TestNauticalTwilightSunrise(t *testing.T) {
	julianDay := ToJulianDay(testDate)
	expected := time.Date(2025, 1, 7, 12, 20, 0, 0, time.UTC) // Example expected time

	result := NauticalTwilightSunrise(julianDay, testLongitude, testLatitude)
	if !result.Equal(expected) {
		t.Errorf("NauticalTwilightSunrise() = %v, want %v", result, expected)
	}
}

func TestNauticalTwilightSunset(t *testing.T) {
	julianDay := ToJulianDay(testDate)
	expected := time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC) // Example expected time

	result := NauticalTwilightSunset(julianDay, testLongitude, testLatitude)
	if !result.Equal(expected) {
		t.Errorf("NauticalTwilightSunset() = %v, want %v", result, expected)
	}
}

func TestAstronomicalTwilightSunrise(t *testing.T) {
	julianDay := ToJulianDay(testDate)
	expected := time.Date(2025, 1, 7, 11, 50, 0, 0, time.UTC) // Example expected time

	result := AstronomicalTwilightSunrise(julianDay, testLongitude, testLatitude)
	if !result.Equal(expected) {
		t.Errorf("AstronomicalTwilightSunrise() = %v, want %v", result, expected)
	}
}

func TestAstronomicalTwilightSunset(t *testing.T) {
	julianDay := ToJulianDay(testDate)
	expected := time.Date(2025, 1, 7, 0, 30, 0, 0, time.UTC) // Example expected time

	result := AstronomicalTwilightSunset(julianDay, testLongitude, testLatitude)
	if !result.Equal(expected) {
		t.Errorf("AstronomicalTwilightSunset() = %v, want %v", result, expected)
	}
}

func TestFromJulianDay(t *testing.T) {
	julianDay := float64(2460680.5)
	expected := testDate

	result := FromJulianDay(julianDay)
	if !result.Equal(expected) {
		t.Errorf("FromJulianDay() = %v, want %v", result, expected)
	}
}

func TestParseDMS(t *testing.T) {
	input := "38° 51' 31.44\" N"
	expectedDMS := DMS{Degrees: 38, Minutes: 51, Seconds: 31.44}
	expectedDirection := "N"

	resultDMS, resultDirection, err := ParseDMS(input)
	if err != nil {
		t.Errorf("ParseDMS() error = %v", err)
	}
	if resultDMS != expectedDMS || resultDirection != expectedDirection {
		t.Errorf(
			"ParseDMS() = %v, %v, want %v, %v", resultDMS, resultDirection, expectedDMS,
			expectedDirection,
		)
	}
}

func TestDMSToDecimal(t *testing.T) {
	dms := DMS{Degrees: 38, Minutes: 51, Seconds: 31.44}
	direction := "N"
	expected := 38.8587333

	result := DmsToDecimal(dms, direction)
	if result != expected {
		t.Errorf("DmsToDecimal() = %v, want %v", result, expected)
	}
}

func TestDecimalToDMS(t *testing.T) {
	decimal := 38.8587333
	isLatitude := true
	expectedDMS := DMS{Degrees: 38, Minutes: 51, Seconds: 31.44}
	expectedDirection := "N"

	resultDMS, resultDirection := DecimalToDMS(decimal, isLatitude)
	if resultDMS != expectedDMS || resultDirection != expectedDirection {
		t.Errorf(
			"DecimalToDMS() = %v, %v, want %v, %v", resultDMS, resultDirection, expectedDMS,
			expectedDirection,
		)
	}
}
