package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/terjesannum/easee-exporter/internal/easee"
)

type ChargerStateCollector struct {
	chargerState               *easee.ChargerState
	smartCharging              *prometheus.Desc
	cableLocked                *prometheus.Desc
	opMode                     *prometheus.Desc
	totalPower                 *prometheus.Desc
	sessionEnergy              *prometheus.Desc
	energyPerHour              *prometheus.Desc
	rssi                       *prometheus.Desc
	outputPhase                *prometheus.Desc
	dynamicCircuitCurrent      *prometheus.Desc
	chargerFirmware            *prometheus.Desc
	chargerRAT                 *prometheus.Desc
	cableLockPermanent         *prometheus.Desc
	inCurrent                  *prometheus.Desc
	outputCurrent              *prometheus.Desc
	inVoltage                  *prometheus.Desc
	ledMode                    *prometheus.Desc
	cableRating                *prometheus.Desc
	dynamicChargerCurrent      *prometheus.Desc
	totalAllocatedPhaseCurrent *prometheus.Desc
	totalPhaseCurrent          *prometheus.Desc
	noCurrentReason            *prometheus.Desc
	wifiAPEnabled              *prometheus.Desc
	lifetimeEnergy             *prometheus.Desc
	offlineMaxCircuitCurrent   *prometheus.Desc
	errorCode                  *prometheus.Desc
	faultErrorCode             *prometheus.Desc
	eqAvailableCurrent         *prometheus.Desc
	deratedCurrent             *prometheus.Desc
	deratingActive             *prometheus.Desc
	connectedToCloud           *prometheus.Desc
	isOnline                   *prometheus.Desc
	voltage                    *prometheus.Desc
	latestPulse                *prometheus.Desc
}

func NewChargerStateCollector(charger string) *ChargerStateCollector {
	return &ChargerStateCollector{
		smartCharging: prometheus.NewDesc(
			"easee_charger_smart_charging",
			"Smart charging state",
			nil,
			prometheus.Labels{"charger": charger},
		),
		cableLocked: prometheus.NewDesc(
			"easee_charger_cable_locked",
			"Cable lock",
			nil,
			prometheus.Labels{"charger": charger},
		),
		opMode: prometheus.NewDesc(
			"easee_charger_operation_mode",
			"Operation mode",
			nil,
			prometheus.Labels{"charger": charger},
		),
		totalPower: prometheus.NewDesc(
			"easee_charger_total_power",
			"Total power",
			nil,
			prometheus.Labels{"charger": charger},
		),
		sessionEnergy: prometheus.NewDesc(
			"easee_charger_session_energy",
			"Session energy",
			nil,
			prometheus.Labels{"charger": charger},
		),
		energyPerHour: prometheus.NewDesc(
			"easee_charger_energy_per_hour",
			"Energy per hour",
			nil,
			prometheus.Labels{"charger": charger},
		),
		rssi: prometheus.NewDesc(
			"easee_charger_signal_strength",
			"Signal strength",
			[]string{"type"},
			prometheus.Labels{"charger": charger},
		),
		outputPhase: prometheus.NewDesc(
			"easee_charger_output_phase",
			"Output phase",
			nil,
			prometheus.Labels{"charger": charger},
		),
		dynamicCircuitCurrent: prometheus.NewDesc(
			"easee_charger_dynamic_circuit_current",
			"Dynamic circuit current",
			[]string{"p"},
			prometheus.Labels{"charger": charger},
		),
		chargerFirmware: prometheus.NewDesc(
			"easee_charger_firmware_version",
			"Charger firmware version",
			nil,
			prometheus.Labels{"charger": charger},
		),
		chargerRAT: prometheus.NewDesc(
			"easee_charger_radio_access_technology",
			"Charger radio access technology",
			nil,
			prometheus.Labels{"charger": charger},
		),
		cableLockPermanent: prometheus.NewDesc(
			"easee_charger_cable_permanent_lock",
			"Permanent cable lock",
			nil,
			prometheus.Labels{"charger": charger},
		),
		inCurrent: prometheus.NewDesc(
			"easee_charger_in_current",
			"In current",
			[]string{"t"},
			prometheus.Labels{"charger": charger},
		),
		outputCurrent: prometheus.NewDesc(
			"easee_charger_output_current",
			"Output current",
			nil,
			prometheus.Labels{"charger": charger},
		),
		inVoltage: prometheus.NewDesc(
			"easee_charger_in_voltage",
			"In voltage",
			[]string{"ta", "tb"},
			prometheus.Labels{"charger": charger},
		),
		ledMode: prometheus.NewDesc(
			"easee_charger_led_mode",
			"LED mode",
			nil,
			prometheus.Labels{"charger": charger},
		),
		cableRating: prometheus.NewDesc(
			"easee_charger_cable_rating",
			"Cable rating",
			nil,
			prometheus.Labels{"charger": charger},
		),
		dynamicChargerCurrent: prometheus.NewDesc(
			"easee_charger_dynamic_current",
			"Dynamic charger current",
			nil,
			prometheus.Labels{"charger": charger},
		),
		totalAllocatedPhaseCurrent: prometheus.NewDesc(
			"easee_charger_total_allocated_phase_current",
			"Total allocated phase current",
			[]string{"l"},
			prometheus.Labels{"charger": charger},
		),
		totalPhaseCurrent: prometheus.NewDesc(
			"easee_charger_total_phase_current",
			"Total phase current",
			[]string{"l"},
			prometheus.Labels{"charger": charger},
		),
		noCurrentReason: prometheus.NewDesc(
			"easee_charger_no_current_reason",
			"Reason for no current",
			nil,
			prometheus.Labels{"charger": charger},
		),
		wifiAPEnabled: prometheus.NewDesc(
			"easee_charger_wifi_access_point_enabled",
			"Access point enabled",
			nil,
			prometheus.Labels{"charger": charger},
		),
		lifetimeEnergy: prometheus.NewDesc(
			"easee_charger_lifetime_energy",
			"Lifetime energy",
			nil,
			prometheus.Labels{"charger": charger},
		),
		offlineMaxCircuitCurrent: prometheus.NewDesc(
			"easee_charger_max_offline_circuit_current",
			"Max offline circuit current",
			[]string{"p"},
			prometheus.Labels{"charger": charger},
		),
		errorCode: prometheus.NewDesc(
			"easee_charger_error_code",
			"Error code",
			nil,
			prometheus.Labels{"charger": charger},
		),
		faultErrorCode: prometheus.NewDesc(
			"easee_charger_fault_error_code",
			"Fault error code",
			nil,
			prometheus.Labels{"charger": charger},
		),
		eqAvailableCurrent: prometheus.NewDesc(
			"easee_charger_eq_available_current",
			"Eq available current",
			[]string{"p"},
			prometheus.Labels{"charger": charger},
		),
		deratedCurrent: prometheus.NewDesc(
			"easee_charger_derated_current",
			"Derated current",
			nil,
			prometheus.Labels{"charger": charger},
		),
		deratingActive: prometheus.NewDesc(
			"easee_charger_derating_active",
			"Derating active",
			nil,
			prometheus.Labels{"charger": charger},
		),
		connectedToCloud: prometheus.NewDesc(
			"easee_charger_cloud_connection",
			"Cloud connection",
			nil,
			prometheus.Labels{"charger": charger},
		),
		isOnline: prometheus.NewDesc(
			"easee_charger_online",
			"Charger online",
			nil,
			prometheus.Labels{"charger": charger},
		),
		voltage: prometheus.NewDesc(
			"easee_charger_voltage",
			"Charger voltage",
			nil,
			prometheus.Labels{"charger": charger},
		),
		latestPulse: prometheus.NewDesc(
			"easee_charger_latest_pulse",
			"Latest pulse",
			nil,
			prometheus.Labels{"charger": charger},
		),
	}
}

func (c *ChargerStateCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.smartCharging
	ch <- c.cableLocked
	ch <- c.opMode
	ch <- c.totalPower
	ch <- c.sessionEnergy
	ch <- c.energyPerHour
	ch <- c.rssi
	ch <- c.outputPhase
	ch <- c.dynamicCircuitCurrent
	ch <- c.chargerFirmware
	ch <- c.chargerRAT
	ch <- c.cableLockPermanent
	ch <- c.inCurrent
	ch <- c.outputCurrent
	ch <- c.inVoltage
	ch <- c.ledMode
	ch <- c.cableRating
	ch <- c.dynamicChargerCurrent
	ch <- c.totalAllocatedPhaseCurrent
	ch <- c.totalPhaseCurrent
	ch <- c.noCurrentReason
	ch <- c.wifiAPEnabled
	ch <- c.lifetimeEnergy
	ch <- c.offlineMaxCircuitCurrent
	ch <- c.errorCode
	ch <- c.faultErrorCode
	ch <- c.eqAvailableCurrent
	ch <- c.deratedCurrent
	ch <- c.deratingActive
	ch <- c.connectedToCloud
	ch <- c.isOnline
	ch <- c.voltage
	ch <- c.latestPulse
}

func (c *ChargerStateCollector) Collect(ch chan<- prometheus.Metric) {
	b2f := map[bool]float64{true: 1}
	if c.chargerState != nil {
		if c.chargerState.SmartCharging != nil {
			ch <- prometheus.MustNewConstMetric(
				c.smartCharging,
				prometheus.GaugeValue,
				b2f[*c.chargerState.SmartCharging],
			)
		}
		if c.chargerState.CableLocked != nil {
			ch <- prometheus.MustNewConstMetric(
				c.cableLocked,
				prometheus.GaugeValue,
				b2f[*c.chargerState.CableLocked],
			)
		}
		if c.chargerState.ChargerOpMode != nil {
			ch <- prometheus.MustNewConstMetric(
				c.opMode,
				prometheus.GaugeValue,
				float64(*c.chargerState.ChargerOpMode),
			)
		}
		if c.chargerState.TotalPower != nil {
			ch <- prometheus.MustNewConstMetric(
				c.totalPower,
				prometheus.GaugeValue,
				*c.chargerState.TotalPower,
			)
		}
		if c.chargerState.SessionEnergy != nil {
			ch <- prometheus.MustNewConstMetric(
				c.sessionEnergy,
				prometheus.GaugeValue,
				*c.chargerState.SessionEnergy,
			)
		}
		if c.chargerState.EnergyPerHour != nil {
			ch <- prometheus.MustNewConstMetric(
				c.energyPerHour,
				prometheus.GaugeValue,
				*c.chargerState.EnergyPerHour,
			)
		}
		if c.chargerState.WifiRSSI != nil {
			ch <- prometheus.MustNewConstMetric(
				c.rssi,
				prometheus.GaugeValue,
				float64(*c.chargerState.WifiRSSI),
				"wifi",
			)
		}
		if c.chargerState.CellRSSI != nil {
			ch <- prometheus.MustNewConstMetric(
				c.rssi,
				prometheus.GaugeValue,
				float64(*c.chargerState.CellRSSI),
				"cell",
			)
		}
		if c.chargerState.LocalRSSI != nil {
			ch <- prometheus.MustNewConstMetric(
				c.rssi,
				prometheus.GaugeValue,
				float64(*c.chargerState.LocalRSSI),
				"local",
			)
		}
		if c.chargerState.OutputPhase != nil {
			ch <- prometheus.MustNewConstMetric(
				c.outputPhase,
				prometheus.GaugeValue,
				float64(*c.chargerState.OutputPhase),
			)
		}
		if c.chargerState.DynamicCircuitCurrentP1 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.dynamicCircuitCurrent,
				prometheus.GaugeValue,
				*c.chargerState.DynamicCircuitCurrentP1,
				"1",
			)
		}
		if c.chargerState.DynamicCircuitCurrentP2 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.dynamicCircuitCurrent,
				prometheus.GaugeValue,
				*c.chargerState.DynamicCircuitCurrentP2,
				"2",
			)
		}
		if c.chargerState.DynamicCircuitCurrentP3 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.dynamicCircuitCurrent,
				prometheus.GaugeValue,
				*c.chargerState.DynamicCircuitCurrentP3,
				"3",
			)
		}
		if c.chargerState.ChargerFirmware != nil {
			ch <- prometheus.MustNewConstMetric(
				c.chargerFirmware,
				prometheus.GaugeValue,
				float64(*c.chargerState.ChargerFirmware),
			)
		}
		if c.chargerState.ChargerRAT != nil {
			ch <- prometheus.MustNewConstMetric(
				c.chargerRAT,
				prometheus.GaugeValue,
				float64(*c.chargerState.ChargerRAT),
			)
		}
		if c.chargerState.LockCablePermanently != nil {
			ch <- prometheus.MustNewConstMetric(
				c.cableLockPermanent,
				prometheus.GaugeValue,
				b2f[*c.chargerState.LockCablePermanently],
			)
		}
		if c.chargerState.InCurrentT2 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.inCurrent,
				prometheus.GaugeValue,
				*c.chargerState.InCurrentT2,
				"2",
			)
		}
		if c.chargerState.InCurrentT3 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.inCurrent,
				prometheus.GaugeValue,
				*c.chargerState.InCurrentT3,
				"3",
			)
		}
		if c.chargerState.InCurrentT4 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.inCurrent,
				prometheus.GaugeValue,
				*c.chargerState.InCurrentT4,
				"4",
			)
		}
		if c.chargerState.InCurrentT5 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.inCurrent,
				prometheus.GaugeValue,
				*c.chargerState.InCurrentT5,
				"5",
			)
		}
		if c.chargerState.OutputCurrent != nil {
			ch <- prometheus.MustNewConstMetric(
				c.outputCurrent,
				prometheus.GaugeValue,
				*c.chargerState.OutputCurrent,
			)
		}
		if c.chargerState.InVoltageT1T2 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.inVoltage,
				prometheus.GaugeValue,
				*c.chargerState.InVoltageT1T2,
				"1",
				"2",
			)
		}
		if c.chargerState.InVoltageT1T3 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.inVoltage,
				prometheus.GaugeValue,
				*c.chargerState.InVoltageT1T3,
				"1",
				"3",
			)
		}
		if c.chargerState.InVoltageT1T4 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.inVoltage,
				prometheus.GaugeValue,
				*c.chargerState.InVoltageT1T4,
				"1",
				"4",
			)
		}
		if c.chargerState.InVoltageT1T5 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.inVoltage,
				prometheus.GaugeValue,
				*c.chargerState.InVoltageT1T5,
				"1",
				"5",
			)
		}
		if c.chargerState.InVoltageT2T3 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.inVoltage,
				prometheus.GaugeValue,
				*c.chargerState.InVoltageT2T3,
				"2",
				"3",
			)
		}
		if c.chargerState.InVoltageT2T4 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.inVoltage,
				prometheus.GaugeValue,
				*c.chargerState.InVoltageT2T4,
				"2",
				"4",
			)
		}
		if c.chargerState.InVoltageT2T5 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.inVoltage,
				prometheus.GaugeValue,
				*c.chargerState.InVoltageT2T5,
				"2",
				"5",
			)
		}
		if c.chargerState.InVoltageT3T4 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.inVoltage,
				prometheus.GaugeValue,
				*c.chargerState.InVoltageT3T4,
				"3",
				"4",
			)
		}
		if c.chargerState.InVoltageT3T5 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.inVoltage,
				prometheus.GaugeValue,
				*c.chargerState.InVoltageT3T5,
				"3",
				"5",
			)
		}
		if c.chargerState.InVoltageT4T5 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.inVoltage,
				prometheus.GaugeValue,
				*c.chargerState.InVoltageT4T5,
				"4",
				"5",
			)
		}
		if c.chargerState.LedMode != nil {
			ch <- prometheus.MustNewConstMetric(
				c.ledMode,
				prometheus.GaugeValue,
				float64(*c.chargerState.LedMode),
			)
		}
		if c.chargerState.CableRating != nil {
			ch <- prometheus.MustNewConstMetric(
				c.cableRating,
				prometheus.GaugeValue,
				*c.chargerState.CableRating,
			)
		}
		if c.chargerState.DynamicChargerCurrent != nil {
			ch <- prometheus.MustNewConstMetric(
				c.dynamicChargerCurrent,
				prometheus.GaugeValue,
				*c.chargerState.DynamicChargerCurrent,
			)
		}
		if c.chargerState.CircuitTotalAllocatedPhaseConductorCurrentL1 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.totalAllocatedPhaseCurrent,
				prometheus.GaugeValue,
				*c.chargerState.CircuitTotalAllocatedPhaseConductorCurrentL1,
				"1",
			)
		}
		if c.chargerState.CircuitTotalAllocatedPhaseConductorCurrentL2 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.totalAllocatedPhaseCurrent,
				prometheus.GaugeValue,
				*c.chargerState.CircuitTotalAllocatedPhaseConductorCurrentL2,
				"2",
			)
		}
		if c.chargerState.CircuitTotalAllocatedPhaseConductorCurrentL3 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.totalAllocatedPhaseCurrent,
				prometheus.GaugeValue,
				*c.chargerState.CircuitTotalAllocatedPhaseConductorCurrentL3,
				"3",
			)
		}
		if c.chargerState.CircuitTotalPhaseConductorCurrentL1 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.totalPhaseCurrent,
				prometheus.GaugeValue,
				*c.chargerState.CircuitTotalPhaseConductorCurrentL1,
				"1",
			)
		}
		if c.chargerState.CircuitTotalPhaseConductorCurrentL2 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.totalPhaseCurrent,
				prometheus.GaugeValue,
				*c.chargerState.CircuitTotalPhaseConductorCurrentL2,
				"2",
			)
		}
		if c.chargerState.CircuitTotalPhaseConductorCurrentL3 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.totalPhaseCurrent,
				prometheus.GaugeValue,
				*c.chargerState.CircuitTotalPhaseConductorCurrentL3,
				"3",
			)
		}
		if c.chargerState.ReasonForNoCurrent != nil {
			ch <- prometheus.MustNewConstMetric(
				c.noCurrentReason,
				prometheus.GaugeValue,
				float64(*c.chargerState.ReasonForNoCurrent),
			)
		}
		if c.chargerState.WifiAPEnabled != nil {
			ch <- prometheus.MustNewConstMetric(
				c.wifiAPEnabled,
				prometheus.GaugeValue,
				b2f[*c.chargerState.WifiAPEnabled],
			)
		}
		if c.chargerState.LifeTimeEnergy != nil {
			ch <- prometheus.MustNewConstMetric(
				c.lifetimeEnergy,
				prometheus.CounterValue,
				*c.chargerState.LifeTimeEnergy,
			)
		}
		if c.chargerState.OfflineMaxCircuitCurrentP1 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.offlineMaxCircuitCurrent,
				prometheus.GaugeValue,
				float64(*c.chargerState.OfflineMaxCircuitCurrentP1),
				"1",
			)
		}
		if c.chargerState.OfflineMaxCircuitCurrentP2 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.offlineMaxCircuitCurrent,
				prometheus.GaugeValue,
				float64(*c.chargerState.OfflineMaxCircuitCurrentP2),
				"2",
			)
		}
		if c.chargerState.OfflineMaxCircuitCurrentP3 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.offlineMaxCircuitCurrent,
				prometheus.GaugeValue,
				float64(*c.chargerState.OfflineMaxCircuitCurrentP3),
				"3",
			)
		}
		if c.chargerState.ErrorCode != nil {
			ch <- prometheus.MustNewConstMetric(
				c.errorCode,
				prometheus.GaugeValue,
				float64(*c.chargerState.ErrorCode),
			)
		}
		if c.chargerState.FaultErrorCode != nil {
			ch <- prometheus.MustNewConstMetric(
				c.faultErrorCode,
				prometheus.GaugeValue,
				float64(*c.chargerState.FaultErrorCode),
			)
		}
		if c.chargerState.EqAvailableCurrentP1 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.eqAvailableCurrent,
				prometheus.GaugeValue,
				*c.chargerState.EqAvailableCurrentP1,
				"1",
			)
		}
		if c.chargerState.EqAvailableCurrentP2 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.eqAvailableCurrent,
				prometheus.GaugeValue,
				*c.chargerState.EqAvailableCurrentP2,
				"2",
			)
		}
		if c.chargerState.EqAvailableCurrentP3 != nil {
			ch <- prometheus.MustNewConstMetric(
				c.eqAvailableCurrent,
				prometheus.GaugeValue,
				*c.chargerState.EqAvailableCurrentP3,
				"3",
			)
		}
		if c.chargerState.DeratedCurrent != nil {
			ch <- prometheus.MustNewConstMetric(
				c.deratedCurrent,
				prometheus.GaugeValue,
				*c.chargerState.DeratedCurrent,
			)
		}
		if c.chargerState.DeratingActive != nil {
			ch <- prometheus.MustNewConstMetric(
				c.deratingActive,
				prometheus.GaugeValue,
				b2f[*c.chargerState.DeratingActive],
			)
		}
		if c.chargerState.ConnectedToCloud != nil {
			ch <- prometheus.MustNewConstMetric(
				c.connectedToCloud,
				prometheus.GaugeValue,
				b2f[*c.chargerState.ConnectedToCloud],
			)
		}
		if c.chargerState.IsOnline != nil {
			ch <- prometheus.MustNewConstMetric(
				c.isOnline,
				prometheus.GaugeValue,
				b2f[*c.chargerState.IsOnline],
			)
		}
		if c.chargerState.Voltage != nil {
			ch <- prometheus.MustNewConstMetric(
				c.voltage,
				prometheus.GaugeValue,
				*c.chargerState.Voltage,
			)
		}
		if c.chargerState.LatestPulse != nil {
			ch <- prometheus.MustNewConstMetric(
				c.latestPulse,
				prometheus.GaugeValue,
				float64((*c.chargerState.LatestPulse).Unix()),
			)
		}
	}
}

func (c *ChargerStateCollector) UpdateState(state *easee.ChargerState) {
	c.chargerState = state
}
