package entity

import "time"

type Meta struct {
	ID            string
	UserID        string
	FileLocation  string
	FileName      string
	FileSize      int
	FileType      string
	FileCreatedAd time.Time
	CreatedAd     time.Time
	UpdatedAt     time.Time
}
