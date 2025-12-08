package model

import "time"

// Ad domain model
// NOTE: sellerRatingCached может быть пустым (nil) если еще не агрегирован рейтинг

type Ad struct {
	ID                 string
	AuthorID           string
	Title              string
	Description        string
	Price              int64
	CategoryID         string
	Condition          string // NEW, USED, REFURBISHED
	Status             string // ACTIVE, ARCHIVED
	SellerRatingCached *float64
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Images             []AdImage
}
