package time

import (
	"errors"
	"time"
)

var (
	// Weekdays ...
	Weekdays = [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
)

const (
	// NanoSecond ...
	NanoSecond Duration = 1
	// MicroSecond ...
	MicroSecond Duration = 1000 * NanoSecond
	// MilliSecond ...
	MilliSecond Duration = 1000 * MicroSecond
	// Second ...
	Second Duration = 1000 * MilliSecond
	// Minute ...
	Minute Duration = 60 * Second
	// Hour ...
	Hour Duration = 60 * Minute
	// Day ...
	Day Duration = 24 * Hour
	// Week ...
	Week Duration = 7 * Day
	// CassandraFormat ...
	CassandraFormat string = "2006-01-02 15:04:05.000Z0700"
)

// Duration ...
type Duration time.Duration

// Location ...
type Location struct {
	time.Location
}

// Month ...
type Month time.Month

// ParseError ...
type ParseError struct {
	time.ParseError
}

// Ticker ...
type Ticker struct {
	time.Ticker
}

// Time ...
type Time struct {
	time.Time
}

// Timer ...
type Timer struct {
	time.Timer
}

// Add ...
func (t Time) Add(d Duration) Time {
	return Time{t.Time.Add(time.Duration(d))}
}

// AddDate ...
func (t Time) AddDate(years, months, days int) Time {
	return Time{t.Time.AddDate(years, months, days)}
}

// Equals ...
func (t Time) Equals(time Time) bool {
	return t.Unix() == time.Unix()
}

// In ...
func (t Time) In(loc *Location) Time {
	return Time{t.Time.In(&((*loc).Location))}
}

// Location ...
func (t Time) Location() *Location {
	return &Location{*(t.Time.Location())}
}

// MarshalJSON ...
func (t Time) MarshalJSON() ([]byte, error) {
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(CassandraFormat)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, CassandraFormat)
	b = append(b, '"')
	return b, nil
}

// Month ...
func (t Time) Month() Month {
	return Month(t.Time.Month())
}

// Round ...
func (t Time) Round(d Duration) Time {
	return Time{t.Time.Round(time.Duration(d))}
}

// Sub ...
func (t Time) Sub(u Time) Duration {
	return Duration(t.Time.Sub(u.Time))
}

// Truncate ...
func (t Time) Truncate(d Duration) Time {
	return Time{t.Time.Truncate(time.Duration(d))}
}

// UTC ...
func (t Time) UTC() Time {
	return Time{t.Time.UTC()}
}

// UnmarshalJSON ...
func (t *Time) UnmarshalJSON(data []byte) error {
	// Fractional seconds are handled implicitly by Parse.
	goTime, err := time.Parse(`2006-01-02 15:04:05.000Z0700`, string(data))
	*t = Time{goTime}
	return err
}

// Weekday ...
func (t Time) Weekday() Weekday {
	return Weekday(t.Time.Weekday())
}

// Weekday ...
type Weekday time.Weekday

// After ...
func After(d Duration) <-chan Time {
	tempChan := time.After(time.Duration(d))
	outChan := make(chan Time)
	for i := 0; i < len(tempChan); i++ {
		tempVal := <-tempChan
		outChan <- Time{tempVal}
	}
	return outChan
}

// Sleep ...
func Sleep(d Duration) {
	time.Sleep(time.Duration(d))
}

// Tick ...
func Tick(d Duration) <-chan Time {
	tempChan := time.Tick(time.Duration(d))
	outChan := make(chan Time)
	for i := 0; i < len(tempChan); i++ {
		tempVal := <-tempChan
		outChan <- Time{tempVal}
	}
	return outChan
}

// ParseDuration ...
func ParseDuration(s string) (Duration, error) {
	dur, err := time.ParseDuration(s)
	return Duration(dur), err
}

// Since ...
func Since(t Time) Duration {
	return Duration(time.Since(t.Time))
}

// FixedZone ...
func FixedZone(name string, offset int) *Location {
	loc := time.FixedZone(name, offset)
	return &(Location{*loc})
}

// LoadLocation ...
func LoadLocation(name string) (*Location, error) {
	loc, err := time.LoadLocation(name)
	return &(Location{*loc}), err
}

// NewTicker ...
func NewTicker(d Duration) *Ticker {
	tick := time.NewTicker(time.Duration(d))
	return &(Ticker{*tick})
}

// Parse ...
func Parse(layout, value string) (Time, error) {
	t, err := time.Parse(layout, value)
	return Time{t}, err
}

// Date ...
func Date(year int, month Month, day, hour, min, sec, nsec int, loc *Location) Time {
	return Time{time.Date(year, time.Month(month), day, hour, min, sec, nsec, &((*loc).Location))}
}

// ParseInLocation ...
func ParseInLocation(layout, value string, loc *Location) (Time, error) {
	t, err := time.ParseInLocation(layout, value, &((*loc).Location))
	return Time{t}, err
}

// Unix ...
func Unix(sec, nsec int64) Time {
	return Time{time.Unix(sec, nsec)}
}

// AfterFunc ...
func AfterFunc(d Duration, f func()) *Timer {
	return &Timer{*(time.AfterFunc(time.Duration(d), f))}
}

// NewTimer ...
func NewTimer(d Duration) *Timer {
	return &Timer{*(time.NewTimer(time.Duration(d)))}
}

// Now ...
func Now() Time {
	return Time{time.Now()}
}
