package entity

import "time"

type Meta struct {
	ID               string
	UserID           string
	FileLocation     string
	FileName         string
	FileSize         int
	FileExt          string
	FileWidth        int
	FileHeight       int
	FileLastModified time.Time
	Labels           []Label
	CreatedAt        time.Time
}

type Label struct {
	ID              string
	Value           string
	MetadataLabelID string
}
