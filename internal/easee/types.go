package easee

import "time"

type Charger struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	ProductCode int    `json:"productCode"`
}

type ChargerState struct {
	SmartCharging                                *bool      `json:"smartCharging"`
	CableLocked                                  *bool      `json:"cableLocked"`
	ChargerOpMode                                *int       `json:"chargerOpMode"`
	TotalPower                                   *float64   `json:"totalPower"`
	SessionEnergy                                *float64   `json:"sessionEnergy"`
	EnergyPerHour                                *float64   `json:"energyPerHour"`
	WifiRSSI                                     *int       `json:"wiFiRSSI"`
	CellRSSI                                     *int       `json:"cellRSSI"`
	LocalRSSI                                    *int       `json:"localRSSI"`
	OutputPhase                                  *int       `json:"outputPhase"`
	DynamicCircuitCurrentP1                      *float64   `json:"dynamicCircuitCurrentP1"`
	DynamicCircuitCurrentP2                      *float64   `json:"dynamicCircuitCurrentP2"`
	DynamicCircuitCurrentP3                      *float64   `json:"dynamicCircuitCurrentP3"`
	LatestPulse                                  *time.Time `json:"latestPulse"`
	ChargerFirmware                              *int       `json:"chargerFirmware"`
	Voltage                                      *float64   `json:"voltage"`
	ChargerRAT                                   *int       `json:"chargerRAT"`
	LockCablePermanently                         *bool      `json:"lockCablePermanently"`
	InCurrentT2                                  *float64   `json:"inCurrentT2"`
	InCurrentT3                                  *float64   `json:"inCurrentT3"`
	InCurrentT4                                  *float64   `json:"inCurrentT4"`
	InCurrentT5                                  *float64   `json:"inCurrentT5"`
	OutputCurrent                                *float64   `json:"outputCurrent"`
	IsOnline                                     *bool      `json:"isOnline"`
	InVoltageT1T2                                *float64   `json:"inVoltageT1T2"`
	InVoltageT1T3                                *float64   `json:"inVoltageT1T3"`
	InVoltageT1T4                                *float64   `json:"inVoltageT1T4"`
	InVoltageT1T5                                *float64   `json:"inVoltageT1T5"`
	InVoltageT2T3                                *float64   `json:"inVoltageT2T3"`
	InVoltageT2T4                                *float64   `json:"inVoltageT2T4"`
	InVoltageT2T5                                *float64   `json:"inVoltageT2T5"`
	InVoltageT3T4                                *float64   `json:"inVoltageT3T4"`
	InVoltageT3T5                                *float64   `json:"inVoltageT3T5"`
	InVoltageT4T5                                *float64   `json:"inVoltageT4T5"`
	LedMode                                      *int       `json:"ledMode"`
	CableRating                                  *float64   `json:"cableRating"`
	DynamicChargerCurrent                        *float64   `json:"dynamicChargerCurrent"`
	CircuitTotalAllocatedPhaseConductorCurrentL1 *float64   `json:"circuitTotalAllocatedPhaseConductorCurrentL1"`
	CircuitTotalAllocatedPhaseConductorCurrentL2 *float64   `json:"circuitTotalAllocatedPhaseConductorCurrentL2"`
	CircuitTotalAllocatedPhaseConductorCurrentL3 *float64   `json:"circuitTotalAllocatedPhaseConductorCurrentL3"`
	CircuitTotalPhaseConductorCurrentL1          *float64   `json:"circuitTotalPhaseConductorCurrentL1"`
	CircuitTotalPhaseConductorCurrentL2          *float64   `json:"circuitTotalPhaseConductorCurrentL2"`
	CircuitTotalPhaseConductorCurrentL3          *float64   `json:"circuitTotalPhaseConductorCurrentL3"`
	ReasonForNoCurrent                           *int       `json:"reasonForNoCurrent"`
	WifiAPEnabled                                *bool      `json:"wiFiAPEnabled"`
	LifeTimeEnergy                               *float64   `json:"lifetimeEnergy"`
	OfflineMaxCircuitCurrentP1                   *int       `json:"offlineMaxCircuitCurrentP1"`
	OfflineMaxCircuitCurrentP2                   *int       `json:"offlineMaxCircuitCurrentP2"`
	OfflineMaxCircuitCurrentP3                   *int       `json:"offlineMaxCircuitCurrentP3"`
	ErrorCode                                    *int       `json:"errorCode"`
	FaultErrorCode                               *int       `json:"faultErrorCode"`
	EqAvailableCurrentP1                         *float64   `json:"eqAvailableCurrentP1"`
	EqAvailableCurrentP2                         *float64   `json:"eqAvailableCurrentP2"`
	EqAvailableCurrentP3                         *float64   `json:"eqAvailableCurrentP3"`
	DeratedCurrent                               *float64   `json:"deratedCurrent"`
	DeratingActive                               *bool      `json:"deratingActive"`
	ConnectedToCloud                             *bool      `json:"connectedToCloud"`
}
