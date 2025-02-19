package xtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	fx "github.com/org/2112-space-lab/org/app-service/pkg/option"
)

// ToUtcTime converts *time.Time to xtime.UtcTime.
func ToUtcTime(t *time.Time) fx.Option[UtcTime] {
	if t == nil {
		return fx.NewEmptyOption[UtcTime]()
	}
	utcTime := NewUtcTimeIgnoreZone(*t) // Ensure the time is standardized to UTC.
	return fx.NewValueOption(utcTime)
}

// toTimePointer converts fx.Option[xtime.UtcTime] to *time.Time.
func ToTimePointer(opt fx.Option[UtcTime]) *time.Time {
	if !opt.HasValue {
		return nil
	}
	utcTime := opt.Value.Inner()
	return &utcTime
}

// UtcTime is a safe wrapper around Time intended to always hold UTC time with timezone correctly applied
type UtcTime struct {
	inner time.Time
}

// Inner returns inner time value in a standard format
func (t UtcTime) Inner() time.Time {
	return StandardizeTime(t.inner)
}

// Add returns a new utc time with added duration
func (t UtcTime) Add(d time.Duration) UtcTime {
	newTime := StandardizeTime(t.inner).Add(d)
	return NewUtcTimeIgnoreZone(newTime)
}

// type PotentialUtcDateTimeISO string
// type PotentialUtcDateTimeARC string

// DateTimeFormat is a type for formatting dates and times according to a language.
type DateTimeFormat string

// IsoFormat is the international standard ISO 8601 specifies the numerical representation of date and time
// - based on the Gregorian calendar and the 24-hour time system, respectively.
const IsoFormat DateTimeFormat = time.RFC3339

// const IsoUTCFormat DateTimeFormat = "2006-01-02T15:04:05.000Z"
// const IpacsFormat DateTimeFormat = "20060102150405.000"
// const PassCalcDopplerFileNameFormat DateTimeFormat = "20060102150405"
// const PassCalcApiInputFormat DateTimeFormat = "2006/01/02 15:04:05.000"

// FromString parse utc string to UtcTime
func FromString(utcTime DateTimeFormat) (UtcTime, error) {
	if utcTime == "" {
		return UtcTime{}, nil
	}

	timeUtc, err := time.Parse(string(IsoFormat), string(utcTime))
	if err != nil {
		return UtcTime{}, err
	}

	return UtcTime{
		inner: StandardizeTime(timeUtc),
	}, nil
}

// String translates UtcTime to String
func String(utcTime UtcTime, format DateTimeFormat) string {
	return string(utcTime.FormatStr(format))
}

// StandardizeTime ensure that time will remain same across de/serialization operations where we can loose precision
func StandardizeTime(t time.Time) time.Time {
	return t.
		Round(1 * time.Microsecond). // postgresql timestamps are limited to ms
		In(time.UTC)                 // to ensure equality before/after serialize
}

// UtcNow returns Now time in UTC
func UtcNow() UtcTime {
	return UtcTime{
		inner: StandardizeTime(time.Now().UTC()),
	}
}

// NewUtcTimeIgnoreZone ensures the given time in UTC in case it was parsed assumed local because of zone info missing
func NewUtcTimeIgnoreZone(t time.Time) UtcTime {
	_, offsetSec := t.Zone()
	t = StandardizeTime(t)
	asUTC := t.Add(time.Duration(offsetSec) * time.Second)
	return UtcTime{
		inner: asUTC,
	}
}

// UnixInLocation helper to convert unix timestamp into time of given Location.
// it solves the issue of time.Unix assuming input timestamp being a Local time
func UnixInLocation(sec int64, nsec int64, wantedLocation *time.Location) time.Time {
	// unixLocalTime := time.Unix(sec, nsec)
	// _, offsetSec := unixLocalTime.Zone()
	// asUTC := unixLocalTime.Add(time.Duration(offsetSec) * time.Second).In(wantedLocation)
	asUTC := time.Unix(sec, nsec).In(wantedLocation)
	return asUTC
}

// UnixUTC converts unix timestamp into a UTC time.
// ! nanosecond will be truncated to milliseconds because of postgreSQL timestamp
func UnixUTC(sec int64, nsec int64) UtcTime {
	utcTime := UnixInLocation(sec, nsec, time.UTC)
	standardized := StandardizeTime(utcTime)
	return UtcTime{
		inner: standardized,
	}
}

// Latest returns the latest time from passed input
func Latest(anchorTime UtcTime, otherTimes ...UtcTime) UtcTime {
	latest := anchorTime
	for _, t := range otherTimes {
		if t.inner.After(latest.inner) {
			latest = t
		}
	}
	return latest
}

// Earliest returns the earliest time from passed input
func Earliest(anchorTime UtcTime, otherTimes ...UtcTime) UtcTime {
	earliest := anchorTime
	for _, t := range otherTimes {
		if t.inner.Before(earliest.inner) {
			earliest = t
		}
	}
	return earliest
}

// FormatStr returns a formatted UtcTIme time
func (t UtcTime) FormatStr(format DateTimeFormat) string {
	return t.inner.Format(string(format))
}

// Format implements Formatter interface
func (t UtcTime) Format(w fmt.State, v rune) {
	_, _ = fmt.Fprint(w, t.inner.Format(string(time.RFC3339Nano)))
}

// DurationUntil calculate the time span duration until given time
func (t UtcTime) DurationUntil(until UtcTime) time.Duration {
	d := t.inner.Sub(until.inner)
	return -d
}

// After check if given time is after
func (t UtcTime) After(other UtcTime) bool {
	isAfter := t.inner.After(other.inner)
	return isAfter
}

// Before check if given time is before
func (t UtcTime) Before(other UtcTime) bool {
	isBefore := t.inner.Before(other.inner)
	return isBefore
}

// IsZero reports whether t represents the zero time instant, January 1, year 1, 00:00:00 UTC.
func (t UtcTime) IsZero() bool {
	return t.inner.IsZero()
}

// MarshalJSON override for time value to json
func (t UtcTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.FormatStr(time.RFC3339Nano))
}

// UnmarshalJSON override for json value to time
func (t *UtcTime) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		t.inner = time.Time{}
		return nil
	}

	var valueStr string
	if err := json.Unmarshal(data, &valueStr); err != nil {
		return err
	}
	tt, err := time.Parse(string(time.RFC3339Nano), valueStr)
	if err != nil {
		return err
	}
	t.inner = NewUtcTimeIgnoreZone(tt).inner
	return nil
}

// String translates UtcTime to String
func (t UtcTime) String() string {
	return string(t.FormatStr(time.RFC3339Nano))
}

// GoString translates UtcTime to String
func (t UtcTime) GoString() string {
	return string(t.FormatStr(time.RFC3339Nano))
}

// ConvertToUtcTime *time.Time to xtime.UtcTime before using fx.AsOption
func ConvertToUtcTime(t *time.Time) fx.Option[UtcTime] {
	if t == nil {
		return fx.NewEmptyOption[UtcTime]()
	}
	return fx.NewValueOption(NewUtcTimeIgnoreZone(*t))
}

// ConvertToTimePtr converts fx.Option[xtime.UtcTime] to *time.Time
func ConvertToTimePtr(opt fx.Option[UtcTime]) *time.Time {
	if !opt.HasValue {
		return nil
	}
	t := time.Time(opt.Value.Inner())
	return &t
}
