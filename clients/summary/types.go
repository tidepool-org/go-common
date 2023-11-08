// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
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
	Date           *time.Time `json:"date,omitempty"`
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
	TotalGlucose *float64 `json:"totalGlucose,omitempty"`
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
	AverageGlucoseMmolDelta  *float64 `json:"averageGlucoseMmolDelta,omitempty"`
	HasAverageDailyRecords   *bool    `json:"hasAverageDailyRecords,omitempty"`
	HasAverageGlucoseMmol    *bool    `json:"hasAverageGlucoseMmol,omitempty"`
	HasTimeInAnyHighPercent  *bool    `json:"hasTimeInAnyHighPercent,omitempty"`
	HasTimeInAnyHighRecords  *bool    `json:"hasTimeInAnyHighRecords,omitempty"`
	HasTimeInAnyLowPercent   *bool    `json:"hasTimeInAnyLowPercent,omitempty"`
	HasTimeInAnyLowRecords   *bool    `json:"hasTimeInAnyLowRecords,omitempty"`
	HasTimeInHighPercent     *bool    `json:"hasTimeInHighPercent,omitempty"`
	HasTimeInHighRecords     *bool    `json:"hasTimeInHighRecords,omitempty"`
	HasTimeInLowPercent      *bool    `json:"hasTimeInLowPercent,omitempty"`
	HasTimeInLowRecords      *bool    `json:"hasTimeInLowRecords,omitempty"`
	HasTimeInTargetPercent   *bool    `json:"hasTimeInTargetPercent,omitempty"`
	HasTimeInTargetRecords   *bool    `json:"hasTimeInTargetRecords,omitempty"`
	HasTimeInVeryHighPercent *bool    `json:"hasTimeInVeryHighPercent,omitempty"`
	HasTimeInVeryHighRecords *bool    `json:"hasTimeInVeryHighRecords,omitempty"`
	HasTimeInVeryLowPercent  *bool    `json:"hasTimeInVeryLowPercent,omitempty"`
	HasTimeInVeryLowRecords  *bool    `json:"hasTimeInVeryLowRecords,omitempty"`
	HasTotalRecords          *bool    `json:"hasTotalRecords,omitempty"`

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
	Buckets *[]Bucket `json:"buckets,omitempty"`

	// OffsetPeriods A map to each supported BGM summary period
	OffsetPeriods *BGMPeriods `json:"offsetPeriods,omitempty"`

	// Periods A map to each supported BGM summary period
	Periods *BGMPeriods `json:"periods,omitempty"`

	// TotalHours Total hours represented in the hourly stats
	TotalHours *int `json:"totalHours,omitempty"`
}

// Bucket bucket containing an hour of bgm or cgm aggregations
type Bucket struct {
	Data           *Bucket_Data `json:"data,omitempty"`
	Date           *time.Time   `json:"date,omitempty"`
	LastRecordTime *time.Time   `json:"lastRecordTime,omitempty"`
}

// Bucket_Data defines model for Bucket.Data.
type Bucket_Data struct {
	union json.RawMessage
}

// CGMBucketData Series of counters which represent one hour of a users data
type CGMBucketData struct {
	// TimeCGMUseMinutes Counter of minutes using a cgm
	TimeCGMUseMinutes *int `json:"timeCGMUseMinutes,omitempty"`

	// TimeCGMUseRecords Counter of records wearing a cgm
	TimeCGMUseRecords *int `json:"timeCGMUseRecords,omitempty"`

	// TimeInHighMinutes Counter of minutes spent in high glucose range
	TimeInHighMinutes *int `json:"timeInHighMinutes,omitempty"`

	// TimeInHighRecords Counter of records in high glucose range
	TimeInHighRecords *int `json:"timeInHighRecords,omitempty"`

	// TimeInLowMinutes Counter of minutes spent in low glucose range
	TimeInLowMinutes *int `json:"timeInLowMinutes,omitempty"`

	// TimeInLowRecords Counter of records in low glucose range
	TimeInLowRecords *int `json:"timeInLowRecords,omitempty"`

	// TimeInTargetMinutes Counter of minutes spent in target glucose range
	TimeInTargetMinutes *int `json:"timeInTargetMinutes,omitempty"`

	// TimeInTargetRecords Counter of records in target glucose range
	TimeInTargetRecords *int `json:"timeInTargetRecords,omitempty"`

	// TimeInVeryHighMinutes Counter of minutes spent in very high glucose range
	TimeInVeryHighMinutes *int `json:"timeInVeryHighMinutes,omitempty"`

	// TimeInVeryHighRecords Counter of records in very high glucose range
	TimeInVeryHighRecords *int `json:"timeInVeryHighRecords,omitempty"`

	// TimeInVeryLowMinutes Counter of minutes spent in very low glucose range
	TimeInVeryLowMinutes *int `json:"timeInVeryLowMinutes,omitempty"`

	// TimeInVeryLowRecords Counter of records in very low glucose range
	TimeInVeryLowRecords *int `json:"timeInVeryLowRecords,omitempty"`

	// TotalGlucose Total value of all glucose records
	TotalGlucose *float64 `json:"totalGlucose,omitempty"`
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
	AverageGlucoseMmolDelta *float64 `json:"averageGlucoseMmolDelta,omitempty"`

	// GlucoseManagementIndicator A derived value which emulates A1C
	GlucoseManagementIndicator *float64 `json:"glucoseManagementIndicator,omitempty"`

	// GlucoseManagementIndicatorDelta Difference between the glucoseManagementIndicator in this period and the other offset version
	GlucoseManagementIndicatorDelta *float64 `json:"glucoseManagementIndicatorDelta,omitempty"`
	HasAverageDailyRecords          *bool    `json:"hasAverageDailyRecords,omitempty"`
	HasAverageGlucoseMmol           *bool    `json:"hasAverageGlucoseMmol,omitempty"`
	HasGlucoseManagementIndicator   *bool    `json:"hasGlucoseManagementIndicator,omitempty"`
	HasTimeCGMUseMinutes            *bool    `json:"hasTimeCGMUseMinutes,omitempty"`
	HasTimeCGMUsePercent            *bool    `json:"hasTimeCGMUsePercent,omitempty"`
	HasTimeCGMUseRecords            *bool    `json:"hasTimeCGMUseRecords,omitempty"`
	HasTimeInAnyHighMinutes         *bool    `json:"hasTimeInAnyHighMinutes,omitempty"`
	HasTimeInAnyHighPercent         *bool    `json:"hasTimeInAnyHighPercent,omitempty"`
	HasTimeInAnyHighRecords         *bool    `json:"hasTimeInAnyHighRecords,omitempty"`
	HasTimeInAnyLowMinutes          *bool    `json:"hasTimeInAnyLowMinutes,omitempty"`
	HasTimeInAnyLowPercent          *bool    `json:"hasTimeInAnyLowPercent,omitempty"`
	HasTimeInAnyLowRecords          *bool    `json:"hasTimeInAnyLowRecords,omitempty"`
	HasTimeInHighMinutes            *bool    `json:"hasTimeInHighMinutes,omitempty"`
	HasTimeInHighPercent            *bool    `json:"hasTimeInHighPercent,omitempty"`
	HasTimeInHighRecords            *bool    `json:"hasTimeInHighRecords,omitempty"`
	HasTimeInLowMinutes             *bool    `json:"hasTimeInLowMinutes,omitempty"`
	HasTimeInLowPercent             *bool    `json:"hasTimeInLowPercent,omitempty"`
	HasTimeInLowRecords             *bool    `json:"hasTimeInLowRecords,omitempty"`
	HasTimeInTargetMinutes          *bool    `json:"hasTimeInTargetMinutes,omitempty"`
	HasTimeInTargetPercent          *bool    `json:"hasTimeInTargetPercent,omitempty"`
	HasTimeInTargetRecords          *bool    `json:"hasTimeInTargetRecords,omitempty"`
	HasTimeInVeryHighMinutes        *bool    `json:"hasTimeInVeryHighMinutes,omitempty"`
	HasTimeInVeryHighPercent        *bool    `json:"hasTimeInVeryHighPercent,omitempty"`
	HasTimeInVeryHighRecords        *bool    `json:"hasTimeInVeryHighRecords,omitempty"`
	HasTimeInVeryLowMinutes         *bool    `json:"hasTimeInVeryLowMinutes,omitempty"`
	HasTimeInVeryLowPercent         *bool    `json:"hasTimeInVeryLowPercent,omitempty"`
	HasTimeInVeryLowRecords         *bool    `json:"hasTimeInVeryLowRecords,omitempty"`
	HasTotalRecords                 *bool    `json:"hasTotalRecords,omitempty"`

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
	Buckets *[]Bucket `json:"buckets,omitempty"`

	// OffsetPeriods A map to each supported CGM summary period
	OffsetPeriods *CGMPeriods `json:"offsetPeriods,omitempty"`

	// Periods A map to each supported CGM summary period
	Periods *CGMPeriods `json:"periods,omitempty"`

	// TotalHours Total hours represented in the hourly stats
	TotalHours *int `json:"totalHours,omitempty"`
}

// Config Summary schema version and calculation configuration
type Config struct {
	// HighGlucoseThreshold Threshold used for determining if a value is high
	HighGlucoseThreshold *float64 `json:"highGlucoseThreshold,omitempty"`

	// LowGlucoseThreshold Threshold used for determining if a value is low
	LowGlucoseThreshold *float64 `json:"lowGlucoseThreshold,omitempty"`

	// SchemaVersion Summary schema version
	SchemaVersion *int `json:"schemaVersion,omitempty"`

	// VeryHighGlucoseThreshold Threshold used for determining if a value is very high
	VeryHighGlucoseThreshold *float64 `json:"veryHighGlucoseThreshold,omitempty"`

	// VeryLowGlucoseThreshold Threshold used for determining if a value is very low
	VeryLowGlucoseThreshold *float64 `json:"veryLowGlucoseThreshold,omitempty"`
}

// Dates dates tracked for summary calculation
type Dates struct {
	// FirstData Date of the first included value
	FirstData         *time.Time `json:"firstData,omitempty"`
	HasFirstData      *bool      `json:"hasFirstData,omitempty"`
	HasLastData       *bool      `json:"hasLastData,omitempty"`
	HasLastUploadDate *bool      `json:"hasLastUploadDate,omitempty"`
	HasOutdatedSince  *bool      `json:"hasOutdatedSince,omitempty"`

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
	Config *Config `json:"config,omitempty"`

	// Dates dates tracked for summary calculation
	Dates *Dates         `json:"dates,omitempty"`
	Stats *Summary_Stats `json:"stats,omitempty"`

	// Type Field which contains a summary type string.
	Type *SummaryTypeSchema `json:"type,omitempty"`

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

	merged, err := runtime.JsonMerge(t.union, b)
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

	merged, err := runtime.JsonMerge(t.union, b)
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

	merged, err := runtime.JsonMerge(t.union, b)
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

	merged, err := runtime.JsonMerge(t.union, b)
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
