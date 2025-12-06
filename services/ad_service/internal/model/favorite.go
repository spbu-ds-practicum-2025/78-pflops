package model

import "time"

type Favorite struct {
	UserID    string
	AdID      string
	CreatedAt time.Time
}
