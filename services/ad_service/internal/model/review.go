package model

import "time"

type Review struct {
	ID         string
	AdID       string
	ReviewerID string
	Rating     int // 1..5
	Comment    *string
	CreatedAt  time.Time
}
