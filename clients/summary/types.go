// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package api

import (
	"time"
)

const (
	SessionTokenScopes = "sessionToken.Scopes"
)

// Defines values for AverageGlucoseUnits.
const (
	AverageGlucoseUnitsMmolL AverageGlucoseUnits = "mmol/L"

	AverageGlucoseUnitsMmoll AverageGlucoseUnits = "mmol/l"
)

// Blood glucose value, in `mmol/L`
type AverageGlucose struct {
	Units AverageGlucoseUnits `json:"units"`

	// A floating point value representing a `mmol/L` value.
	Value float32 `json:"value"`
}

// AverageGlucoseUnits defines model for AverageGlucose.Units.
type AverageGlucoseUnits string

// Series of counters which represent one hour of a users data
type BGMBucketData struct {
	Date           *time.Time `json:"date,omitempty"`
	LastRecordTime *time.Time `json:"lastRecordTime,omitempty"`

	// Counter of records in high glucose range
	TimeInHighRecords *int `json:"timeInHighRecords,omitempty"`

	// Counter of records in low glucose range
	TimeInLowRecords *int `json:"timeInLowRecords,omitempty"`

	// Counter of records in target glucose range
	TimeInTargetRecords *int `json:"timeInTargetRecords,omitempty"`

	// Counter of records in very high glucose range
	TimeInVeryHighRecords *int `json:"timeInVeryHighRecords,omitempty"`

	// Counter of records in very low glucose range
	TimeInVeryLowRecords *int `json:"timeInVeryLowRecords,omitempty"`

	// Total value of all glucose records
	TotalGlucose *float64 `json:"totalGlucose,omitempty"`
}

// Summary of a specific BGM time period (currently: 1d, 7d, 14d, 30d)
type BGMPeriod struct {
	// Blood glucose value, in `mmol/L`
	AverageGlucose           *AverageGlucose `json:"averageGlucose,omitempty"`
	HasAverageGlucose        *bool           `json:"hasAverageGlucose,omitempty"`
	HasTimeInHighPercent     *bool           `json:"hasTimeInHighPercent,omitempty"`
	HasTimeInLowPercent      *bool           `json:"hasTimeInLowPercent,omitempty"`
	HasTimeInTargetPercent   *bool           `json:"hasTimeInTargetPercent,omitempty"`
	HasTimeInVeryHighPercent *bool           `json:"hasTimeInVeryHighPercent,omitempty"`
	HasTimeInVeryLowPercent  *bool           `json:"hasTimeInVeryLowPercent,omitempty"`

	// Percentage of time spent in high glucose range
	TimeInHighPercent *float64 `json:"timeInHighPercent,omitempty"`

	// Counter of records in high glucose range
	TimeInHighRecords *int `json:"timeInHighRecords,omitempty"`

	// Percentage of time spent in low glucose range
	TimeInLowPercent *float64 `json:"timeInLowPercent,omitempty"`

	// Counter of records in low glucose range
	TimeInLowRecords *int `json:"timeInLowRecords,omitempty"`

	// Percentage of time spent in target glucose range
	TimeInTargetPercent *float64 `json:"timeInTargetPercent,omitempty"`

	// Counter of records in target glucose range
	TimeInTargetRecords *int `json:"timeInTargetRecords,omitempty"`

	// Percentage of time spent in very high glucose range
	TimeInVeryHighPercent *float64 `json:"timeInVeryHighPercent,omitempty"`

	// Counter of records in very high glucose range
	TimeInVeryHighRecords *int `json:"timeInVeryHighRecords,omitempty"`

	// Percentage of time spent in very low glucose range
	TimeInVeryLowPercent *float64 `json:"timeInVeryLowPercent,omitempty"`

	// Counter of records in very low glucose range
	TimeInVeryLowRecords *int `json:"timeInVeryLowRecords,omitempty"`
}

// A map to each supported BGM summary period
type BGMPeriods struct {
	// Summary of a specific BGM time period (currently: 1d, 7d, 14d, 30d)
	N14d *BGMPeriod `json:"14d,omitempty"`

	// Summary of a specific BGM time period (currently: 1d, 7d, 14d, 30d)
	N1d *BGMPeriod `json:"1d,omitempty"`

	// Summary of a specific BGM time period (currently: 1d, 7d, 14d, 30d)
	N30d *BGMPeriod `json:"30d,omitempty"`

	// Summary of a specific BGM time period (currently: 1d, 7d, 14d, 30d)
	N7d *BGMPeriod `json:"7d,omitempty"`
}

// A summary of a users recent BGM glucose values
type BGMStats struct {
	// Rotating list containing the stats for each currently tracked hour in order
	Buckets *[]Bucket `json:"buckets,omitempty"`

	// A map to each supported BGM summary period
	Periods *BGMPeriods `json:"periods,omitempty"`

	// Total hours represented in the hourly stats
	TotalHours *int `json:"totalHours,omitempty"`
}

// bucket containing an hour of bgm or cgm aggregations
type Bucket struct {
	Data           *interface{} `json:"data,omitempty"`
	Date           *time.Time   `json:"date,omitempty"`
	LastRecordTime *time.Time   `json:"lastRecordTime,omitempty"`
}

// Series of counters which represent one hour of a users data
type CGMBucketData struct {
	// Counter of minutes using a cgm
	TimeCGMUseMinutes *int `json:"timeCGMUseMinutes,omitempty"`

	// Counter of records wearing a cgm
	TimeCGMUseRecords *int `json:"timeCGMUseRecords,omitempty"`

	// Counter of minutes spent in high glucose range
	TimeInHighMinutes *int `json:"timeInHighMinutes,omitempty"`

	// Counter of records in high glucose range
	TimeInHighRecords *int `json:"timeInHighRecords,omitempty"`

	// Counter of minutes spent in low glucose range
	TimeInLowMinutes *int `json:"timeInLowMinutes,omitempty"`

	// Counter of records in low glucose range
	TimeInLowRecords *int `json:"timeInLowRecords,omitempty"`

	// Counter of minutes spent in target glucose range
	TimeInTargetMinutes *int `json:"timeInTargetMinutes,omitempty"`

	// Counter of records in target glucose range
	TimeInTargetRecords *int `json:"timeInTargetRecords,omitempty"`

	// Counter of minutes spent in very high glucose range
	TimeInVeryHighMinutes *int `json:"timeInVeryHighMinutes,omitempty"`

	// Counter of records in very high glucose range
	TimeInVeryHighRecords *int `json:"timeInVeryHighRecords,omitempty"`

	// Counter of minutes spent in very low glucose range
	TimeInVeryLowMinutes *int `json:"timeInVeryLowMinutes,omitempty"`

	// Counter of records in very low glucose range
	TimeInVeryLowRecords *int `json:"timeInVeryLowRecords,omitempty"`

	// Total value of all glucose records
	TotalGlucose *float64 `json:"totalGlucose,omitempty"`
}

// Summary of a specific CGM time period (currently: 1d, 7d, 14d, 30d)
type CGMPeriod struct {
	// Blood glucose value, in `mmol/L`
	AverageGlucose *AverageGlucose `json:"averageGlucose,omitempty"`

	// A derived value which emulates A1C
	GlucoseManagementIndicator    *float64 `json:"glucoseManagementIndicator,omitempty"`
	HasAverageGlucose             *bool    `json:"hasAverageGlucose,omitempty"`
	HasGlucoseManagementIndicator *bool    `json:"hasGlucoseManagementIndicator,omitempty"`
	HasTimeCGMUsePercent          *bool    `json:"hasTimeCGMUsePercent,omitempty"`
	HasTimeInHighPercent          *bool    `json:"hasTimeInHighPercent,omitempty"`
	HasTimeInLowPercent           *bool    `json:"hasTimeInLowPercent,omitempty"`
	HasTimeInTargetPercent        *bool    `json:"hasTimeInTargetPercent,omitempty"`
	HasTimeInVeryHighPercent      *bool    `json:"hasTimeInVeryHighPercent,omitempty"`
	HasTimeInVeryLowPercent       *bool    `json:"hasTimeInVeryLowPercent,omitempty"`

	// Counter of minutes spent wearing a cgm
	TimeCGMUseMinutes *int `json:"timeCGMUseMinutes,omitempty"`

	// Percentage of time spent wearing a cgm
	TimeCGMUsePercent *float64 `json:"timeCGMUsePercent,omitempty"`

	// Counter of minutes spent wearing a cgm
	TimeCGMUseRecords *int `json:"timeCGMUseRecords,omitempty"`

	// Counter of minutes spent in high glucose range
	TimeInHighMinutes *int `json:"timeInHighMinutes,omitempty"`

	// Percentage of time spent in high glucose range
	TimeInHighPercent *float64 `json:"timeInHighPercent,omitempty"`

	// Counter of records in high glucose range
	TimeInHighRecords *int `json:"timeInHighRecords,omitempty"`

	// Counter of minutes spent in low glucose range
	TimeInLowMinutes *int `json:"timeInLowMinutes,omitempty"`

	// Percentage of time spent in low glucose range
	TimeInLowPercent *float64 `json:"timeInLowPercent,omitempty"`

	// Counter of records in low glucose range
	TimeInLowRecords *int `json:"timeInLowRecords,omitempty"`

	// Counter of minutes spent in target glucose range
	TimeInTargetMinutes *int `json:"timeInTargetMinutes,omitempty"`

	// Percentage of time spent in target glucose range
	TimeInTargetPercent *float64 `json:"timeInTargetPercent,omitempty"`

	// Counter of records in target glucose range
	TimeInTargetRecords *int `json:"timeInTargetRecords,omitempty"`

	// Counter of minutes spent in very high glucose range
	TimeInVeryHighMinutes *int `json:"timeInVeryHighMinutes,omitempty"`

	// Percentage of time spent in very high glucose range
	TimeInVeryHighPercent *float64 `json:"timeInVeryHighPercent,omitempty"`

	// Counter of records in very high glucose range
	TimeInVeryHighRecords *int `json:"timeInVeryHighRecords,omitempty"`

	// Counter of minutes spent in very low glucose range
	TimeInVeryLowMinutes *int `json:"timeInVeryLowMinutes,omitempty"`

	// Percentage of time spent in very low glucose range
	TimeInVeryLowPercent *float64 `json:"timeInVeryLowPercent,omitempty"`

	// Counter of records in very low glucose range
	TimeInVeryLowRecords *int `json:"timeInVeryLowRecords,omitempty"`
}

// A map to each supported CGM summary period
type CGMPeriods struct {
	// Summary of a specific CGM time period (currently: 1d, 7d, 14d, 30d)
	N14d *CGMPeriod `json:"14d,omitempty"`

	// Summary of a specific CGM time period (currently: 1d, 7d, 14d, 30d)
	N1d *CGMPeriod `json:"1d,omitempty"`

	// Summary of a specific CGM time period (currently: 1d, 7d, 14d, 30d)
	N30d *CGMPeriod `json:"30d,omitempty"`

	// Summary of a specific CGM time period (currently: 1d, 7d, 14d, 30d)
	N7d *CGMPeriod `json:"7d,omitempty"`
}

// A summary of a users recent CGM glucose values
type CGMStats struct {
	// Rotating list containing the stats for each currently tracked hour in order
	Buckets *[]Bucket `json:"buckets,omitempty"`

	// A map to each supported CGM summary period
	Periods *CGMPeriods `json:"periods,omitempty"`

	// Total hours represented in the hourly stats
	TotalHours *int `json:"totalHours,omitempty"`
}

// Summary schema version and calculation configuration
type Config struct {
	// Threshold used for determining if a value is high
	HighGlucoseThreshold *float64 `json:"highGlucoseThreshold,omitempty"`

	// Threshold used for determining if a value is low
	LowGlucoseThreshold *float64 `json:"lowGlucoseThreshold,omitempty"`

	// Summary schema version
	SchemaVersion *int `json:"schemaVersion,omitempty"`

	// Threshold used for determining if a value is very high
	VeryHighGlucoseThreshold *float64 `json:"veryHighGlucoseThreshold,omitempty"`

	// Threshold used for determining if a value is very low
	VeryLowGlucoseThreshold *float64 `json:"veryLowGlucoseThreshold,omitempty"`
}

// dates tracked for summary calculation
type Dates struct {
	// Date of the first included value
	FirstData         *time.Time `json:"firstData,omitempty"`
	HasLastUploadDate *bool      `json:"hasLastUploadDate,omitempty"`

	// Date of the last calculated value
	LastData *time.Time `json:"lastData,omitempty"`

	// Date of the last calculation
	LastUpdatedDate *time.Time `json:"lastUpdatedDate,omitempty"`

	// Created date of the last calculated value
	LastUploadDate *time.Time `json:"lastUploadDate,omitempty"`

	// Date of the first user upload after lastData, removed when calculated
	OutdatedSince *time.Time `json:"outdatedSince,omitempty"`
}

// A summary of a users recent data
type Summary struct {
	// Summary schema version and calculation configuration
	Config *Config `json:"config,omitempty"`

	// dates tracked for summary calculation
	Dates *Dates       `json:"dates,omitempty"`
	Stats *interface{} `json:"stats,omitempty"`

	// Field which contains a summary type string.
	Type *SummaryType `json:"type,omitempty"`

	// String representation of a Tidepool User ID. Old style IDs are 10-digit strings consisting of only hexadeximcal digits. New style IDs are 36-digit [UUID v4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random))
	UserId *TidepoolUserId `json:"userId,omitempty"`
}

// Field which contains a summary type string.
type SummaryType string

// String representation of a Tidepool User ID. Old style IDs are 10-digit strings consisting of only hexadeximcal digits. New style IDs are 36-digit [UUID v4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random))
type TidepoolUserId string

// String representation of a Tidepool User ID. Old style IDs are 10-digit strings consisting of only hexadeximcal digits. New style IDs are 36-digit [UUID v4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random))
type UserId TidepoolUserId

