package model

import "time"

type Launch struct {
	ID        int
	UserID    int
	ContestID int
	Date      time.Time
	Code      string
	TestsResult
}
