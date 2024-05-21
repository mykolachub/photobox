package entity

import "time"

type Meta struct {
	ID               string
	UserID           string
	FileLocation     string
	FileName         string
	FileSize         int
	FileExt          string
	FileLastModified time.Time
	CreatedAt        time.Time
}
