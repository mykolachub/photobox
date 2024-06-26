package postgres

import (
	"database/sql"
	"photobox-meta/internal/entity"
	"time"

	_ "github.com/lib/pq"
)

type MetaRepo struct {
	db *sql.DB
}

func InitMetaRepo(db *sql.DB) *MetaRepo {
	return &MetaRepo{db: db}
}

func (r *MetaRepo) CreateMeta(data entity.Meta) (entity.Meta, error) {
	meta := entity.Meta{}
	now := time.Now().UTC()

	query := `
	INSERT INTO
		metadata(user_id, file_location, file_name, file_size, file_ext, file_last_modified, created_at, file_width, file_height)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING
		id, user_id,
		file_location, file_name,
		file_size, file_ext, file_last_modified,
		file_width, file_height,
		created_at`

	rows := r.db.QueryRow(query, data.UserID, data.FileLocation, data.FileName, data.FileSize, data.FileExt, data.FileLastModified, now, data.FileWidth, data.FileHeight)
	err := rows.Scan(
		&meta.ID,
		&meta.UserID,
		&meta.FileLocation,
		&meta.FileName,
		&meta.FileSize,
		&meta.FileExt,
		&meta.FileLastModified,
		&meta.FileWidth,
		&meta.FileHeight,
		&meta.CreatedAt,
	)
	if err != nil {
		return entity.Meta{}, err
	}

	return meta, nil
}

func (r *MetaRepo) GetMetaById(id string) (entity.Meta, error) {
	meta := entity.Meta{}

	query := `
	SELECT
		id, user_id,
		file_location, file_name,
		file_size, file_ext, file_last_modified,
		file_width, file_height,
		created_at
	FROM
		metadata
	WHERE
		id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&meta.ID,
		&meta.UserID,
		&meta.FileLocation,
		&meta.FileName,
		&meta.FileSize,
		&meta.FileExt,
		&meta.FileLastModified,
		&meta.FileWidth,
		&meta.FileHeight,
		&meta.CreatedAt,
	)
	if err != nil {
		return entity.Meta{}, err
	}

	return meta, nil
}

func (r *MetaRepo) GetMetaByFileLocation(fileLocation string) (entity.Meta, error) {
	meta := entity.Meta{}

	query := `
	SELECT
		id, user_id,
		file_location, file_name,
		file_size, file_ext, file_last_modified,
		file_width, file_height,
		created_at
	FROM
		metadata
	WHERE
		file_location = $1
	`

	err := r.db.QueryRow(query, fileLocation).Scan(
		&meta.ID,
		&meta.UserID,
		&meta.FileLocation,
		&meta.FileName,
		&meta.FileSize,
		&meta.FileExt,
		&meta.FileLastModified,
		&meta.FileWidth,
		&meta.FileHeight,
		&meta.CreatedAt,
	)
	if err != nil {
		return entity.Meta{}, err
	}

	return meta, nil
}

func (r *MetaRepo) GetAllMetaByUserId(user_id string) ([]entity.Meta, error) {
	metas := []entity.Meta{}

	query := `
		SELECT
			id, user_id,
			file_location, file_name,
			file_size, file_ext, file_last_modified,
			file_width, file_height,
			created_at
		FROM
			metadata
		WHERE
			user_id = $1
	`

	rows, err := r.db.Query(query, user_id)
	if err != nil {
		return []entity.Meta{}, err
	}
	for rows.Next() {
		meta := entity.Meta{}
		err := rows.Scan(
			&meta.ID,
			&meta.UserID,
			&meta.FileLocation,
			&meta.FileName,
			&meta.FileSize,
			&meta.FileExt,
			&meta.FileLastModified,
			&meta.FileWidth,
			&meta.FileHeight,
			&meta.CreatedAt,
		)
		if err != nil {
			return []entity.Meta{}, err
		}

		query = `
			SELECT
				l.value as value,
				ml.label_id as label_id,
				ml.id as metadata_label_id
			FROM
				metadata_labels ml
			JOIN
				labels l
			ON
				l.id = ml.label_id
			WHERE
				ml.metadata_id = $1
		`
		rows, err := r.db.Query(query, meta.ID)
		if err != nil {
			return []entity.Meta{}, err
		}

		labels := []entity.Label{}
		for rows.Next() {
			label := entity.Label{}
			err := rows.Scan(
				&label.Value,
				&label.ID,
				&label.MetadataLabelID,
			)
			if err != nil {
				return []entity.Meta{}, err
			}
			labels = append(labels, label)
		}

		meta.Labels = labels
		metas = append(metas, meta)
	}

	return metas, nil
}

func (r *MetaRepo) GetAllMeta() ([]entity.Meta, error) {
	return []entity.Meta{}, nil
}

func (r *MetaRepo) UpdateMeta(id string, data entity.Meta) (entity.Meta, error) {
	return entity.Meta{}, nil
}

func (r *MetaRepo) DeleteMeta(id string) (entity.Meta, error) {
	meta := entity.Meta{}

	query := `
	DELETE FROM
		metadata
	WHERE
		id = $1
	RETURNING
		id, user_id,
		file_location, file_name,
		file_size, file_ext, file_last_modified,
		file_width, file_height,
		created_at
	`
	err := r.db.QueryRow(query, id).Scan(
		&meta.ID,
		&meta.UserID,
		&meta.FileLocation,
		&meta.FileName,
		&meta.FileSize,
		&meta.FileExt,
		&meta.FileLastModified,
		&meta.FileWidth,
		&meta.FileHeight,
		&meta.CreatedAt,
	)
	if err != nil {
		return entity.Meta{}, err
	}

	return meta, nil
}

func (r *MetaRepo) DeleteMetaByUser(user_id string) ([]entity.Meta, error) {
	return []entity.Meta{}, nil
}
