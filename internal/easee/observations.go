package easee

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Observation is a single value reported by a charger, as returned by
// /state/{serialNumber}/observations. Value is typed by DataType
// (2=Boolean, 3=Double, 4=Integer, 5=Position, 6=String), so it is kept raw
// and decoded by the field it maps to.
type Observation struct {
	Id        int             `json:"id"`
	Timestamp time.Time       `json:"timestamp"`
	DataType  int             `json:"dataType"`
	Value     json.RawMessage `json:"value"`
}

type observationsResponse struct {
	Observations []Observation `json:"observations"`
}

func (o *Observation) decodeBool(dst **bool) error {
	var v bool
	if err := json.Unmarshal(o.Value, &v); err != nil {
		return fmt.Errorf("observation %d: %w", o.Id, err)
	}
	*dst = &v
	return nil
}

// decodeInt and decodeFloat both go through float64 so that an observation
// documented as Integer still decodes if the API reports it as a double.
func (o *Observation) decodeInt(dst **int) error {
	var v float64
	if err := json.Unmarshal(o.Value, &v); err != nil {
		return fmt.Errorf("observation %d: %w", o.Id, err)
	}
	i := int(v)
	*dst = &i
	return nil
}

func (o *Observation) decodeFloat(dst **float64) error {
	var v float64
	if err := json.Unmarshal(o.Value, &v); err != nil {
		return fmt.Errorf("observation %d: %w", o.Id, err)
	}
	*dst = &v
	return nil
}

type observationDecoder func(*ChargerState, *Observation) error

// chargerStateObservations maps the observation ids that replace the removed
// /api/chargers/{id}/state endpoint onto the ChargerState fields they feed.
// The ids requested from the API are derived from these keys, so the two
// cannot drift apart.
var chargerStateObservations = map[int]observationDecoder{
	30: func(s *ChargerState, o *Observation) error { return o.decodeBool(&s.LockCablePermanently) },      // LOCK CABLE PERMANENTLY
	46: func(s *ChargerState, o *Observation) error { return o.decodeInt(&s.LedMode) },                    // LEDMODE
	48: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.DynamicChargerCurrent) },    // DYNAMIC CHARGER CURRENT
	50: func(s *ChargerState, o *Observation) error { return o.decodeInt(&s.OfflineMaxCircuitCurrentP1) }, // MAX CURRENT OFFLINE FALLBACK P1
	51: func(s *ChargerState, o *Observation) error { return o.decodeInt(&s.OfflineMaxCircuitCurrentP2) }, // MAX CURRENT OFFLINE FALLBACK P2
	52: func(s *ChargerState, o *Observation) error { return o.decodeInt(&s.OfflineMaxCircuitCurrentP3) }, // MAX CURRENT OFFLINE FALLBACK P3
	70: func(s *ChargerState, o *Observation) error {
		return o.decodeFloat(&s.CircuitTotalAllocatedPhaseConductorCurrentL1)
	}, // Circuit Total Allocated Phase Conductor Current L1
	71: func(s *ChargerState, o *Observation) error {
		return o.decodeFloat(&s.CircuitTotalAllocatedPhaseConductorCurrentL2)
	}, // ... L2
	72: func(s *ChargerState, o *Observation) error {
		return o.decodeFloat(&s.CircuitTotalAllocatedPhaseConductorCurrentL3)
	}, // ... L3
	73: func(s *ChargerState, o *Observation) error {
		return o.decodeFloat(&s.CircuitTotalPhaseConductorCurrentL1)
	}, // Circuit Total Phase Conductor Current L1
	74: func(s *ChargerState, o *Observation) error {
		return o.decodeFloat(&s.CircuitTotalPhaseConductorCurrentL2)
	}, // ... L2
	75: func(s *ChargerState, o *Observation) error {
		return o.decodeFloat(&s.CircuitTotalPhaseConductorCurrentL3)
	}, // ... L3
	80:  func(s *ChargerState, o *Observation) error { return o.decodeInt(&s.ChargerFirmware) },           // SOFTWARE RELEASE
	96:  func(s *ChargerState, o *Observation) error { return o.decodeInt(&s.ReasonForNoCurrent) },        // REASON FOR NO CURRENT
	102: func(s *ChargerState, o *Observation) error { return o.decodeBool(&s.SmartCharging) },            // SMART CHARGING
	103: func(s *ChargerState, o *Observation) error { return o.decodeBool(&s.CableLocked) },              // CABLE LOCKED
	104: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.CableRating) },             // CABLE RATING
	109: func(s *ChargerState, o *Observation) error { return o.decodeInt(&s.ChargerOpMode) },             // CHARGER OP MODE
	110: func(s *ChargerState, o *Observation) error { return o.decodeInt(&s.OutputPhase) },               // OUTPUT PHASE
	111: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.DynamicCircuitCurrentP1) }, // Dynamic Circuit Current P1
	112: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.DynamicCircuitCurrentP2) }, // Dynamic Circuit Current P2
	113: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.DynamicCircuitCurrentP3) }, // Dynamic Circuit Current P3
	114: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.OutputCurrent) },           // Output Current
	115: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.DeratedCurrent) },          // Derated Current
	116: func(s *ChargerState, o *Observation) error { return o.decodeBool(&s.DeratingActive) },           // DERATING ACTIVE
	119: func(s *ChargerState, o *Observation) error { return o.decodeInt(&s.ErrorCode) },                 // ERROR CODE
	120: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.TotalPower) },              // TOTAL POWER
	121: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.SessionEnergy) },           // SESSION ENERGY
	122: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.EnergyPerHour) },           // ENERGY PER HOUR
	124: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.LifeTimeEnergy) },          // LIFETIME ENERGY
	130: func(s *ChargerState, o *Observation) error { return o.decodeInt(&s.CellRSSI) },                  // CELL RSSI
	131: func(s *ChargerState, o *Observation) error { return o.decodeInt(&s.ChargerRAT) },                // CELL RAT
	132: func(s *ChargerState, o *Observation) error { return o.decodeInt(&s.WifiRSSI) },                  // WI FI RSSI
	136: func(s *ChargerState, o *Observation) error { return o.decodeInt(&s.LocalRSSI) },                 // LOCAL RSSI
	182: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.InCurrentT2) },             // INT CURRENT T2
	183: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.InCurrentT3) },             // INT CURRENT T3
	184: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.InCurrentT4) },             // INT CURRENT T4
	185: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.InCurrentT5) },             // INT CURRENT T5
	190: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.InVoltageT1T2) },           // IN VOLT T1T2
	191: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.InVoltageT1T3) },           // IN VOLT T1T3
	192: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.InVoltageT1T4) },           // IN VOLT T1T4
	193: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.InVoltageT1T5) },           // IN VOLT T1T5
	194: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.InVoltageT2T3) },           // IN VOLT T2T3
	195: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.InVoltageT2T4) },           // IN VOLT T2T4
	196: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.InVoltageT2T5) },           // IN VOLT T2T5
	197: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.InVoltageT3T4) },           // IN VOLT T3T4
	198: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.InVoltageT3T5) },           // IN VOLT T3T5
	199: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.InVoltageT4T5) },           // IN VOLT T4T5
	230: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.EqAvailableCurrentP1) },    // EQ AVAILABLE CURRENT P1
	231: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.EqAvailableCurrentP2) },    // EQ AVAILABLE CURRENT P2
	232: func(s *ChargerState, o *Observation) error { return o.decodeFloat(&s.EqAvailableCurrentP3) },    // EQ AVAILABLE CURRENT P3
	250: func(s *ChargerState, o *Observation) error { return o.decodeBool(&s.ConnectedToCloud) },         // CONNECTED TO CLOUD
}

// chargerStateObservationIds is the comma separated ids query parameter for
// all observations in chargerStateObservations, so a charger's full state is
// fetched in a single request.
var chargerStateObservationIds = observationIds(chargerStateObservations)

func observationIds(observations map[int]observationDecoder) string {
	ids := make([]int, 0, len(observations))
	for id := range observations {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	s := make([]string, len(ids))
	for i, id := range ids {
		s[i] = strconv.Itoa(id)
	}
	return strings.Join(s, ",")
}

// newChargerState folds an observations response into a ChargerState. A value
// that fails to decode is skipped rather than failing the whole update, so one
// unexpected observation cannot blank out every metric for a charger.
func newChargerState(res *observationsResponse) ChargerState {
	var state ChargerState
	for i := range res.Observations {
		o := &res.Observations[i]
		decode, ok := chargerStateObservations[o.Id]
		if !ok {
			continue
		}
		if err := decode(&state, o); err != nil {
			log.Printf("Ignoring observation: %v\n", err)
			continue
		}
		if state.LatestPulse == nil || o.Timestamp.After(*state.LatestPulse) {
			t := o.Timestamp
			state.LatestPulse = &t
		}
	}
	state.Voltage = maxInVoltage(&state)
	return state
}

// Voltage has no observation of its own. The removed state endpoint reported
// the highest of the input voltage measurements: checked against a week of
// recorded data for seven chargers, the old value equalled the maximum in
// every sample. Take the maximum of whichever pairs a charger reports rather
// than a fixed pair, since which terminals are wired varies by installation.
func maxInVoltage(s *ChargerState) *float64 {
	var max *float64
	for _, v := range []*float64{
		s.InVoltageT1T2, s.InVoltageT1T3, s.InVoltageT1T4, s.InVoltageT1T5,
		s.InVoltageT2T3, s.InVoltageT2T4, s.InVoltageT2T5,
		s.InVoltageT3T4, s.InVoltageT3T5, s.InVoltageT4T5,
	} {
		if v != nil && (max == nil || *v > *max) {
			max = v
		}
	}
	if max == nil {
		return nil
	}
	v := *max
	return &v
}
