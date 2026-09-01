package validation

import (
	"fmt"
	"time"
)

func FormatWorkingHours(d time.Duration) string { //suppose d = 7 hours 35 min
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60                //divide the minutes %60
	return fmt.Sprintf("%02d:%02d", hours, minutes) //%02d -> no with atleast 2 digit
}
