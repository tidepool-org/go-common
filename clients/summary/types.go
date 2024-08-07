// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package api

import (
	"encoding/json"
	"time"

	"github.com/oapi-codegen/runtime"
)

const (
	SessionTokenScopes = "sessionToken.Scopes"
)

// BGMBucketData Series of counters which represent one hour of a users data
type BGMBucketData struct {
	LastRecordTime *time.Time `json:"lastRecordTime,omitempty"`

	// TimeInHighRecords Counter of records in high glucose range
	TimeInHighRecords *int `json:"timeInHighRecords,omitempty"`

	// TimeInLowRecords Counter of records in low glucose range
	TimeInLowRecords *int `json:"timeInLowRecords,omitempty"`

	// TimeInTargetRecords Counter of records in target glucose range
	TimeInTargetRecords *int `json:"timeInTargetRecords,omitempty"`

	// TimeInVeryHighRecords Counter of records in very high glucose range
	TimeInVeryHighRecords *int `json:"timeInVeryHighRecords,omitempty"`

	// TimeInVeryLowRecords Counter of records in very low glucose range
	TimeInVeryLowRecords *int `json:"timeInVeryLowRecords,omitempty"`

	// TotalGlucose Total value of all glucose records
	TotalGlucose float64 `json:"totalGlucose"`
}

// BGMPeriod Summary of a specific BGM time period (currently: 1d, 7d, 14d, 30d)
type BGMPeriod struct {
	// AverageDailyRecords Average daily readings
	AverageDailyRecords *float64 `json:"averageDailyRecords,omitempty"`

	// AverageDailyRecordsDelta Difference between the averageDailyRecords in this period and version in the opposite offset
	AverageDailyRecordsDelta *float64 `json:"averageDailyRecordsDelta,omitempty"`

	// AverageGlucoseMmol Average Glucose of records in this period
	AverageGlucoseMmol *float64 `json:"averageGlucoseMmol,omitempty"`

	// AverageGlucoseMmolDelta Difference between the averageGlucose in this period and the other offset version
	AverageGlucoseMmolDelta     *float64 `json:"averageGlucoseMmolDelta,omitempty"`
	HasAverageDailyRecords      bool     `json:"hasAverageDailyRecords"`
	HasAverageGlucoseMmol       bool     `json:"hasAverageGlucoseMmol"`
	HasTimeInAnyHighPercent     bool     `json:"hasTimeInAnyHighPercent"`
	HasTimeInAnyHighRecords     bool     `json:"hasTimeInAnyHighRecords"`
	HasTimeInAnyLowPercent      bool     `json:"hasTimeInAnyLowPercent"`
	HasTimeInAnyLowRecords      bool     `json:"hasTimeInAnyLowRecords"`
	HasTimeInExtremeHighPercent bool     `json:"hasTimeInExtremeHighPercent"`
	HasTimeInExtremeHighRecords bool     `json:"hasTimeInExtremeHighRecords"`
	HasTimeInHighPercent        bool     `json:"hasTimeInHighPercent"`
	HasTimeInHighRecords        bool     `json:"hasTimeInHighRecords"`
	HasTimeInLowPercent         bool     `json:"hasTimeInLowPercent"`
	HasTimeInLowRecords         bool     `json:"hasTimeInLowRecords"`
	HasTimeInTargetPercent      bool     `json:"hasTimeInTargetPercent"`
	HasTimeInTargetRecords      bool     `json:"hasTimeInTargetRecords"`
	HasTimeInVeryHighPercent    bool     `json:"hasTimeInVeryHighPercent"`
	HasTimeInVeryHighRecords    bool     `json:"hasTimeInVeryHighRecords"`
	HasTimeInVeryLowPercent     bool     `json:"hasTimeInVeryLowPercent"`
	HasTimeInVeryLowRecords     bool     `json:"hasTimeInVeryLowRecords"`
	HasTotalRecords             bool     `json:"hasTotalRecords"`

	// TimeInAnyHighPercent Percentage of time spent in Any high glucose range
	TimeInAnyHighPercent *float64 `json:"timeInAnyHighPercent,omitempty"`

	// TimeInAnyHighPercentDelta Difference between the timeInAnyHighPercent in this period and version in the opposite offset
	TimeInAnyHighPercentDelta *float64 `json:"timeInAnyHighPercentDelta,omitempty"`

	// TimeInAnyHighRecords Counter of records in Any high glucose range
	TimeInAnyHighRecords *int `json:"timeInAnyHighRecords,omitempty"`

	// TimeInAnyHighRecordsDelta Difference between the timeInAnyHighRecords in this period and version in the opposite offset
	TimeInAnyHighRecordsDelta *int `json:"timeInAnyHighRecordsDelta,omitempty"`

	// TimeInAnyLowPercent Percentage of time spent in Any low glucose range
	TimeInAnyLowPercent *float64 `json:"timeInAnyLowPercent,omitempty"`

	// TimeInAnyLowPercentDelta Difference between the timeInAnyLowPercent in this period and version in the opposite offset
	TimeInAnyLowPercentDelta *float64 `json:"timeInAnyLowPercentDelta,omitempty"`

	// TimeInAnyLowRecords Counter of records in Any low glucose range
	TimeInAnyLowRecords *int `json:"timeInAnyLowRecords,omitempty"`

	// TimeInAnyLowRecordsDelta Difference between the timeInAnyLowRecords in this period and version in the opposite offset
	TimeInAnyLowRecordsDelta *int `json:"timeInAnyLowRecordsDelta,omitempty"`

	// TimeInExtremeHighPercent Percentage of time spent in extreme high glucose range
	TimeInExtremeHighPercent *float64 `json:"timeInExtremeHighPercent,omitempty"`

	// TimeInExtremeHighPercentDelta Difference between the timeInExtremeHighPercent in this period and version in the opposite offset
	TimeInExtremeHighPercentDelta *float64 `json:"timeInExtremeHighPercentDelta,omitempty"`

	// TimeInExtremeHighRecords Counter of records in extreme high glucose range
	TimeInExtremeHighRecords *int `json:"timeInExtremeHighRecords,omitempty"`

	// TimeInExtremeHighRecordsDelta Difference between the timeInExtremeHighRecords in this period and version in the opposite offset
	TimeInExtremeHighRecordsDelta *int `json:"timeInExtremeHighRecordsDelta,omitempty"`

	// TimeInHighPercent Percentage of time spent in high glucose range
	TimeInHighPercent *float64 `json:"timeInHighPercent,omitempty"`

	// TimeInHighPercentDelta Difference between the timeInHighPercent in this period and version in the opposite offset
	TimeInHighPercentDelta *float64 `json:"timeInHighPercentDelta,omitempty"`

	// TimeInHighRecords Counter of records in high glucose range
	TimeInHighRecords *int `json:"timeInHighRecords,omitempty"`

	// TimeInHighRecordsDelta Difference between the timeInHighRecords in this period and version in the opposite offset
	TimeInHighRecordsDelta *int `json:"timeInHighRecordsDelta,omitempty"`

	// TimeInLowPercent Percentage of time spent in low glucose range
	TimeInLowPercent *float64 `json:"timeInLowPercent,omitempty"`

	// TimeInLowPercentDelta Difference between the timeInLowPercent in this period and version in the opposite offset
	TimeInLowPercentDelta *float64 `json:"timeInLowPercentDelta,omitempty"`

	// TimeInLowRecords Counter of records in low glucose range
	TimeInLowRecords *int `json:"timeInLowRecords,omitempty"`

	// TimeInLowRecordsDelta Difference between the timeInLowRecords in this period and version in the opposite offset
	TimeInLowRecordsDelta *int `json:"timeInLowRecordsDelta,omitempty"`

	// TimeInTargetPercent Percentage of time spent in target glucose range
	TimeInTargetPercent *float64 `json:"timeInTargetPercent,omitempty"`

	// TimeInTargetPercentDelta Difference between the timeInTargetPercent in this period and version in the opposite offset
	TimeInTargetPercentDelta *float64 `json:"timeInTargetPercentDelta,omitempty"`

	// TimeInTargetRecords Counter of records in target glucose range
	TimeInTargetRecords *int `json:"timeInTargetRecords,omitempty"`

	// TimeInTargetRecordsDelta Difference between the timeInTargetRecords in this period and version in the opposite offset
	TimeInTargetRecordsDelta *int `json:"timeInTargetRecordsDelta,omitempty"`

	// TimeInVeryHighPercent Percentage of time spent in very high glucose range
	TimeInVeryHighPercent *float64 `json:"timeInVeryHighPercent,omitempty"`

	// TimeInVeryHighPercentDelta Difference between the timeInVeryHighPercent in this period and version in the opposite offset
	TimeInVeryHighPercentDelta *float64 `json:"timeInVeryHighPercentDelta,omitempty"`

	// TimeInVeryHighRecords Counter of records in very high glucose range
	TimeInVeryHighRecords *int `json:"timeInVeryHighRecords,omitempty"`

	// TimeInVeryHighRecordsDelta Difference between the timeInVeryHighRecords in this period and version in the opposite offset
	TimeInVeryHighRecordsDelta *int `json:"timeInVeryHighRecordsDelta,omitempty"`

	// TimeInVeryLowPercent Percentage of time spent in very low glucose range
	TimeInVeryLowPercent *float64 `json:"timeInVeryLowPercent,omitempty"`

	// TimeInVeryLowPercentDelta Difference between the timeInVeryLowPercent in this period and version in the opposite offset
	TimeInVeryLowPercentDelta *float64 `json:"timeInVeryLowPercentDelta,omitempty"`

	// TimeInVeryLowRecords Counter of records in very low glucose range
	TimeInVeryLowRecords *int `json:"timeInVeryLowRecords,omitempty"`

	// TimeInVeryLowRecordsDelta Difference between the timeInVeryLowRecords in this period and version in the opposite offset
	TimeInVeryLowRecordsDelta *int `json:"timeInVeryLowRecordsDelta,omitempty"`

	// TotalRecords Counter of records
	TotalRecords *int `json:"totalRecords,omitempty"`

	// TotalRecordsDelta Difference between the totalRecords in this period and version in the opposite offset
	TotalRecordsDelta *int `json:"totalRecordsDelta,omitempty"`
}

// BGMPeriods A map to each supported BGM summary period
type BGMPeriods map[string]BGMPeriod

// BGMStats A summary of a users recent BGM glucose values
type BGMStats struct {
	// Buckets Rotating list containing the stats for each currently tracked hour in order
	Buckets []Bucket `json:"buckets,omitempty"`

	// OffsetPeriods A map to each supported BGM summary period
	OffsetPeriods BGMPeriods `json:"offsetPeriods,omitempty"`

	// Periods A map to each supported BGM summary period
	Periods BGMPeriods `json:"periods,omitempty"`

	// TotalHours Total hours represented in the hourly stats
	TotalHours int `json:"totalHours"`
}

// Bucket bucket containing an hour of bgm or cgm aggregations
type Bucket struct {
	Data           *Bucket_Data `json:"data,omitempty"`
	Date           time.Time    `json:"date"`
	LastRecordTime *time.Time   `json:"lastRecordTime,omitempty"`
}

// Bucket_Data defines model for Bucket.Data.
type Bucket_Data struct {
	union json.RawMessage
}

// CGMBucketData Series of counters which represent one hour of a users data
type CGMBucketData struct {
	// HighMinutes Counter of minutes spent in high glucose range
	HighMinutes int `json:"highMinutes"`

	// HighRecords Counter of records in high glucose range
	HighRecords int `json:"highRecords"`

	// LowMinutes Counter of minutes spent in low glucose range
	LowMinutes int `json:"lowMinutes"`

	// LowRecords Counter of records in low glucose range
	LowRecords int `json:"lowRecords"`

	// TargetMinutes Counter of minutes spent in target glucose range
	TargetMinutes int `json:"targetMinutes"`

	// TargetRecords Counter of records in target glucose range
	TargetRecords int `json:"targetRecords"`

	// TotalGlucose Total value of all glucose records
	TotalGlucose float64 `json:"totalGlucose"`

	// TotalMinutes Counter of minutes using a cgm
	TotalMinutes int `json:"totalMinutes"`

	// TotalRecords Counter of records using a cgm
	TotalRecords int `json:"totalRecords"`

	// TotalVariance Total variance of all glucose records
	TotalVariance float64 `json:"totalVariance"`

	// VeryHighMinutes Counter of minutes spent in very high glucose range
	VeryHighMinutes int `json:"veryHighMinutes"`

	// VeryHighRecords Counter of records in very high glucose range
	VeryHighRecords int `json:"veryHighRecords"`

	// VeryLowMinutes Counter of minutes spent in very low glucose range
	VeryLowMinutes int `json:"veryLowMinutes"`

	// VeryLowRecords Counter of records in very low glucose range
	VeryLowRecords int `json:"veryLowRecords"`
}

// CGMPeriod Summary of a specific CGM time period (currently: 1d, 7d, 14d, 30d)
type CGMPeriod struct {
	// AverageDailyRecords Average daily readings
	AverageDailyRecords *float64 `json:"averageDailyRecords,omitempty"`

	// AverageDailyRecordsDelta Difference between the averageDailyRecords in this period and version in the opposite offset
	AverageDailyRecordsDelta *float64 `json:"averageDailyRecordsDelta,omitempty"`

	// AverageGlucoseMmol Average Glucose of records in this period
	AverageGlucoseMmol *float64 `json:"averageGlucoseMmol,omitempty"`

	// AverageGlucoseMmolDelta Difference between the averageGlucose in this period and the other offset version
	AverageGlucoseMmolDelta     *float64 `json:"averageGlucoseMmolDelta,omitempty"`
	CoefficientOfVariation      float64  `json:"coefficientOfVariation"`
	CoefficientOfVariationDelta float64  `json:"coefficientOfVariationDelta"`
	DaysWithData                int      `json:"daysWithData"`
	DaysWithDataDelta           int      `json:"daysWithDataDelta"`

	// GlucoseManagementIndicator A derived value which emulates A1C
	GlucoseManagementIndicator *float64 `json:"glucoseManagementIndicator,omitempty"`

	// GlucoseManagementIndicatorDelta Difference between the glucoseManagementIndicator in this period and the other offset version
	GlucoseManagementIndicatorDelta *float64 `json:"glucoseManagementIndicatorDelta,omitempty"`
	HasAverageDailyRecords          bool     `json:"hasAverageDailyRecords"`
	HasAverageGlucoseMmol           bool     `json:"hasAverageGlucoseMmol"`
	HasGlucoseManagementIndicator   bool     `json:"hasGlucoseManagementIndicator"`
	HasTimeCGMUseMinutes            bool     `json:"hasTimeCGMUseMinutes"`
	HasTimeCGMUsePercent            bool     `json:"hasTimeCGMUsePercent"`
	HasTimeCGMUseRecords            bool     `json:"hasTimeCGMUseRecords"`
	HasTimeInAnyHighMinutes         bool     `json:"hasTimeInAnyHighMinutes"`
	HasTimeInAnyHighPercent         bool     `json:"hasTimeInAnyHighPercent"`
	HasTimeInAnyHighRecords         bool     `json:"hasTimeInAnyHighRecords"`
	HasTimeInAnyLowMinutes          bool     `json:"hasTimeInAnyLowMinutes"`
	HasTimeInAnyLowPercent          bool     `json:"hasTimeInAnyLowPercent"`
	HasTimeInAnyLowRecords          bool     `json:"hasTimeInAnyLowRecords"`
	HasTimeInExtremeHighMinutes     bool     `json:"hasTimeInExtremeHighMinutes"`
	HasTimeInExtremeHighPercent     bool     `json:"hasTimeInExtremeHighPercent"`
	HasTimeInExtremeHighRecords     bool     `json:"hasTimeInExtremeHighRecords"`
	HasTimeInHighMinutes            bool     `json:"hasTimeInHighMinutes"`
	HasTimeInHighPercent            bool     `json:"hasTimeInHighPercent"`
	HasTimeInHighRecords            bool     `json:"hasTimeInHighRecords"`
	HasTimeInLowMinutes             bool     `json:"hasTimeInLowMinutes"`
	HasTimeInLowPercent             bool     `json:"hasTimeInLowPercent"`
	HasTimeInLowRecords             bool     `json:"hasTimeInLowRecords"`
	HasTimeInTargetMinutes          bool     `json:"hasTimeInTargetMinutes"`
	HasTimeInTargetPercent          bool     `json:"hasTimeInTargetPercent"`
	HasTimeInTargetRecords          bool     `json:"hasTimeInTargetRecords"`
	HasTimeInVeryHighMinutes        bool     `json:"hasTimeInVeryHighMinutes"`
	HasTimeInVeryHighPercent        bool     `json:"hasTimeInVeryHighPercent"`
	HasTimeInVeryHighRecords        bool     `json:"hasTimeInVeryHighRecords"`
	HasTimeInVeryLowMinutes         bool     `json:"hasTimeInVeryLowMinutes"`
	HasTimeInVeryLowPercent         bool     `json:"hasTimeInVeryLowPercent"`
	HasTimeInVeryLowRecords         bool     `json:"hasTimeInVeryLowRecords"`
	HasTotalRecords                 bool     `json:"hasTotalRecords"`
	HoursWithData                   int      `json:"hoursWithData"`
	HoursWithDataDelta              int      `json:"hoursWithDataDelta"`
	StandardDeviation               float64  `json:"standardDeviation"`
	StandardDeviationDelta          float64  `json:"standardDeviationDelta"`

	// TimeCGMUseMinutes Counter of minutes spent wearing a cgm
	TimeCGMUseMinutes *int `json:"timeCGMUseMinutes,omitempty"`

	// TimeCGMUseMinutesDelta Difference between the timeCGMUseMinutes in this period and version in the opposite offset
	TimeCGMUseMinutesDelta *int `json:"timeCGMUseMinutesDelta,omitempty"`

	// TimeCGMUsePercent Percentage of time spent wearing a cgm
	TimeCGMUsePercent *float64 `json:"timeCGMUsePercent,omitempty"`

	// TimeCGMUsePercentDelta Difference between the timeCGMUsePercent in this period and version in the opposite offset
	TimeCGMUsePercentDelta *float64 `json:"timeCGMUsePercentDelta,omitempty"`

	// TimeCGMUseRecords Counter of minutes spent wearing a cgm
	TimeCGMUseRecords *int `json:"timeCGMUseRecords,omitempty"`

	// TimeCGMUseRecordsDelta Difference between the timeCGMUseRecords in this period and version in the opposite offset
	TimeCGMUseRecordsDelta *int `json:"timeCGMUseRecordsDelta,omitempty"`

	// TimeInAnyHighMinutes Counter of minutes spent in Any high glucose range
	TimeInAnyHighMinutes *int `json:"timeInAnyHighMinutes,omitempty"`

	// TimeInAnyHighMinutesDelta Difference between the timeInAnyHighMinutes in this period and version in the opposite offset
	TimeInAnyHighMinutesDelta *int `json:"timeInAnyHighMinutesDelta,omitempty"`

	// TimeInAnyHighPercent Percentage of time spent in Any high glucose range
	TimeInAnyHighPercent *float64 `json:"timeInAnyHighPercent,omitempty"`

	// TimeInAnyHighPercentDelta Difference between the timeInAnyHighPercent in this period and version in the opposite offset
	TimeInAnyHighPercentDelta *float64 `json:"timeInAnyHighPercentDelta,omitempty"`

	// TimeInAnyHighRecords Counter of records in Any high glucose range
	TimeInAnyHighRecords *int `json:"timeInAnyHighRecords,omitempty"`

	// TimeInAnyHighRecordsDelta Difference between the timeInAnyHighRecords in this period and version in the opposite offset
	TimeInAnyHighRecordsDelta *int `json:"timeInAnyHighRecordsDelta,omitempty"`

	// TimeInAnyLowMinutes Counter of minutes spent in Any low glucose range
	TimeInAnyLowMinutes *int `json:"timeInAnyLowMinutes,omitempty"`

	// TimeInAnyLowMinutesDelta Difference between the timeInAnyLowMinutes in this period and version in the opposite offset
	TimeInAnyLowMinutesDelta *int `json:"timeInAnyLowMinutesDelta,omitempty"`

	// TimeInAnyLowPercent Percentage of time spent in Any low glucose range
	TimeInAnyLowPercent *float64 `json:"timeInAnyLowPercent,omitempty"`

	// TimeInAnyLowPercentDelta Difference between the timeInAnyLowPercent in this period and version in the opposite offset
	TimeInAnyLowPercentDelta *float64 `json:"timeInAnyLowPercentDelta,omitempty"`

	// TimeInAnyLowRecords Counter of records in Any low glucose range
	TimeInAnyLowRecords *int `json:"timeInAnyLowRecords,omitempty"`

	// TimeInAnyLowRecordsDelta Difference between the timeInAnyLowRecords in this period and version in the opposite offset
	TimeInAnyLowRecordsDelta *int `json:"timeInAnyLowRecordsDelta,omitempty"`

	// TimeInExtremeHighMinutes Counter of minutes spent in extreme high glucose range
	TimeInExtremeHighMinutes *int `json:"timeInExtremeHighMinutes,omitempty"`

	// TimeInExtremeHighMinutesDelta Difference between the timeInExtremeHighMinutes in this period and version in the opposite offset
	TimeInExtremeHighMinutesDelta *int `json:"timeInExtremeHighMinutesDelta,omitempty"`

	// TimeInExtremeHighPercent Percentage of time spent in extreme high glucose range
	TimeInExtremeHighPercent *float64 `json:"timeInExtremeHighPercent,omitempty"`

	// TimeInExtremeHighPercentDelta Difference between the timeInExtremeHighPercent in this period and version in the opposite offset
	TimeInExtremeHighPercentDelta *float64 `json:"timeInExtremeHighPercentDelta,omitempty"`

	// TimeInExtremeHighRecords Counter of records in extreme high glucose range
	TimeInExtremeHighRecords *int `json:"timeInExtremeHighRecords,omitempty"`

	// TimeInExtremeHighRecordsDelta Difference between the timeInExtremeHighRecords in this period and version in the opposite offset
	TimeInExtremeHighRecordsDelta *int `json:"timeInExtremeHighRecordsDelta,omitempty"`

	// TimeInHighMinutes Counter of minutes spent in high glucose range
	TimeInHighMinutes *int `json:"timeInHighMinutes,omitempty"`

	// TimeInHighMinutesDelta Difference between the timeInHighMinutes in this period and version in the opposite offset
	TimeInHighMinutesDelta *int `json:"timeInHighMinutesDelta,omitempty"`

	// TimeInHighPercent Percentage of time spent in high glucose range
	TimeInHighPercent *float64 `json:"timeInHighPercent,omitempty"`

	// TimeInHighPercentDelta Difference between the timeInHighPercent in this period and version in the opposite offset
	TimeInHighPercentDelta *float64 `json:"timeInHighPercentDelta,omitempty"`

	// TimeInHighRecords Counter of records in high glucose range
	TimeInHighRecords *int `json:"timeInHighRecords,omitempty"`

	// TimeInHighRecordsDelta Difference between the timeInHighRecords in this period and version in the opposite offset
	TimeInHighRecordsDelta *int `json:"timeInHighRecordsDelta,omitempty"`

	// TimeInLowMinutes Counter of minutes spent in low glucose range
	TimeInLowMinutes *int `json:"timeInLowMinutes,omitempty"`

	// TimeInLowMinutesDelta Difference between the timeInLowMinutes in this period and version in the opposite offset
	TimeInLowMinutesDelta *int `json:"timeInLowMinutesDelta,omitempty"`

	// TimeInLowPercent Percentage of time spent in low glucose range
	TimeInLowPercent *float64 `json:"timeInLowPercent,omitempty"`

	// TimeInLowPercentDelta Difference between the timeInLowPercent in this period and version in the opposite offset
	TimeInLowPercentDelta *float64 `json:"timeInLowPercentDelta,omitempty"`

	// TimeInLowRecords Counter of records in low glucose range
	TimeInLowRecords *int `json:"timeInLowRecords,omitempty"`

	// TimeInLowRecordsDelta Difference between the timeInLowRecords in this period and version in the opposite offset
	TimeInLowRecordsDelta *int `json:"timeInLowRecordsDelta,omitempty"`

	// TimeInTargetMinutes Counter of minutes spent in target glucose range
	TimeInTargetMinutes *int `json:"timeInTargetMinutes,omitempty"`

	// TimeInTargetMinutesDelta Difference between the timeInTargetMinutes in this period and version in the opposite offset
	TimeInTargetMinutesDelta *int `json:"timeInTargetMinutesDelta,omitempty"`

	// TimeInTargetPercent Percentage of time spent in target glucose range
	TimeInTargetPercent *float64 `json:"timeInTargetPercent,omitempty"`

	// TimeInTargetPercentDelta Difference between the timeInTargetPercent in this period and version in the opposite offset
	TimeInTargetPercentDelta *float64 `json:"timeInTargetPercentDelta,omitempty"`

	// TimeInTargetRecords Counter of records in target glucose range
	TimeInTargetRecords *int `json:"timeInTargetRecords,omitempty"`

	// TimeInTargetRecordsDelta Difference between the timeInTargetRecords in this period and version in the opposite offset
	TimeInTargetRecordsDelta *int `json:"timeInTargetRecordsDelta,omitempty"`

	// TimeInVeryHighMinutes Counter of minutes spent in very high glucose range
	TimeInVeryHighMinutes *int `json:"timeInVeryHighMinutes,omitempty"`

	// TimeInVeryHighMinutesDelta Difference between the timeInVeryHighMinutes in this period and version in the opposite offset
	TimeInVeryHighMinutesDelta *int `json:"timeInVeryHighMinutesDelta,omitempty"`

	// TimeInVeryHighPercent Percentage of time spent in very high glucose range
	TimeInVeryHighPercent *float64 `json:"timeInVeryHighPercent,omitempty"`

	// TimeInVeryHighPercentDelta Difference between the timeInVeryHighPercent in this period and version in the opposite offset
	TimeInVeryHighPercentDelta *float64 `json:"timeInVeryHighPercentDelta,omitempty"`

	// TimeInVeryHighRecords Counter of records in very high glucose range
	TimeInVeryHighRecords *int `json:"timeInVeryHighRecords,omitempty"`

	// TimeInVeryHighRecordsDelta Difference between the timeInVeryHighRecords in this period and version in the opposite offset
	TimeInVeryHighRecordsDelta *int `json:"timeInVeryHighRecordsDelta,omitempty"`

	// TimeInVeryLowMinutes Counter of minutes spent in very low glucose range
	TimeInVeryLowMinutes *int `json:"timeInVeryLowMinutes,omitempty"`

	// TimeInVeryLowMinutesDelta Difference between the timeInVeryLowMinutes in this period and version in the opposite offset
	TimeInVeryLowMinutesDelta *int `json:"timeInVeryLowMinutesDelta,omitempty"`

	// TimeInVeryLowPercent Percentage of time spent in very low glucose range
	TimeInVeryLowPercent *float64 `json:"timeInVeryLowPercent,omitempty"`

	// TimeInVeryLowPercentDelta Difference between the timeInVeryLowPercent in this period and version in the opposite offset
	TimeInVeryLowPercentDelta *float64 `json:"timeInVeryLowPercentDelta,omitempty"`

	// TimeInVeryLowRecords Counter of records in very low glucose range
	TimeInVeryLowRecords *int `json:"timeInVeryLowRecords,omitempty"`

	// TimeInVeryLowRecordsDelta Difference between the timeInVeryLowRecords in this period and version in the opposite offset
	TimeInVeryLowRecordsDelta *int `json:"timeInVeryLowRecordsDelta,omitempty"`

	// TotalRecords Counter of records
	TotalRecords *int `json:"totalRecords,omitempty"`

	// TotalRecordsDelta Difference between the totalRecords in this period and version in the opposite offset
	TotalRecordsDelta *int `json:"totalRecordsDelta,omitempty"`
}

// CGMPeriods A map to each supported CGM summary period
type CGMPeriods map[string]CGMPeriod

// CGMStats A summary of a users recent CGM glucose values
type CGMStats struct {
	// Buckets Rotating list containing the stats for each currently tracked hour in order
	Buckets []Bucket `json:"buckets,omitempty"`

	// OffsetPeriods A map to each supported CGM summary period
	OffsetPeriods CGMPeriods `json:"offsetPeriods,omitempty"`

	// Periods A map to each supported CGM summary period
	Periods CGMPeriods `json:"periods,omitempty"`

	// TotalHours Total hours represented in the hourly stats
	TotalHours int `json:"totalHours"`
}

// Config Summary schema version and calculation configuration
type Config struct {
	// HighGlucoseThreshold Threshold used for determining if a value is high
	HighGlucoseThreshold float64 `json:"highGlucoseThreshold"`

	// LowGlucoseThreshold Threshold used for determining if a value is low
	LowGlucoseThreshold float64 `json:"lowGlucoseThreshold"`

	// SchemaVersion Summary schema version
	SchemaVersion int `json:"schemaVersion"`

	// VeryHighGlucoseThreshold Threshold used for determining if a value is very high
	VeryHighGlucoseThreshold float64 `json:"veryHighGlucoseThreshold"`

	// VeryLowGlucoseThreshold Threshold used for determining if a value is very low
	VeryLowGlucoseThreshold float64 `json:"veryLowGlucoseThreshold"`
}

// ContinuousBucketData Series of counters which represent one hour of a users data
type ContinuousBucketData struct {
	// DeferredRecords Counter of records uploaded later than 24 hours of their timestamp
	DeferredRecords int `json:"deferredRecords"`

	// RealtimeRecords Counter of records uploaded within 24 hours of their timestamp
	RealtimeRecords int `json:"realtimeRecords"`

	// TotalRecords Counter of records from continuous uploads
	TotalRecords int `json:"totalRecords"`
}

// ContinuousPeriod Summary of a specific continuous upload time period (currently: 1d, 7d, 14d, 30d)
type ContinuousPeriod struct {
	AverageDailyRecords float64 `json:"averageDailyRecords"`
	DeferredPercent     float64 `json:"deferredPercent"`
	DeferredRecords     int     `json:"deferredRecords"`
	RealtimePercent     float64 `json:"realtimePercent"`
	RealtimeRecords     int     `json:"realtimeRecords"`
	TotalRecords        int     `json:"totalRecords"`
}

// ContinuousPeriods A map to each supported CGM summary period
type ContinuousPeriods map[string]ContinuousPeriod

// ContinuousStats A summary of a users recent CGM glucose values
type ContinuousStats struct {
	// Buckets Rotating list containing the stats for each currently tracked hour in order
	Buckets []Bucket `json:"buckets,omitempty"`

	// OffsetPeriods A map to each supported CGM summary period
	OffsetPeriods ContinuousPeriods `json:"offsetPeriods,omitempty"`

	// Periods A map to each supported CGM summary period
	Periods ContinuousPeriods `json:"periods,omitempty"`

	// TotalHours Total hours represented in the hourly stats
	TotalHours int `json:"totalHours"`
}

// Dates dates tracked for summary calculation
type Dates struct {
	// FirstData Date of the first included value
	FirstData         *time.Time `json:"firstData,omitempty"`
	HasFirstData      bool       `json:"hasFirstData"`
	HasLastData       bool       `json:"hasLastData"`
	HasLastUploadDate bool       `json:"hasLastUploadDate"`
	HasOutdatedSince  bool       `json:"hasOutdatedSince"`

	// LastData Date of the last calculated value
	LastData *time.Time `json:"lastData,omitempty"`

	// LastUpdatedDate Date of the last calculation
	LastUpdatedDate *time.Time `json:"lastUpdatedDate,omitempty"`

	// LastUpdatedReason List of reasons the summary was updated for
	LastUpdatedReason *[]string `json:"lastUpdatedReason,omitempty"`

	// LastUploadDate Created date of the last calculated value
	LastUploadDate *time.Time `json:"lastUploadDate,omitempty"`

	// OutdatedReason List of reasons the summary was marked outdated for
	OutdatedReason *[]string `json:"outdatedReason,omitempty"`

	// OutdatedSince Date of the first user upload after lastData, removed when calculated
	OutdatedSince *time.Time `json:"outdatedSince,omitempty"`

	// OutdatedSinceLimit Upper limit of the OutdatedSince value to prevent infinite queue duration
	OutdatedSinceLimit *time.Time `json:"outdatedSinceLimit,omitempty"`
}

// Summary A summary of a users recent data
type Summary struct {
	// Config Summary schema version and calculation configuration
	Config Config `json:"config,omitempty"`

	// Dates dates tracked for summary calculation
	Dates Dates          `json:"dates,omitempty"`
	Stats *Summary_Stats `json:"stats,omitempty"`

	// Type Field which contains a summary type string.
	Type                     SummaryTypeSchema `json:"type,omitempty"`
	UpdateWithoutChangeCount int               `json:"updateWithoutChangeCount"`

	// UserId String representation of a Tidepool User ID. Old style IDs are 10-digit strings consisting of only hexadeximcal digits. New style IDs are 36-digit [UUID v4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random))
	UserId *TidepoolUserId `json:"userId,omitempty"`
}

// Summary_Stats defines model for Summary.Stats.
type Summary_Stats struct {
	union json.RawMessage
}

// SummaryTypeSchema Field which contains a summary type string.
type SummaryTypeSchema = string

// TidepoolUserId String representation of a Tidepool User ID. Old style IDs are 10-digit strings consisting of only hexadeximcal digits. New style IDs are 36-digit [UUID v4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random))
type TidepoolUserId = string

// SummaryType Field which contains a summary type string.
type SummaryType = SummaryTypeSchema

// UserId String representation of a Tidepool User ID. Old style IDs are 10-digit strings consisting of only hexadeximcal digits. New style IDs are 36-digit [UUID v4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random))
type UserId = TidepoolUserId

// AsBGMBucketData returns the union data inside the Bucket_Data as a BGMBucketData
func (t Bucket_Data) AsBGMBucketData() (BGMBucketData, error) {
	var body BGMBucketData
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromBGMBucketData overwrites any union data inside the Bucket_Data as the provided BGMBucketData
func (t *Bucket_Data) FromBGMBucketData(v BGMBucketData) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeBGMBucketData performs a merge with any union data inside the Bucket_Data, using the provided BGMBucketData
func (t *Bucket_Data) MergeBGMBucketData(v BGMBucketData) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsCGMBucketData returns the union data inside the Bucket_Data as a CGMBucketData
func (t Bucket_Data) AsCGMBucketData() (CGMBucketData, error) {
	var body CGMBucketData
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromCGMBucketData overwrites any union data inside the Bucket_Data as the provided CGMBucketData
func (t *Bucket_Data) FromCGMBucketData(v CGMBucketData) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeCGMBucketData performs a merge with any union data inside the Bucket_Data, using the provided CGMBucketData
func (t *Bucket_Data) MergeCGMBucketData(v CGMBucketData) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsContinuousBucketData returns the union data inside the Bucket_Data as a ContinuousBucketData
func (t Bucket_Data) AsContinuousBucketData() (ContinuousBucketData, error) {
	var body ContinuousBucketData
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromContinuousBucketData overwrites any union data inside the Bucket_Data as the provided ContinuousBucketData
func (t *Bucket_Data) FromContinuousBucketData(v ContinuousBucketData) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeContinuousBucketData performs a merge with any union data inside the Bucket_Data, using the provided ContinuousBucketData
func (t *Bucket_Data) MergeContinuousBucketData(v ContinuousBucketData) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

func (t Bucket_Data) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *Bucket_Data) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// AsCGMStats returns the union data inside the Summary_Stats as a CGMStats
func (t Summary_Stats) AsCGMStats() (CGMStats, error) {
	var body CGMStats
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromCGMStats overwrites any union data inside the Summary_Stats as the provided CGMStats
func (t *Summary_Stats) FromCGMStats(v CGMStats) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeCGMStats performs a merge with any union data inside the Summary_Stats, using the provided CGMStats
func (t *Summary_Stats) MergeCGMStats(v CGMStats) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsBGMStats returns the union data inside the Summary_Stats as a BGMStats
func (t Summary_Stats) AsBGMStats() (BGMStats, error) {
	var body BGMStats
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromBGMStats overwrites any union data inside the Summary_Stats as the provided BGMStats
func (t *Summary_Stats) FromBGMStats(v BGMStats) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeBGMStats performs a merge with any union data inside the Summary_Stats, using the provided BGMStats
func (t *Summary_Stats) MergeBGMStats(v BGMStats) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsContinuousStats returns the union data inside the Summary_Stats as a ContinuousStats
func (t Summary_Stats) AsContinuousStats() (ContinuousStats, error) {
	var body ContinuousStats
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromContinuousStats overwrites any union data inside the Summary_Stats as the provided ContinuousStats
func (t *Summary_Stats) FromContinuousStats(v ContinuousStats) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeContinuousStats performs a merge with any union data inside the Summary_Stats, using the provided ContinuousStats
func (t *Summary_Stats) MergeContinuousStats(v ContinuousStats) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

func (t Summary_Stats) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *Summary_Stats) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}
