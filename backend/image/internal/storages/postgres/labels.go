package postgres

import (
	"database/sql"
	"photobox/image/internal/entities"

	_ "github.com/lib/pq"
)

type LabelsRepo struct {
	db *sql.DB
}

func InitLabelsRepo(db *sql.DB) *LabelsRepo {
	return &LabelsRepo{db}
}

func (r *LabelsRepo) CreateLabel(data entities.Label) (entities.Label, error) {
	var label entities.Label
	tx, err := r.db.Begin()
	if err != nil {
		return entities.Label{}, err
	}
	defer tx.Rollback()

	// Check if the label already exists
	query := `SELECT id, value FROM labels WHERE value = $1`
	err = tx.QueryRow(query, data.Name).Scan(&label.ID, &label.Name)
	if err != sql.ErrNoRows {
		tx.Commit()
		return label, nil
	}

	// If the label doesn't exist, insert a new one
	query = `INSERT INTO labels (value) VALUES ($1) RETURNING id, value`
	row := tx.QueryRow(query, data.Name)
	if row == nil {
		tx.Rollback()
		return entities.Label{}, sql.ErrNoRows
	}
	err = row.Scan(&label.ID, &label.Name)
	if err != nil {
		tx.Rollback()
		return entities.Label{}, err
	}

	tx.Commit()
	return label, nil
}

func (r *LabelsRepo) CreateMetadataLabel(metadata_id, label_id string) (entities.MetadataLabel, error) {
	metadataLabel := entities.MetadataLabel{}

	query := `
	INSERT INTO
		metadata_labels (metadata_id, label_id)
	VALUES
		($1, $2)
	RETURNING
		id, metadata_id, label_id
	`

	err := r.db.QueryRow(query, metadata_id, label_id).Scan(
		&metadataLabel.ID,
		&metadataLabel.MetadataID,
		&metadataLabel.LabelID,
	)
	if err != nil {
		return entities.MetadataLabel{}, err
	}

	return metadataLabel, nil
}

func (r *LabelsRepo) GetMetadataLabelByMetalId(metadata_id string) ([]entities.MetadataLabel, error) {
	metadataLabels := []entities.MetadataLabel{}

	query := `
	SELECT
		id, metadata_id, label_id
	FROM
		metadata_labels
	WHERE
		metadata_id = $1
	`

	rows, err := r.db.Query(query, metadata_id)
	if err != nil {
		return []entities.MetadataLabel{}, err
	}

	for rows.Next() {
		metadataLabel := entities.MetadataLabel{}
		err := rows.Scan(
			&metadataLabel.ID,
			&metadataLabel.MetadataID,
			&metadataLabel.LabelID,
		)
		if err != nil {
			return []entities.MetadataLabel{}, err
		}
		metadataLabels = append(metadataLabels, metadataLabel)
	}

	return metadataLabels, nil
}
