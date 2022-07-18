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

// Summary defines model for Summary.
type Summary struct {
	FirstData                *time.Time      `json:"firstData,omitempty"`
	HighGlucoseThreshold     *float64        `json:"highGlucoseThreshold,omitempty"`
	HourlyStats              *[]SummaryStat  `json:"hourlyStats,omitempty"`
	LastData                 *time.Time      `json:"lastData,omitempty"`
	LastUpdatedDate          *time.Time      `json:"lastUpdatedDate,omitempty"`
	LastUploadDate           *time.Time      `json:"lastUploadDate,omitempty"`
	LowGlucoseThreshold      *float64        `json:"lowGlucoseThreshold,omitempty"`
	OutdatedSince            *time.Time      `json:"outdatedSince,omitempty"`
	Periods                  *SummaryPeriods `json:"periods,omitempty"`
	TotalHours               *int            `json:"totalHours,omitempty"`
	VeryHighGlucoseThreshold *float64        `json:"veryHighGlucoseThreshold,omitempty"`
	VeryLowGlucoseThreshold  *float64        `json:"veryLowGlucoseThreshold,omitempty"`
}

// SummaryPeriod defines model for SummaryPeriod.
type SummaryPeriod struct {
	// Blood glucose value, in `mmol/L`
	AvgGlucose                 *AverageGlucose `json:"avgGlucose,omitempty"`
	GlucoseManagementIndicator *float64        `json:"glucoseManagementIndicator,omitempty"`
	TimeCGMUseMinutes          *int            `json:"timeCGMUseMinutes,omitempty"`
	TimeCGMUsePercent          *float64        `json:"timeCGMUsePercent,omitempty"`
	TimeCGMUseRecords          *int            `json:"timeCGMUseRecords,omitempty"`
	TimeInHighMinutes          *int            `json:"timeInHighMinutes,omitempty"`
	TimeInHighPercent          *float64        `json:"timeInHighPercent,omitempty"`
	TimeInHighRecords          *int            `json:"timeInHighRecords,omitempty"`
	TimeInLowMinutes           *int            `json:"timeInLowMinutes,omitempty"`
	TimeInLowPercent           *float64        `json:"timeInLowPercent,omitempty"`
	TimeInLowRecords           *int            `json:"timeInLowRecords,omitempty"`
	TimeInTargetMinutes        *int            `json:"timeInTargetMinutes,omitempty"`
	TimeInTargetPercent        *float64        `json:"timeInTargetPercent,omitempty"`
	TimeInTargetRecords        *int            `json:"timeInTargetRecords,omitempty"`
	TimeInVeryHighMinutes      *int            `json:"timeInVeryHighMinutes,omitempty"`
	TimeInVeryHighPercent      *float64        `json:"timeInVeryHighPercent,omitempty"`
	TimeInVeryHighRecords      *int            `json:"timeInVeryHighRecords,omitempty"`
	TimeInVeryLowMinutes       *int            `json:"timeInVeryLowMinutes,omitempty"`
	TimeInVeryLowPercent       *float64        `json:"timeInVeryLowPercent,omitempty"`
	TimeInVeryLowRecords       *int            `json:"timeInVeryLowRecords,omitempty"`
}

// SummaryPeriods defines model for SummaryPeriods.
type SummaryPeriods struct {
	N14d *SummaryPeriod `json:"14d,omitempty"`
	N1d  *SummaryPeriod `json:"1d,omitempty"`
	N30d *SummaryPeriod `json:"30d,omitempty"`
	N7d  *SummaryPeriod `json:"7d,omitempty"`
}

// SummaryStat defines model for SummaryStat.
type SummaryStat struct {
	Date                  *time.Time `json:"date,omitempty"`
	DeviceId              *string    `json:"deviceId,omitempty"`
	LastRecordTime        *time.Time `json:"lastRecordTime,omitempty"`
	TimeCGMUseMinutes     *int       `json:"timeCGMUseMinutes,omitempty"`
	TimeCGMUseRecords     *int       `json:"timeCGMUseRecords,omitempty"`
	TimeInHighMinutes     *int       `json:"timeInHighMinutes,omitempty"`
	TimeInHighRecords     *int       `json:"timeInHighRecords,omitempty"`
	TimeInLowMinutes      *int       `json:"timeInLowMinutes,omitempty"`
	TimeInLowRecords      *int       `json:"timeInLowRecords,omitempty"`
	TimeInTargetMinutes   *int       `json:"timeInTargetMinutes,omitempty"`
	TimeInTargetRecords   *int       `json:"timeInTargetRecords,omitempty"`
	TimeInVeryHighMinutes *int       `json:"timeInVeryHighMinutes,omitempty"`
	TimeInVeryHighRecords *int       `json:"timeInVeryHighRecords,omitempty"`
	TimeInVeryLowMinutes  *int       `json:"timeInVeryLowMinutes,omitempty"`
	TimeInVeryLowRecords  *int       `json:"timeInVeryLowRecords,omitempty"`
	TotalGlucose          *float64   `json:"totalGlucose,omitempty"`
}

// String representation of a Tidepool User ID. Old style IDs are 10-digit strings consisting of only hexadeximcal digits. New style IDs are 36-digit [UUID v4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random))
type TidepoolUserId string

// String representation of a Tidepool User ID. Old style IDs are 10-digit strings consisting of only hexadeximcal digits. New style IDs are 36-digit [UUID v4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random))
type UserId TidepoolUserId

