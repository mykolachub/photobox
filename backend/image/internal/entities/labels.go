package entities

type Label struct {
	ID   string
	Name string
}

type MetadataLabel struct {
	ID         string
	MetadataID string
	LabelID    string
}
