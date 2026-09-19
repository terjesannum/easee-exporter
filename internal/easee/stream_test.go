package easee

import (
	"encoding/json"
	"testing"
	"time"
)

func collectingStream() (*Stream, map[string]*ChargerState) {
	states := map[string]*ChargerState{}
	s := &Stream{apply: func(charger string, o *Observation) {
		if states[charger] == nil {
			states[charger] = &ChargerState{}
		}
		_ = states[charger].Apply(o)
	}}
	return s, states
}

// The stream sends every value as a string regardless of dataType, unlike the
// observations endpoint which sends native json types.
func TestHandleProductUpdateStringValues(t *testing.T) {
	s, states := collectingStream()

	s.handleProductUpdate(json.RawMessage(`{"mid":"ECRTU44X","id":109,"dataType":4,"timestamp":"2026-09-10T12:00:00Z","value":"3"}`))
	s.handleProductUpdate(json.RawMessage(`{"mid":"ECRTU44X","id":121,"dataType":3,"timestamp":"2026-09-10T12:00:01Z","value":"12.5"}`))
	s.handleProductUpdate(json.RawMessage(`{"mid":"ECRTU44X","id":102,"dataType":2,"timestamp":"2026-09-10T12:00:02Z","value":"1"}`))
	s.handleProductUpdate(json.RawMessage(`{"mid":"ECRTU44X","id":250,"dataType":2,"timestamp":"2026-09-10T12:00:03Z","value":"0"}`))

	st := states["ECRTU44X"]
	if st == nil {
		t.Fatal("no state for ECRTU44X")
	}
	if st.ChargerOpMode == nil || *st.ChargerOpMode != 3 {
		t.Errorf("ChargerOpMode = %v, want 3", st.ChargerOpMode)
	}
	if st.SessionEnergy == nil || *st.SessionEnergy != 12.5 {
		t.Errorf("SessionEnergy = %v, want 12.5", st.SessionEnergy)
	}
	if st.SmartCharging == nil || !*st.SmartCharging {
		t.Errorf("SmartCharging = %v, want true", st.SmartCharging)
	}
	if st.ConnectedToCloud == nil || *st.ConnectedToCloud {
		t.Errorf("ConnectedToCloud = %v, want false", st.ConnectedToCloud)
	}
}

// Observations arrive one at a time, so each must build on what is known
// rather than replace it.
func TestObservationsAccumulate(t *testing.T) {
	s, states := collectingStream()
	s.handleProductUpdate(json.RawMessage(`{"mid":"ECRTU44X","id":109,"dataType":4,"value":"3"}`))
	s.handleProductUpdate(json.RawMessage(`{"mid":"ECRTU44X","id":121,"dataType":3,"value":"12.5"}`))

	st := states["ECRTU44X"]
	if st.ChargerOpMode == nil {
		t.Error("earlier observation was lost when a later one arrived")
	}
	if st.SessionEnergy == nil {
		t.Error("later observation was not applied")
	}
}

func TestHandleProductUpdateBatchAndRouting(t *testing.T) {
	s, states := collectingStream()
	s.handleProductUpdate(json.RawMessage(`[
		{"mid":"ECRTU44X","id":109,"dataType":4,"value":"3"},
		{"mid":"EC8TCZYT","id":109,"dataType":4,"value":"1"}
	]`))

	if st := states["ECRTU44X"]; st == nil || st.ChargerOpMode == nil || *st.ChargerOpMode != 3 {
		t.Errorf("ECRTU44X opMode = %v, want 3", st)
	}
	if st := states["EC8TCZYT"]; st == nil || st.ChargerOpMode == nil || *st.ChargerOpMode != 1 {
		t.Errorf("EC8TCZYT opMode = %v, want 1", st)
	}
}

// Voltage is derived, so it must keep working when the in voltages arrive
// separately over time rather than together in one response.
func TestStreamDerivesVoltageIncrementally(t *testing.T) {
	s, states := collectingStream()
	s.handleProductUpdate(json.RawMessage(`{"mid":"ECRTU44X","id":195,"dataType":3,"value":"239.47"}`))
	if v := states["ECRTU44X"].Voltage; v == nil || *v != 239.47 {
		t.Fatalf("Voltage = %v, want 239.47", v)
	}
	s.handleProductUpdate(json.RawMessage(`{"mid":"ECRTU44X","id":194,"dataType":3,"value":"246.225"}`))
	if v := states["ECRTU44X"].Voltage; v == nil || *v != 246.225 {
		t.Errorf("Voltage = %v, want 246.225 after a higher reading arrived", v)
	}
}

func TestTimestamplessObservationCountsAsNow(t *testing.T) {
	s, states := collectingStream()
	before := time.Now()
	s.handleProductUpdate(json.RawMessage(`{"mid":"ECRTU44X","id":109,"dataType":4,"value":"3"}`))
	p := states["ECRTU44X"].LatestPulse
	if p == nil || p.Before(before) {
		t.Errorf("LatestPulse = %v, want the receive time", p)
	}
}

// Junk on the stream must not take the process down or corrupt known state.
func TestMalformedUpdatesAreSurvivable(t *testing.T) {
	s, states := collectingStream()
	s.handleProductUpdate(json.RawMessage(`{"mid":"ECRTU44X","id":109,"dataType":4,"value":"3"}`))

	for _, junk := range []string{
		`not json at all`,
		`{}`,
		`null`,
		`[]`,
		`{"mid":"ECRTU44X","id":120,"dataType":3,"value":"banana"}`,
		`{"id":109,"dataType":4,"value":"1"}`,
		`{"mid":"ECRTU44X","id":999999,"dataType":4,"value":"1"}`,
	} {
		s.handleProductUpdate(json.RawMessage(junk))
	}

	st := states["ECRTU44X"]
	if st.ChargerOpMode == nil || *st.ChargerOpMode != 3 {
		t.Errorf("known state corrupted by junk input: %v", st.ChargerOpMode)
	}
	if st.TotalPower != nil {
		t.Errorf("TotalPower = %v, want nil after an unparseable value", st.TotalPower)
	}
}
