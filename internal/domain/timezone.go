package domain

import (
	"time"
)

// ISTLocation is the IST timezone (Asia/Kolkata, UTC+5:30)
// Shared across all domain services
var ISTLocation *time.Location

func init() {
	var err error
	ISTLocation, err = time.LoadLocation("Asia/Kolkata")
	if err != nil {
		// Fallback to UTC+5:30 if location loading fails
		ISTLocation = time.FixedZone("IST", 5*3600+30*60)
	}
}
