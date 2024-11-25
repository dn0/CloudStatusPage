package s3

import (
	"time"
)

const (
	waiterMinDelay = 100 * time.Millisecond
	waiterMaxDelay = 1000 * time.Millisecond
)
