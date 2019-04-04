package util

import (
	"sort"
	"time"
)

func FindNewest(dates []time.Time) time.Time {
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].After(dates[j])
	})

	return dates[0]
}
