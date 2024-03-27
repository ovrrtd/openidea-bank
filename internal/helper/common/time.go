package common

import "time"

// unix milli to time ISO 8601 format
func UnixMilliToISO8601(unixMilli int64) string {
	return time.UnixMilli(unixMilli).Format(time.RFC3339)
}
