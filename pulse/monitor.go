package pulse

import "time"

type Monitor struct {
	Status map[string]time.Time
}
