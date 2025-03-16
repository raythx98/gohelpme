package timehelper

import "time"

func TimePtr(t time.Time) *time.Time {
	return &t
}
