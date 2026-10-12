package jetspotter

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"jetspotter/internal/configuration"

	"github.com/jftuga/geodist"
)

// candidateTrackers holds the active confirmation goroutines keyed by ICAO.
// This prevents spawning duplicate goroutines for the same aircraft and lets
// the main loop find and signal trackers when needed.
var (
	candidateTrackers   = make(map[string]chan struct{})
	candidateTrackersMu sync.Mutex

	// OverheadNotifier is called when a candidate is confirmed. It is set by
	// cmd/jetspotter at startup to route to the configured notification
	// channels. Using a callback avoids an import cycle between the
	// jetspotter and notification packages.
	OverheadNotifier func([]OverheadCandidate, configuration.Config)
)

// EvaluateSample decides whether a single sub-loop sample confirms or rejects
// the candidate. It is a pure function so it can be unit-tested without any
// network I/O.
//
// Returns:
//   - stillOnCourse: true if the aircraft is still inbound and its predicted
//                   CPA remains within the overhead radius
//   - confirmed:     true if the cumulative on-course time now meets the
//                   configured confirmation threshold
func EvaluateSample(target geodist.Coord, raw AircraftRaw, candidate *OverheadCandidate, config configuration.Config) (stillOnCourse, confirmed bool) {
	// If the aircraft is on the ground or has no ground speed, it can no
	// longer be inbound.
	if raw.GS <= 0 {
		candidate.ConsecutiveOnCourseSeconds = 0
		return false, false
	}

	// Check whether the aircraft is still heading toward the target within
	// the configured (tight) inbound margin.
	if !IsAircraftInbound(target, raw, float64(config.OverheadInboundMarginDegrees)) {
		candidate.ConsecutiveOnCourseSeconds = 0
		return false, false
	}

	// Re-predict the CPA using the fresh sample.
	horizon := time.Duration(config.OverheadLookAheadMinutes) * time.Minute
	cpaAt, cpaDist, ok := PredictCPA(target, raw, horizon)
	if !ok || cpaDist > config.OverheadRadiusKilometers {
		candidate.ConsecutiveOnCourseSeconds = 0
		return false, false
	}

	// Still on course: accumulate the confirmation time.
	candidate.ConsecutiveOnCourseSeconds += config.OverheadSubloopPollSeconds
	candidate.PredictedCPATime = cpaAt
	candidate.CPADistanceKm = cpaDist
	candidate.CompassWord = HeadingToCompassWord(raw.Track)
	candidate.MinutesUntilOverhead = int(math.Round(time.Until(cpaAt).Minutes()))
	candidate.LastSampleTime = time.Now()
	candidate.Aircraft = updateAircraftFromRaw(candidate.Aircraft, raw)

	confirmed = candidate.ConsecutiveOnCourseSeconds >= config.OverheadConfirmSeconds
	return true, confirmed
}

// updateAircraftFromRaw refreshes the position-derived fields of an Aircraft
// snapshot from a fresh AircraftRaw sample, without re-fetching flight route
// or image data (which are stable and expensive to re-query).
func updateAircraftFromRaw(ac Aircraft, raw AircraftRaw) Aircraft {
	ac.Latitude = raw.Lat
	ac.Longitude = raw.Lon
	ac.Heading = raw.Track
	ac.Speed = int(raw.GS)
	if raw.AltBaro != nil {
		if f, ok := raw.AltBaro.(float64); ok {
			ac.Altitude = f
		}
	}
	ac.OnGround = ac.Altitude == 0
	return ac
}

// HandleOverheadPrediction scans the full in-range aircraft list for new
// overhead-prediction candidates and spawns a confirmation goroutine for
// each one. It is called from the main loop after HandleAircraft when
// OVERHEAD_PREDICTION_ENABLED is true.
//
// allAircraftInRange is the full list of aircraft within MaxScanRangeKilometers
// (i.e. SpottedAircraft.Aircraft), which is larger than the notification range
// so that we can detect candidates before they enter the close-range zone.
func HandleOverheadPrediction(allAircraftInRange []Aircraft, config configuration.Config) {
	// Prune candidates whose CPA has passed or that were abandoned, so the
	// UI and API reflect only live candidates.
	pruneCandidates(config)

	// Apply the same aircraft type filter as the "spotted" notification path
	// so overhead predictions respect the AIRCRAFT_TYPES configuration.
	filteredAircraft := filterAircraftByTypes(allAircraftInRange, config.AircraftTypes)

	for _, ac := range filteredAircraft {
		// Skip aircraft on the ground: they cannot be inbound.
		if ac.OnGround {
			continue
		}

		// Skip aircraft that are already being tracked.
		if isTracked(ac.ICAO) {
			continue
		}

		// Reconstruct an AircraftRaw-like view from the Aircraft snapshot for
		// the inbound check and CPA prediction.
		raw := aircraftToRaw(ac)

		// Tighter inbound margin than the existing 30-degree default.
		if !IsAircraftInbound(config.Location, raw, float64(config.OverheadInboundMarginDegrees)) {
			continue
		}

		horizon := time.Duration(config.OverheadLookAheadMinutes) * time.Minute
		cpaAt, cpaDist, ok := PredictCPA(config.Location, raw, horizon)
		if !ok {
			continue
		}
		if cpaDist > config.OverheadRadiusKilometers {
			continue
		}

		// Register the candidate and spawn a confirmation tracker.
		candidate := OverheadCandidate{
			ICAO:                       ac.ICAO,
			Aircraft:                   ac,
			PredictedCPATime:           cpaAt,
			CPADistanceKm:              cpaDist,
			CompassWord:                HeadingToCompassWord(raw.Track),
			MinutesUntilOverhead:       int(math.Round(time.Until(cpaAt).Minutes())),
			ConfirmationState:          CandidatePending,
			ConsecutiveOnCourseSeconds: 0,
			LastSampleTime:            time.Now(),
			Notified:                  false,
		}

		upsertCandidate(candidate)
		spawnConfirmationTracker(candidate, config)
	}
}

// aircraftToRaw builds a minimal AircraftRaw from an Aircraft snapshot so the
// existing IsAircraftInbound and PredictCPA functions (which operate on
// AircraftRaw) can be reused.
func aircraftToRaw(ac Aircraft) AircraftRaw {
	altBaro := interface{}(ac.Altitude)
	if ac.OnGround {
		altBaro = float64(0)
	}
	return AircraftRaw{
		ICAO:         ac.ICAO,
		Callsign:     ac.Callsign,
		Registration: ac.Registration,
		PlaneType:    ac.Type,
		Desc:         ac.Description,
		AltBaro:      altBaro,
		GS:           float64(ac.Speed),
		Track:        ac.Heading,
		Lat:          ac.Latitude,
		Lon:          ac.Longitude,
	}
}

// isTracked reports whether a confirmation goroutine is already running for
// the given ICAO.
func isTracked(icao string) bool {
	candidateTrackersMu.Lock()
	defer candidateTrackersMu.Unlock()
	_, ok := candidateTrackers[icao]
	return ok
}

// spawnConfirmationTracker starts a per-candidate goroutine that polls the
// ADS-B API at the configured sub-loop interval until the candidate is
// confirmed, abandoned, or its CPA time has passed.
func spawnConfirmationTracker(candidate OverheadCandidate, config configuration.Config) {
	candidateTrackersMu.Lock()
	if _, exists := candidateTrackers[candidate.ICAO]; exists {
		candidateTrackersMu.Unlock()
		return
	}
	stop := make(chan struct{}, 1)
	candidateTrackers[candidate.ICAO] = stop
	candidateTrackersMu.Unlock()

	go func() {
		ticker := time.NewTicker(time.Duration(config.OverheadSubloopPollSeconds) * time.Second)
		defer ticker.Stop()

		missStreak := 0
		maxMisses := 6 // tolerate rate-limited/failed samples before abandoning

		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				done := sampleCandidate(candidate.ICAO, config, &missStreak, maxMisses)
				if done {
					return
				}
			}
		}
	}()
}

// sampleCandidate fetches a fresh position for the candidate's ICAO from the
// ADS-B API and evaluates it. Returns true if the tracker should stop
// (candidate confirmed+notified, abandoned, or passed its CPA).
func sampleCandidate(icao string, config configuration.Config, missStreak *int, maxMisses int) (stop bool) {
	raw, err := fetchAircraftByICAO(icao, config.Location, config.MaxScanRangeKilometers)
	if err != nil {
		log.Printf("Overhead tracker: error fetching %s: %v", icao, err)
		*missStreak++
		if *missStreak > maxMisses {
			markCandidateState(icao, CandidateAbandoned)
			untrack(icao)
			return true
		}
		return false
	}

	cand := getCandidate(icao)
	if cand == nil {
		// Candidate was pruned by the main loop; stop the tracker.
		untrack(icao)
		return true
	}

	stillOnCourse, confirmed := EvaluateSample(config.Location, raw, cand, config)

	if !stillOnCourse {
		*missStreak++
		cand.ConsecutiveOnCourseSeconds = 0
		if *missStreak > maxMisses {
			markCandidateState(icao, CandidateAbandoned)
			untrack(icao)
			return true
		}
		upsertCandidate(*cand)
		return false
	}

	*missStreak = 0

	// If the CPA time has passed, mark as passed and stop.
	if time.Now().After(cand.PredictedCPATime.Add(30 * time.Second)) {
		markCandidateState(icao, CandidatePassed)
		untrack(icao)
		return true
	}

	if confirmed && !cand.Notified {
		cand.ConfirmationState = CandidateConfirmed
		cand.Notified = true
		upsertCandidate(*cand)

		// Send the "look up in N minutes" notification via the registered
		// notifier callback.
		if OverheadNotifier != nil {
			OverheadNotifier([]OverheadCandidate{*cand}, config)
		} else {
			log.Printf("Overhead candidate %s confirmed but no notifier is registered", icao)
		}
	} else {
		upsertCandidate(*cand)
	}

	return false
}

// fetchAircraftByICAO queries the ADS-B API for a single aircraft by ICAO.
// It reuses the point/lat/lon/radius endpoint (which is known to exist on
// both api.adsb.one and api.adsb.lol) with the configured scan range, then
// filters the response down to the requested ICAO. This avoids depending on
// a /icao/{hex} endpoint that may not exist on all providers.
func fetchAircraftByICAO(icao string, location geodist.Coord, scanRangeKm int) (AircraftRaw, error) {
	miles := convertKilometersToNauticalMiles(float64(scanRangeKm))
	endpoint, err := url.JoinPath(baseURL, "point",
		strconv.FormatFloat(location.Lat, 'f', -1, 64),
		strconv.FormatFloat(location.Lon, 'f', -1, 64),
		strconv.Itoa(miles))
	if err != nil {
		return AircraftRaw{}, err
	}

	all, err := getAllAircrafRawInRangeFromURL(endpoint)
	if err != nil {
		return AircraftRaw{}, err
	}

	for _, ac := range all {
		if ac.ICAO == icao {
			return ac, nil
		}
	}
	return AircraftRaw{}, fmt.Errorf("aircraft %s not found in range", icao)
}

// getAllAircrafRawInRangeFromURL fetches and parses aircraft data from a
// fully-constructed ADS-B API URL. It is a small refactor of
// getAllAircrafRawInRange that allows callers to reuse the parsing logic
// with a custom URL (used here by the per-ICAO confirmation tracker).
func getAllAircrafRawInRangeFromURL(endpoint string) ([]AircraftRaw, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		if res.StatusCode == 429 {
			return nil, fmt.Errorf("API rate limit exceeded: %s", res.Status)
		}
		return nil, fmt.Errorf("API call to %s returned error: %s", endpoint, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var flightData FlightData
	if err := json.Unmarshal(body, &flightData); err != nil {
		return nil, fmt.Errorf("failed to parse API response from %s: %w", baseURL, err)
	}

	return flightData.AC, nil
}

// --- Candidate store helpers (thin wrappers around OverheadCandidates) ---

func upsertCandidate(c OverheadCandidate) {
	OverheadCandidates.Lock()
	defer OverheadCandidates.Unlock()
	for i := range OverheadCandidates.Candidates {
		if OverheadCandidates.Candidates[i].ICAO == c.ICAO {
			OverheadCandidates.Candidates[i] = c
			return
		}
	}
	OverheadCandidates.Candidates = append(OverheadCandidates.Candidates, c)
}

func getCandidate(icao string) *OverheadCandidate {
	OverheadCandidates.Lock()
	defer OverheadCandidates.Unlock()
	for i := range OverheadCandidates.Candidates {
		if OverheadCandidates.Candidates[i].ICAO == icao {
			c := OverheadCandidates.Candidates[i]
			return &c
		}
	}
	return nil
}

func markCandidateState(icao, state string) {
	OverheadCandidates.Lock()
	defer OverheadCandidates.Unlock()
	for i := range OverheadCandidates.Candidates {
		if OverheadCandidates.Candidates[i].ICAO == icao {
			OverheadCandidates.Candidates[i].ConfirmationState = state
			return
		}
	}
}

func untrack(icao string) {
	candidateTrackersMu.Lock()
	defer candidateTrackersMu.Unlock()
	delete(candidateTrackers, icao)
}

// pruneCandidates removes candidates whose CPA time has passed by more than
// 60 seconds, or that were abandoned. This keeps the /api/overhead response
// focused on live candidates.
func pruneCandidates(config configuration.Config) {
	OverheadCandidates.Lock()
	defer OverheadCandidates.Unlock()

	kept := OverheadCandidates.Candidates[:0]
	for _, c := range OverheadCandidates.Candidates {
		if c.ConfirmationState == CandidateAbandoned {
			continue
		}
		if c.ConfirmationState == CandidatePassed && time.Now().After(c.PredictedCPATime.Add(60*time.Second)) {
			continue
		}
		// Refresh minutes-until-overhead for the UI.
		if !c.PredictedCPATime.IsZero() {
			c.MinutesUntilOverhead = int(math.Round(time.Until(c.PredictedCPATime).Minutes()))
		}
		kept = append(kept, c)
	}
	OverheadCandidates.Candidates = kept
}
