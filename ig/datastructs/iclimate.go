package datastructs

// IClimate
type iClimateShadow struct {
	State StateIClimate `json:"state"`
}

// StateIClimate represents the State data structure from an IntelliClimate packet
type StateIClimate struct {
	Reported ReportedIClimate `json:"reported"`
}

// ReportedIClimate represents the Reported data structure from an IntelliClimate packet
type ReportedIClimate struct {
	Config    ConfigIClimate  `json:"config"`
	Metrics   MetricsIClimate `json:"metrics"`
	Status    StatusIClimate  `json:"status"`
	Source    string          `json:"source"`
	Device    string          `json:"device"`
	Timestamp int64           `json:"timestamp"`
	Connected bool            `json:"connected"`
}

// ConfigIClimate represents the Config data structure from an IntelliClimate packet
type ConfigIClimate struct {
	Units     UnitsIClimate     `json:"units"`
	Functions FunctionsIClimate `json:"functions"`
	Advanced  AdvancedIClimate  `json:"advanced"`
	General   GeneralIClimate   `json:"general"`
}

// MetricsIClimate represents the Metrics data structure from an IntelliClimate packet
type MetricsIClimate struct {
	AirTemp        float64 `json:"air_temp"`
	DayNight       string  `json:"day_night"`
	FailSafeAlarms bool    `json:"fail_safe_alarms"`
	Light          float64 `json:"light"`
	PowerFail      bool    `json:"power_fail"`
	Rh             float64 `json:"rh"`
	Vpd            float64 `json:"vpd"`
	Co2            float64 `json:"co2"`
	Intruder       bool    `json:"intruder_alarm"`
	OutsideTemp    float64 `json:"outside_temp_sensor"`
	EnviroAirTemp1 float64 `json:"enviro_air_temp_1"`
	EnviroAirTemp2 float64 `json:"enviro_air_temp_2"`
	EnviroRH1      float64 `json:"enviro_rh_1"`
	EnviroRH2      float64 `json:"enviro_rh_2"`
	EnviroCO21     float64 `json:"enviro_co2_1"`
	EnviroCO22     float64 `json:"enviro_co2_2"`
	EnviroLight1   float64 `json:"enviro_light_1"`
	EnviroLight2   float64 `json:"enviro_light_2"`
}

// StatusIClimate represents the Status data structure from an IntelliClimate packet
type StatusIClimate struct {
	Readings         ReadingsIClimate         `json:"readings"`
	Statistics       StatisticsIClimate       `json:"statistics"`
	ModeAlarmHistory ModeAlarmHistoryIClimate `json:"mode_alarm_history"`
	SetPoints        []SetPointIClimate       `json:"set_points"`
	Status           []StatusStatusIClimate   `json:"status"`
}

// ReadingsIClimate represents the Readings data structure from an IntelliClimate packet
type ReadingsIClimate struct {
	AirTemp        AirTempIClimate        `json:"air_temp"`
	Detent         byte                   `json:"detent"`
	FailSafeAlarms FailSafeAlarmsIClimate `json:"fail_safe_alarms"`
	Light          LightIClimate          `json:"light"`
	PowerFail      PowerFailIClimate      `json:"power_fail"`
	Rh             RhIClimate             `json:"rh"`
	CO2            CO2IClimate            `json:"co2"`
	IntruderAlarm  IntruderAlarmIClimate  `json:"intruder"`
}

// IntruderAlarmIClimate represents the IntruderAlarm data structure from an IntelliClimate packet
type IntruderAlarmIClimate struct {
	Enabled bool `json:"enabled"`
	Page    bool `json:"page"`
}

// CO2IClimate represents the CO2 data structure from an IntelliClimate packet
type CO2IClimate struct {
	Target  float64 `json:"target"`
	Enabled bool    `json:"enabled"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
	Page    bool    `json:"page"`
}

// AirTempIClimate represents the AirTemp data structure from an IntelliClimate packet
type AirTempIClimate struct {
	Cool    float64 `json:"cool"`
	Enabled bool    `json:"enabled"`
	Heat    float64 `json:"heat"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
	Page    bool    `json:"page"`
}

// FailSafeAlarmsIClimate represents the FailSafeAlarms data structure from an IntelliClimate packet
type FailSafeAlarmsIClimate struct {
	Enabled bool `json:"enabled"`
	Page    bool `json:"page"`
}

// LightIClimate represents the Light data structure from an IntelliClimate packet
type LightIClimate struct {
	Enabled bool    `json:"enabled"`
	Min     float64 `json:"min"`
	Page    bool    `json:"page"`
}

// PowerFailIClimate represents the PowerFail data structure from an IntelliClimate packet
type PowerFailIClimate struct {
	Enabled bool `json:"enabled"`
	Page    bool `json:"page"`
}

// RhIClimate represents the Rh data structure from an IntelliClimate packet
type RhIClimate struct {
	Enabled bool `json:"enabled"`
	Max     byte `json:"max"`
	Min     byte `json:"min"`
	Page    bool `json:"page"`
	Target  byte `json:"target"`
}

// StatisticsIClimate represents the Statistics data structure from an IntelliClimate packet
type StatisticsIClimate struct {
	Lights float64 `json:"lights"`
	CO2    float64 `json:"CO2"`
}

// ModeAlarmHistoryIClimate represents the ModeAlarmHistory data structure from an IntelliClimate packet
type ModeAlarmHistoryIClimate struct {
	Alarms []AlarmsIClimate `json:"alarms"`
	Mode   []ModeIClimate   `json:"mode"`
}

// AlarmsIClimate represents the Alarms data structure from an IntelliClimate packet
type AlarmsIClimate struct {
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
}

// ModeIClimate represents the Mode data structure from an IntelliClimate packet
type ModeIClimate struct {
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
}

// SetPointIClimate represents the SetPoint data structure from an IntelliClimate packet
type SetPointIClimate struct {
	LightBank     string  `json:"light_bank"`
	LightOn       int     `json:"light_on"`
	LightDuration int     `json:"light_duration"`
	DayTemp       float64 `json:"day_temp"`
	NightDropDeg  float64 `json:"night_drop_deg"`
	RhDay         int     `json:"rh_day"`
	RhMax         int     `json:"rh_max"`
	RhNight       int     `json:"rh_night"`
	CO2           int     `json:"co2"`
}

// StatusStatusIClimate represents the StatusStatus data structure from an IntelliClimate packet
type StatusStatusIClimate struct {
	Active    bool   `json:"active"`
	Enabled   bool   `json:"enabled"`
	ForceOn   bool   `json:"force_on"`
	Function  string `json:"function"`
	Installed bool   `json:"installed"`
}

// UnitsIClimate represents the Units data structure from an IntelliClimate packet
type UnitsIClimate struct {
	DateFormat  string `json:"date_format"`
	Temperature string `json:"temperature"`
}

// FunctionsIClimate represents the Functions data structure from an IntelliClimate packet
type FunctionsIClimate struct {
	Fan1                        bool   `json:"fan_1"`
	Fan2                        bool   `json:"fan_2"`
	AirConditioner              bool   `json:"air_conditioner"`
	Heater                      bool   `json:"heater"`
	Co2Sensor                   bool   `json:"co2_sensor"`
	Co2SensorRange              string `json:"co2_sensor_range"`
	Co2Injection                bool   `json:"co2_injection"`
	Co2Extraction               bool   `json:"co2_extraction"`
	Dehumidifier                bool   `json:"dehumidifier"`
	Humidifier                  bool   `json:"humidifier"`
	PulsedFogger                bool   `json:"pulsed_fogger"`
	LightBank1                  bool   `json:"light_bank_1"`
	LightsAirColored            bool   `json:"lights_air_colored"`
	LightBank2                  bool   `json:"light_bank_2"`
	LampOverTempShutdownSensors bool   `json:"lamp_over_temp_shutdown_sensors"`
	OutsideTempSensor           bool   `json:"outside_temp_sensor"`
	SecondEnviroSensor          bool   `json:"second_enviro_sensor"`
	IntruderAlarm               bool   `json:"intruder_alarm"`
	DehumidifyBy                string `json:"dehumidify_by"`
	Setup                       string `json:"setup"`
	MuteBuzzer                  bool   `json:"mute_buzzer"`
}

// AdvancedIClimate represents the Advanced data structure from an IntelliClimate packet
type AdvancedIClimate struct {
	ViewAdvancedSetting bool                     `json:"view_advanced_setting"`
	SwitchingOffsets    SwitchingOffsetsIClimate `json:"switching_offsets"`
	FailSafeSettings    FailSafeSettingsIClimate `json:"fail_safe_settings"`
	Rules               RulesIClimate            `json:"rules"`
}

// RulesIClimate represents the Rules data structure from an IntelliClimate packet
type RulesIClimate struct {
	HumidifyTempRules     HumidifyTempRulesIClimate     `json:"humidify_temp_rules"`
	MinimumAirChangeRules MinimumAirChangeRulesIClimate `json:"minimum_air_change_rules"`
	AllowAirCon           bool                          `json:"allow_air_con"`
	SetpointRamping       SetpointRampingIClimate       `json:"setpoint_ramping"`
	AirCon                AirConIClimate                `json:"air_con"`
	CO2Rules              CO2RulesIClimate              `json:"co2_rules"`
	Humidification        HumidificationIClimate        `json:"humidification"`
	Lighting              LightingIClimate              `json:"lighting"`
	FoggingRules          FoggingRulesIClimate          `json:"fogging_rules"`
	PurgingRules          PurgingRulesIClimate          `json:"purging_rules"`
}

// HumidifyTempRulesIClimate represents the HumidifyTempRules data structure from an IntelliClimate packet
type HumidifyTempRulesIClimate struct {
	LowerCoolingTemp float64 `json:"lower_cooling_temp"`
	RaiseHeatingTemp float64 `json:"raise_heating_temp"`
	RhLowThenRaise   float64 `json:"rh_low_then_raise"`
	PreventHeater    byte    `json:"prevent_heater"`
	HeatingOffset    float64 `json:"heating_offset"`
}

// MinimumAirChangeRulesIClimate represents the MinimumAirChangeRules data structure from an IntelliClimate packet
type MinimumAirChangeRulesIClimate struct {
	DaySecs        int `json:"day_secs"`
	EveryDayMins   int `json:"every_day_mins"`
	NightSecs      int `json:"night_secs"`
	EveryNightMins int `json:"every_night_mins"`
}

// SetpointRampingIClimate represents the SetpointRamping data structure from an IntelliClimate packet
type SetpointRampingIClimate struct {
	RampSetpoints byte `json:"ramp_setpoints"`
}

// AirConIClimate represents the AirCon data structure from an IntelliClimate packet
type AirConIClimate struct {
	ForceAirCon      bool    `json:"force_air_con"`
	AutoChangeAirCon float64 `json:"auto_change_air_con"`
	StartBefore      byte    `json:"start_before"`
	AutoStartAirCon  float64 `json:"auto_start_air_con"`
}

// CO2RulesIClimate represents the CO2Rules data structure from an IntelliClimate packet
type CO2RulesIClimate struct {
	Co2InjectionAllowed  bool    `json:"co2_injection_allowed"`
	InjectIfLightGreater float64 `json:"inject_if_light_greater"`
	Co2InjectionAvoid    bool    `json:"co2_injection_avoid"`
	Co2Cycling           float64 `json:"co2_cycling"`
	RiseVentTemp         float64 `json:"rise_vent_temp"`
	InjectTimeMin        byte    `json:"inject_time_min"`
	InjectTimeMax        byte    `json:"inject_time_max"`
	WaitTimeMin          byte    `json:"wait_time_min"`
	WaitTimeMax          byte    `json:"wait_time_max"`
	VentTimeMin          byte    `json:"vent_time_min"`
	VentTimeMax          byte    `json:"vent_time_max"`
}

// HumidificationIClimate represents the Humidification data structure from an IntelliClimate packet
type HumidificationIClimate struct {
	AllowHumidification  bool `json:"allow_humidification"`
	ChangeHumidification byte `json:"change_humidification"`
}

// LightingIClimate represents the Lighting data structure from an IntelliClimate packet
type LightingIClimate struct {
	LampCoolDownTime  byte `json:"lamp_cool_down_time"`
	SwOnNextLightBank byte `json:"sw_on_next_light_bank"`
}

// FoggingRulesIClimate represents the FoggingRules data structure from an IntelliClimate packet
type FoggingRulesIClimate struct {
	FogToCool      byte    `json:"fog_to_cool"`
	FogToAchieveRh float64 `json:"fog_to_achieve_rh"`
	FogTimes       int     `json:"fog_times"`
	FogTimeMax     byte    `json:"fog_time_max"`
	FogTimeMin     byte    `json:"fog_time_min"`
}

// PurgingRulesIClimate represents the PurgingRules data structure from an IntelliClimate packet
type PurgingRulesIClimate struct {
	PurgeMins byte `json:"purge_mins"`
	PurgeMin  byte `json:"purge_min"`
	PurgeMax  byte `json:"purge_max"`
}

// SwitchingOffsetsIClimate represents the SwitchingOffsets data structure from an IntelliClimate packet
type SwitchingOffsetsIClimate struct {
	AirConditionerOn  float64 `json:"air_conditioner_on"`
	AirConditionerOff float64 `json:"air_conditioner_off"`
	CO2On             float64 `json:"co2_on"`
	CO2Off            float64 `json:"co2_off"`
	DehumidifierOn    float64 `json:"dehumidifier_on"`
	DehumidifierOff   float64 `json:"dehumidifier_off"`
	FansOn            float64 `json:"fans_on"`
	FansOff           float64 `json:"fans_off"`
	HeaterOn          float64 `json:"heater_on"`
	HeaterOff         float64 `json:"heater_off"`
	HumidifierOn      float64 `json:"humidifier_on"`
	HumidifierOff     float64 `json:"humidifier_off"`
	PulsedFoggerOn    float64 `json:"pulsed_fogger_on"`
	PulsedFoggerOff   float64 `json:"pulsed_fogger_off"`
}

// FailSafeSettingsIClimate represents the FailSafeSettings data structure from an IntelliClimate packet
type FailSafeSettingsIClimate struct {
	FanFailOverride      FanFailOverrideIClimate      `json:"fan_fail_override"`
	AirConOverride       AirConOverrideIClimate       `json:"air_con_override"`
	DehumidifierOverride DehumidifierOverrideIClimate `json:"dehumidifier_override"`
	Co2FailSafe          Co2FailSafeIClimate          `json:"co2_fail_safe"`
	Co2InjectionOverride Co2InjectionOverrideIClimate `json:"co2_injection_override"`
	PowerFailure         PowerFailureIClimate         `json:"power_failure"`
	LightingOverride     LightingOverrideIClimate     `json:"light_falls_alarm_minimum"`
}

// FanFailOverrideIClimate represents the FanFailOverride data structure from an IntelliClimate packet
type FanFailOverrideIClimate struct {
	SwOffLightTempExceed  float64 `json:"sw_off_light_temp_exceed"`
	SwOffLightsTempExceed float64 `json:"sw_off_lights_temp_exceed"`
}

// AirConOverrideIClimate represents the AirConOverride data structure from an IntelliClimate packet
type AirConOverrideIClimate struct {
	SwAllExhaustFans float64 `json:"sw_all_exhaust_fans"`
}

// DehumidifierOverrideIClimate represents the DehumidifierOverride data structure from an IntelliClimate packet
type DehumidifierOverrideIClimate struct {
	SwOnFansRhExceed byte `json:"sw_on_fans_rh_exceed"`
	SwAcRhExceed     byte `json:"sw_ac_rh_exceed"`
}

// Co2FailSafeIClimate represents the Co2FailSafe data structure from an IntelliClimate packet
type Co2FailSafeIClimate struct {
	SwOnFansCo2Exceed int `json:"sw_on_fans_co2_exceed"`
}

// Co2InjectionOverrideIClimate represents the Co2InjectionOverride data structure from an IntelliClimate packet
type Co2InjectionOverrideIClimate struct {
	RevertFansCo2Falls int `json:"revert_fans_co2_falls"`
}

// PowerFailureIClimate represents the PowerFailure data structure from an IntelliClimate packet
type PowerFailureIClimate struct {
	SwLightsAfterCoolDown byte `json:"sw_lights_after_cool_down"`
}

// LightingOverrideIClimate represents the LightingOverride data structure from an IntelliClimate packet
type LightingOverrideIClimate struct {
	LightFallsAlarmMinimum bool `json:"light_falls_alarm_minimum"`
}

// GeneralIClimate represents the General data structure from an IntelliClimate packet
type GeneralIClimate struct {
	DeviceName string  `json:"device_name"`
	Firmware   float64 `json:"firmware"`
}

// ClimateHistory - consists of a slice of history points
type ClimateHistory struct {
	Points []*ClimateHistoryPoint `json:"points"`
}

// ClimateMetricsHistory - Metrics the history point contains
type ClimateMetricsHistory struct {
	AirTemp float64 `json:"air_temp"`
	Rh      float64 `json:"rh"`
	Vpd     float64 `json:"vpd"`
	CO2     float64 `json:"co2"`
	Light   float64 `json:"light"`
}

// ClimateHistoryPoint - defines a single history point reported for a IntelliClimate
type ClimateHistoryPoint struct {
	Timestamp float64               `json:"timestamp"`
	Status    Status                `json:"status"`
	Metrics   ClimateMetricsHistory `json:"metrics"`
}
