package distance

import (
	"math"
)

// Unit constants for distance calculation.
const (
	// UnitKilometers returns the result in meters using the kilometers conversion factor.
	UnitKilometers = "K"
)

// CalculateDistance computes the distance in meters between two geographic coordinates
// using the Haversine formula.
//
// By default, the result uses the statute miles conversion (* 1.1515) then multiplied by 1000.
// Pass UnitKilometers ("K") as a unit to use the kilometers conversion factor instead.
func CalculateDistance(lat1, lng1, lat2, lng2 float64, units ...string) float64 {
	radlat1 := math.Pi * lat1 / 180
	radlat2 := math.Pi * lat2 / 180

	theta := lng1 - lng2
	radtheta := math.Pi * theta / 180

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi * 60 * 1.1515

	if len(units) > 0 && units[0] == UnitKilometers {
		dist = dist * 1.609344
	}

	dist = dist * 1000

	return dist
}
