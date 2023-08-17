package util

import (
	"fmt"
	"time"
)

// Formats seconds to HH:MM
func SecondsToHHMM(seconds int64) string {
	duration := time.Second * time.Duration(seconds)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	return fmt.Sprintf("%02d:%02d", hours, minutes)
}
