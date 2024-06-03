package constants

import "time"

const (
	EXPIRE_AFTER   time.Duration = 15 * time.Minute
	CLEAN_INTERVAL time.Duration = 30 * time.Second
)
