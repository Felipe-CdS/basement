package sqlite

import (
	"database/sql"

	"nugu.dev/basement/pkg/models"
)

type TagRepository struct {
	Db *sql.DB
}

func (r *TagRepository) GetActivityTags() ([]models.Tag, error) {
	var tags []models.Tag

	return tags, nil
}
