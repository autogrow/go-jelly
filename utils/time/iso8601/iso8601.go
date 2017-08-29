package iso8601

import "time"

// Format of ISO 8601 time in golang (http://stackoverflow.com/a/42358468/5323316)
const Format = "2006-01-02T15:04:05Z07:00"

// Now returns the time in ISO 8601 format
func Now() string {
	return time.Now().Format(Format)
}

// Convert the given time to ISO 8601
func Convert(t time.Time) string {
	return t.Format(Format)
}

// Parse returns a time for a give ISO 8601 formated time
func Parse(t string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, t)
}
