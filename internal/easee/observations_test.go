package easee

import (
	"encoding/json"
	"testing"
	"time"
)

// Response shaped like the documented /state/{serialNumber}/observations
// payload: value is typed by dataType, and includes observations the exporter
// does not map.
const observationsJson = `{
  "observations": [
    {"value": true,  "id": 102, "timestamp": "2026-09-04T11:11:27.018Z", "dataType": 2},
    {"value": false, "id": 250, "timestamp": "2026-09-04T11:11:27.018Z", "dataType": 2},
    {"value": 1.23,  "id": 121, "timestamp": "2026-09-04T11:12:00.000Z", "dataType": 3},
    {"value": 3,     "id": 109, "timestamp": "2026-09-04T11:11:27.018Z", "dataType": 4},
    {"value": 16,    "id": 111, "timestamp": "2026-09-04T11:11:27.018Z", "dataType": 4},
    {"value": 20.0,  "id": 50,  "timestamp": "2026-09-04T11:11:27.018Z", "dataType": 4},
    {"value": 230.5, "id": 190, "timestamp": "2026-09-04T11:11:27.018Z", "dataType": 3},
    {"value": "{\"Id\":1}", "id": 129, "timestamp": "2026-09-04T11:30:00.000Z", "dataType": 6}
  ]
}`

func parseState(t *testing.T, body string) ChargerState {
	t.Helper()
	var res observationsResponse
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return newChargerState(&res)
}

func TestNewChargerStateDecodesObservations(t *testing.T) {
	s := parseState(t, observationsJson)

	if s.SmartCharging == nil || !*s.SmartCharging {
		t.Errorf("SmartCharging = %v, want true", s.SmartCharging)
	}
	if s.ConnectedToCloud == nil || *s.ConnectedToCloud {
		t.Errorf("ConnectedToCloud = %v, want false", s.ConnectedToCloud)
	}
	if s.SessionEnergy == nil || *s.SessionEnergy != 1.23 {
		t.Errorf("SessionEnergy = %v, want 1.23", s.SessionEnergy)
	}
	if s.ChargerOpMode == nil || *s.ChargerOpMode != 3 {
		t.Errorf("ChargerOpMode = %v, want 3", s.ChargerOpMode)
	}
	if s.InVoltageT1T2 == nil || *s.InVoltageT1T2 != 230.5 {
		t.Errorf("InVoltageT1T2 = %v, want 230.5", s.InVoltageT1T2)
	}
	// Documented as Integer but mapped to a float64 field.
	if s.DynamicCircuitCurrentP1 == nil || *s.DynamicCircuitCurrentP1 != 16 {
		t.Errorf("DynamicCircuitCurrentP1 = %v, want 16", s.DynamicCircuitCurrentP1)
	}
	// A whole number arriving as a double still fills an int field.
	if s.OfflineMaxCircuitCurrentP1 == nil || *s.OfflineMaxCircuitCurrentP1 != 20 {
		t.Errorf("OfflineMaxCircuitCurrentP1 = %v, want 20", s.OfflineMaxCircuitCurrentP1)
	}
	// Observations with no mapping must not populate anything.
	if s.CableLocked != nil {
		t.Errorf("CableLocked = %v, want nil", s.CableLocked)
	}
}

// LatestPulse stands in for the removed state endpoint's field, so it must
// track the newest mapped observation and ignore unmapped ones.
func TestNewChargerStateLatestPulse(t *testing.T) {
	s := parseState(t, observationsJson)
	want := time.Date(2026, 9, 4, 11, 12, 0, 0, time.UTC)
	if s.LatestPulse == nil {
		t.Fatal("LatestPulse = nil, want newest mapped observation timestamp")
	}
	if !s.LatestPulse.Equal(want) {
		t.Errorf("LatestPulse = %v, want %v", s.LatestPulse.UTC(), want)
	}
}

// One unexpected value must not blank out the rest of a charger's metrics.
func TestNewChargerStateSkipsUndecodableValue(t *testing.T) {
	s := parseState(t, `{"observations": [
		{"value": "not-a-number", "id": 120, "timestamp": "2026-09-04T11:11:27.018Z", "dataType": 6},
		{"value": 42.0, "id": 121, "timestamp": "2026-09-04T11:11:27.018Z", "dataType": 3}
	]}`)
	if s.TotalPower != nil {
		t.Errorf("TotalPower = %v, want nil", s.TotalPower)
	}
	if s.SessionEnergy == nil || *s.SessionEnergy != 42 {
		t.Errorf("SessionEnergy = %v, want 42", s.SessionEnergy)
	}
}

// Voltage is derived, not observed: the removed state endpoint's value equalled
// the highest input voltage measurement in every recorded sample.
func TestNewChargerStateDerivesVoltage(t *testing.T) {
	s := parseState(t, `{"observations": [
		{"value": 234.28199768066406, "id": 196, "timestamp": "2026-09-04T11:11:27.018Z", "dataType": 3},
		{"value": 246.22500610351562, "id": 194, "timestamp": "2026-09-04T11:11:27.018Z", "dataType": 3},
		{"value": 239.47000122070312, "id": 195, "timestamp": "2026-09-04T11:11:27.018Z", "dataType": 3}
	]}`)
	if s.Voltage == nil {
		t.Fatal("Voltage = nil, want the highest in voltage")
	}
	if *s.Voltage != 246.22500610351562 {
		t.Errorf("Voltage = %v, want 246.22500610351562", *s.Voltage)
	}
	// Must be a copy, not an alias of the in voltage field it came from.
	if s.Voltage == s.InVoltageT2T3 {
		t.Error("Voltage aliases InVoltageT2T3")
	}
}

func TestNewChargerStateVoltageAbsentWithoutInVoltage(t *testing.T) {
	s := parseState(t, `{"observations": [
		{"value": 1.23, "id": 121, "timestamp": "2026-09-04T11:11:27.018Z", "dataType": 3}
	]}`)
	if s.Voltage != nil {
		t.Errorf("Voltage = %v, want nil", *s.Voltage)
	}
}

func TestChargerStateObservationIds(t *testing.T) {
	if chargerStateObservationIds == "" {
		t.Fatal("no observation ids requested")
	}
	// Sorted, comma separated, and covering every mapped observation.
	want := "30,46,48,50,51,52,70,71,72,73,74,75,80,96,102,103,104,109,110,111,112,113,114,115,116,119,120,121,122,124,130,131,132,136,182,183,184,185,190,191,192,193,194,195,196,197,198,199,230,231,232,250"
	if chargerStateObservationIds != want {
		t.Errorf("ids = %q, want %q", chargerStateObservationIds, want)
	}
}
