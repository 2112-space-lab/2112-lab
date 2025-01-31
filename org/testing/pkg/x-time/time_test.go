package xtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFromString(t *testing.T) {
	tests := []struct {
		name     string
		inputStr DateTimeFormat
		expected UtcTime
	}{
		{
			name:     "valid utc time from string",
			inputStr: DateTimeFormat("2024-02-21T01:10:30Z"),
			expected: UtcTime{
				inner: StandardizeTime(time.Date(2024, time.Month(2), 21, 1, 10, 30, 0, time.FixedZone("", 0))),
			},
		},
		{
			name:     "invalid utc time from string",
			inputStr: DateTimeFormat("1"),
			expected: UtcTime{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := FromString(tt.inputStr)
			empty := UtcTime{}
			if tt.expected == empty {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tt.expected, res)
			}
		})
	}
}

func TestTimeAdd(t *testing.T) {
	input := NewUtcTimeIgnoreZone(StandardizeTime(time.Date(2024, time.Month(2), 21, 1, 10, 15, 0, time.FixedZone("", 0))))
	expected := NewUtcTimeIgnoreZone(StandardizeTime(time.Date(2024, time.Month(2), 21, 1, 10, 45, 0, time.FixedZone("", 0))))
	require.Equal(t, expected, input.Add(30*time.Second))
}

func TestToString(t *testing.T) {
	defaultTime := DateTimeFormat("0001-01-01T00:00:00Z")
	tests := []struct {
		name     string
		inputStr UtcTime
		expected string
	}{
		{
			name:     "invalid utc time to string",
			inputStr: UtcTime{},
			expected: "0001-01-01T00:00:00Z",
		},
		{
			name: "valid utc time to string",
			inputStr: UtcTime{
				inner: time.Date(2024, time.Month(2), 21, 1, 10, 30, 0, time.FixedZone("", 0)),
			},
			expected: "2024-02-21T01:10:30Z",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.inputStr.FormatStr(DateTimeFormat(IsoFormat))
			if tt.expected == "" {
				require.Equal(t, defaultTime, res)
			} else {
				require.NotEmpty(t, res)
				require.Equal(t, tt.expected, res)
			}
		})
	}
}

func TestDurationUntil(t *testing.T) {
	t1 := UtcNow()
	dur := 200 * time.Millisecond
	t2 := NewUtcTimeIgnoreZone(t1.inner.Add(dur))

	res1 := t1.DurationUntil(t2)
	assert.Equal(t, dur, res1)

	t3 := NewUtcTimeIgnoreZone(t1.inner.Add(res1))
	assert.Equal(t, t2, t3)
}

func TestXxx(t *testing.T) {
	x := UnixUTC(1685089549, 0)
	y := time.Unix(1685089549, 0)
	require.Equal(t, x.inner.Unix(), y.Unix())

	n := time.Now()
	u := n.UTC()
	z := timestamppb.Timestamp{
		Seconds: 1685089549,
	}
	zz := time.Unix(z.GetSeconds(), 0)

	require.Equal(t, zz.UTC().UTC().Unix(), y.Unix())

	require.Equal(t, n.Unix(), u.Unix())
}

func TestLatest(t *testing.T) {
	// // f := "2006-01-02T15:04:05Z07:00"
	// // t1, _ := time.ParseInLocation(f, "2014-07-16T22:55:46Z", time.UTC)
	// // t2, _ := time.ParseInLocation(f, "2014-07-16T22:55:46+02:00", time.Local)
	// t1 := time.Unix(1405544146, 0).In(time.FixedZone("", 0))
	// // t1 := time.Date(0, 0, 0, 0, 0, 1405544146, 0, time.UTC)
	// // t1.
	// t2 := time.Unix(1405544146, 0).UTC().In(time.FixedZone("", 0))
	// a, b := t2.Zone()
	// // t1 = t1.Add(time.Second * time.Duration(b))
	// assert.NotNil(t, a)
	// assert.NotNil(t, b)

	// // x1 := time.Parse("2006-01-02T15:04:05Z07:00")
	// assert.Equal(t, t1, t2)

	tests := []struct {
		name        string
		inputAnchor UtcTime
		inputOthers []UtcTime
		expected    UtcTime
	}{ // cases
		{
			name:        "single input is latest",
			inputAnchor: UnixUTC(1405544146, 0),
			inputOthers: []UtcTime{},
			expected:    UnixUTC(1405544146, 0),
		},
		{
			name:        "anchor is latest",
			inputAnchor: UnixUTC(1405544146, 0),
			inputOthers: []UtcTime{
				UnixUTC(1405544140, 0),
			},
			expected: UnixUTC(1405544146, 0),
		},
		{
			name:        "first of other is latest",
			inputAnchor: UnixUTC(1405544146, 0),
			inputOthers: []UtcTime{
				UnixUTC(1405544149, 0),
				UnixUTC(1405544140, 0),
			},
			expected: UnixUTC(1405544149, 0),
		},
		{
			name:        "second of other is latest",
			inputAnchor: UnixUTC(1405544146, 0),
			inputOthers: []UtcTime{
				UnixUTC(1405544140, 0),
				UnixUTC(1405544149, 0),
				UnixUTC(1405544147, 0),
			},
			expected: UnixUTC(1405544149, 0),
		},
		{
			name:        "last of other is latest",
			inputAnchor: UnixUTC(1405544146, 0),
			inputOthers: []UtcTime{
				UnixUTC(1405544140, 0),
				UnixUTC(1405544147, 0),
				UnixUTC(1405544149, 0),
			},
			expected: UnixUTC(1405544149, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Latest(tt.inputAnchor, tt.inputOthers...)
			require.Equal(t, tt.expected, res)
		})
	}
}

func TestEarliest(t *testing.T) {
	tests := []struct {
		name        string
		inputAnchor UtcTime
		inputOthers []UtcTime
		expected    UtcTime
	}{ // cases
		{
			name:        "single input is earliest",
			inputAnchor: UnixUTC(1405544146, 0),
			inputOthers: []UtcTime{},
			expected:    UnixUTC(1405544146, 0),
		},
		{
			name:        "anchor is earliest",
			inputAnchor: UnixUTC(1405544140, 0),
			inputOthers: []UtcTime{
				UnixUTC(1405544146, 0),
			},
			expected: UnixUTC(1405544140, 0),
		},
		{
			name:        "first of other is earliest",
			inputAnchor: UnixUTC(1405544146, 0),
			inputOthers: []UtcTime{
				UnixUTC(1405544140, 0),
				UnixUTC(1405544149, 0),
			},
			expected: UnixUTC(1405544140, 0),
		},
		{
			name:        "second of other is earliest",
			inputAnchor: UnixUTC(1405544146, 0),
			inputOthers: []UtcTime{
				UnixUTC(1405544149, 0),
				UnixUTC(1405544140, 0),
				UnixUTC(1405544147, 0),
			},
			expected: UnixUTC(1405544140, 0),
		},
		{
			name:        "last of other is earliest",
			inputAnchor: UnixUTC(1405544145, 0),
			inputOthers: []UtcTime{
				UnixUTC(1405544142, 0),
				UtcNow(),
				UnixUTC(1405544140, 0),
			},
			expected: UnixUTC(1405544140, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Earliest(tt.inputAnchor, tt.inputOthers...)
			require.Equal(t, tt.expected, res)
		})
	}
}
