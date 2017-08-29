package time

import "time"

// Midday gives the epoch of the midday for the given day
func Midday(epoch float64) float64 {
	s := time.Unix(int64(epoch), 0)
	midday := time.Date(s.Year(), s.Month(), s.Day(), 12, 0, 0, 0, s.Location())
	return float64(midday.Unix())
}
