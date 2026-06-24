package jetspotter

import (
	"math"
	"time"

	"github.com/jftuga/geodist"
)

// PredictCPA dead-reckons the aircraft along its current track at its current
// ground speed and returns the closest point of approach (CPA) to the target
// location within the look-ahead horizon.
//
// Returns:
//   - cpaAt:       the wall-clock time at which the CPA occurs
//   - cpaDistance: the lateral ground distance in kilometers between the
//                  projected aircraft position and the target at the CPA
//   - ok:          false if the aircraft is stationary or no CPA within the
//                  horizon is found (e.g. the aircraft is heading away and
//                  the distance only grows)
func PredictCPA(target geodist.Coord, ac AircraftRaw, horizon time.Duration) (cpaAt time.Time, cpaDistance int, ok bool) {
	// Stationary aircraft cannot be projected.
	if ac.GS <= 0 {
		return time.Time{}, 0, false
	}

	// Convert ground speed from knots to kilometers per hour, then to km/s.
	kmph := float64(ConvertKnotsToKilometersPerHour(int(ac.GS)))
	kmPerSec := kmph / 3600.0

	// Use a 5-second stepping granularity. Smaller steps improve accuracy
	// marginally at the cost of more iterations; 5s is a good balance for
	// aircraft that move ~0.1-0.3 km per step.
	const stepSec = 5.0
	steps := int(horizon.Seconds() / stepSec)
	if steps < 1 {
		steps = 1
	}

	// If no valid track is available, fall back to true heading, then mag heading.
	track := ac.Track
	if track == 0 {
		if ac.TrueHeading != 0 {
			track = ac.TrueHeading
		} else {
			track = ac.MagHeading
		}
	}

	lat := ac.Lat
	lon := ac.Lon

	now := time.Now().UTC()
	bestDist := CalculateDistance(target, geodist.Coord{Lat: lat, Lon: lon})
	bestAt := now

	for i := 1; i <= steps; i++ {
		elapsedSec := float64(i) * stepSec
		// Distance traveled along the track in this step.
		km := kmPerSec * elapsedSec
		// Advance the position using an equirectangular approximation.
		// This is consistent with how CalculateBearing treats coordinates.
		newLat, newLon := advancePosition(lat, lon, track, km)

		dist := CalculateDistance(target, geodist.Coord{Lat: newLat, Lon: newLon})
		if dist < bestDist {
			bestDist = dist
			bestAt = now.Add(time.Duration(elapsedSec*1000) * time.Millisecond)
		}

		// Early exit: if distance is growing for several consecutive steps we
		// have passed the CPA. Allow a small number of increasing steps to
		// handle minor jitter, then stop.
		if i > 2 && dist > bestDist {
			break
		}
	}

	// If the best distance is the starting distance (i.e. the distance never
	// decreased), the aircraft is heading away. Report no CPA.
	if bestDist >= CalculateDistance(target, geodist.Coord{Lat: lat, Lon: lon}) {
		return time.Time{}, 0, false
	}

	return bestAt, bestDist, true
}

// advancePosition moves a lat/lon position by `km` kilometers along the given
// bearing (in degrees) using an equirectangular approximation. This is the
// inverse of the bearing calculation and is accurate enough over the short
// distances and small steps used by PredictCPA.
func advancePosition(lat, lon, bearingDeg, km float64) (newLat, newLon float64) {
	// 1 degree of latitude ~= 111.32 km.
	const kmPerDeg = 111.32

	// Convert bearing to radians.
	brg := toRadians(bearingDeg)

	// North-south component (latitude change).
	dLat := (km * math.Cos(brg)) / kmPerDeg

	// East-west component (longitude change), scaled by the cosine of latitude
	// to account for convergence of meridians.
	cosLat := math.Cos(toRadians(lat))
	if cosLat < 1e-6 {
		cosLat = 1e-6
	}
	dLon := (km * math.Sin(brg)) / (kmPerDeg * cosLat)

	return lat + dLat, lon + dLon
}

// HeadingToCompassWord converts a heading in degrees (0-359) to an 8-point
// compass word: N, NE, E, SE, S, SW, W, NW.
func HeadingToCompassWord(heading float64) string {
	// Normalize to [0, 360)
	h := math.Mod(heading, 360.0)
	if h < 0 {
		h += 360
	}

	// Each compass sector is 45 degrees wide, centered on the cardinal point.
	// Sector boundaries: [0, 22.5) -> N, [22.5, 67.5) -> NE, etc.
	sectors := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	index := int((h + 22.5) / 45.0)
	if index >= len(sectors) {
		index = 0
	}
	return sectors[index]
}
