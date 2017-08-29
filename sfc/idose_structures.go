package sfc

type iDoseShadow struct {
	State StateIDose `json:"state"`
}

// StateIDose represents the state of the IntelliDose
type StateIDose struct {
	Reported ReportedIDose `json:"reported"`
}

// ReportedIDose represents the top level report from the IntelliDose
type ReportedIDose struct {
	Config    ConfigIDose   `json:"config"`
	Metrics   MetricsIDose  `json:"metrics"`
	Settings  SettingsIDose `json:"status"`
	Source    string        `json:"source"`
	Device    string        `json:"device"`
	Timestamp int64         `json:"timestamp"`
}

// ConfigIDose represents the IntelliDose config
type ConfigIDose struct {
	Units     UnitsIDose     `json:"units"`
	Times     TimesIDose     `json:"times"`
	Functions FunctionsIDose `json:"functions"`
	Advanced  AdvancedIDose  `json:"advanced"`
	General   GeneralIDose   `json:"general"`
}

// MetricsIDose represents the IntelliDose metrics
type MetricsIDose struct {
	Ec      float64 `json:"ec"`
	NutTemp float64 `json:"nut_temp"`
	PH      float64 `json:"pH"`
}

// SettingsIDose represents the top level of the IntelliDose status
type SettingsIDose struct {
	General   GeneralStatusIDose `json:"general"`
	Nutrient  AlarmIDose         `json:"nutrient"`
	SetPoints SetPointsIDose     `json:"set_points"`
	Status    []StatusIDose      `json:"status"`
	Units     UnitsIDose         `json:"units"`
}

// GeneralStatusIDose represents the general status of the intellidose
type GeneralStatusIDose struct {
	DoseInterval        byte                    `json:"dose_interval"`
	NutrientDoseTime    byte                    `json:"nutrient_dose_time"`
	IrrigationInterval1 IrrigationIntervalIDose `json:"irrigation_interval_1"`
	IrrigationInterval2 IrrigationIntervalIDose `json:"irrigation_interval_2"`
	IrrigationInterval3 IrrigationIntervalIDose `json:"irrigation_interval_3"`
	IrrigationInterval4 IrrigationIntervalIDose `json:"irrigation_interval_4"`
	IrrigationDuration1 int                     `json:"irrigation_duration_1"`
	IrrigationDuration2 int                     `json:"irrigation_duration_2"`
	IrrigationDuration3 int                     `json:"irrigation_duration_3"`
	IrrigationDuration4 int                     `json:"irrigation_duration_4"`
	MaxNutrientDoseTime byte                    `json:"max_nutrient_dose_time"`
	MaxPhDoseTime       byte                    `json:"max_ph_dose_time"`
	Mix1                byte                    `json:"mix_1"`
	Mix2                byte                    `json:"mix_2"`
	Mix3                byte                    `json:"mix_3"`
	Mix4                byte                    `json:"mix_4"`
	Mix5                byte                    `json:"mix_5"`
	Mix6                byte                    `json:"mix_6"`
	Mix7                byte                    `json:"mix_7"`
	Mix8                byte                    `json:"mix_8"`
	PhDoseTime          byte                    `json:"ph_dose_time"`
}

// IrrigationIntervalIDose represents the irrigation interval settings of an IntelliDose
type IrrigationIntervalIDose struct {
	Days  int `json:"days"`
	Day   int `json:"day"`
	Night int `json:"night"`
}

// AlarmIDose represents the alarm settings for an IntelliDose
type AlarmIDose struct {
	Detent  byte              `json:"detent"`
	Ec      AlarmEcIDose      `json:"ec"`
	NutTemp AlarmNutTempIDose `json:"nut_temp"`
	Ph      AlarmPhIDose      `json:"ph"`
}

// AlarmEcIDose represents the EC settings of an IntelliDose
type AlarmEcIDose struct {
	Enabled bool    `json:"enabled"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
}

// AlarmNutTempIDose represents the nutrient temp settings of an IntelliDose
type AlarmNutTempIDose struct {
	Enabled bool    `json:"enabled"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
}

// AlarmPhIDose represents the pH temp settings of an IntelliDose
type AlarmPhIDose struct {
	Enabled bool    `json:"enabled"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
}

// SetPointsIDose represents the set points settings for an IntelliDose
type SetPointsIDose struct {
	Nutrient float64 `json:"nutrient"`
	PhDosing string  `json:"ph_dosing"`
	Ph       float64 `json:"ph"`
}

// StatusIDose represents the current status of an IntelliDose
type StatusIDose struct {
	Active   bool   `json:"active"`
	Enabled  bool   `json:"enabled"`
	ForceOn  bool   `json:"force_on"`
	Function string `json:"function"`
}

// UnitsIDose represents the reading units used
type UnitsIDose struct {
	DateFormat              string `json:"date_format"`
	Temperature             string `json:"temperature"`
	Ec                      string `json:"ec"`
	TdsConversationStandard string `json:"tds_conversation_standart"`
}

// TimesIDose represents the day start and day end times
type TimesIDose struct {
	DayStart string `json:"day_start"`
	DayEnd   string `json:"day_end"`
}

// FunctionsIDose represents the configuration of the functions in the IntelliDose
type FunctionsIDose struct {
	NutrientsParts     byte   `json:"nutrients_parts"`
	PhDosing           string `json:"ph_dosing"`
	IrrigationMode     string `json:"irrigation_mode"`
	IrrigationStations byte   `json:"irrigation_stations"`
	SeparatePumpOutput bool   `json:"separate_pump_output"`
	UseWater           bool   `json:"use_water"`
	ExternalAlarm      bool   `json:"external_alarm"`
	DayNightEc         bool   `json:"day_night_ec"`
	IrrigationStation1 string `json:"irrigation_station_1"`
	IrrigationStation2 string `json:"irrigation_station_2"`
	IrrigationStation3 string `json:"irrigation_station_3"`
	IrrigationStation4 string `json:"irrigation_station_4"`
	Scheduling         bool   `json:"scheduling"`
	MuteBuzzer         bool   `json:"mute_buzzer"`
}

// AdvancedIDose represents some advanced settings for the IntelliDose
type AdvancedIDose struct {
	ProportinalDosing bool   `json:"proportinal_dosing"`
	SequentialDosing  bool   `json:"sequential_dosing"`
	DisableEc         bool   `json:"disable_ec"`
	DisablePh         bool   `json:"disable_ph"`
	MntnReminderFreq  string `json:"mntn_reminder_freq"`
}

// GeneralIDose represents some general settings of the IntelliDose
type GeneralIDose struct {
	Growroom   string `json:"growroom"`
	DeviceName string `json:"device_name"`
}
