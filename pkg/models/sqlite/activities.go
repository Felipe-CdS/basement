package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"nugu.dev/basement/pkg/models"
)

type ActivityRepository struct {
	Db *sql.DB
}

func (a *ActivityRepository) StartActivity() (int, error) {
	return 1, nil
}

func (a *ActivityRepository) EndActivity() error {
	return nil
}

func (a *ActivityRepository) NewCompleteActivity(x models.Activity) (int, error) {
	return 1, nil
}

func (a *ActivityRepository) GetLastActivity() (models.Activity, error) {
	var search models.Activity

	return search, nil
}

func (a *ActivityRepository) GetDailyLog(date time.Time) ([]models.Activity, error) {
	var search []models.Activity

	return search, nil
}

func (a *ActivityRepository) GetSingleDetailedLogById(id string) (models.Activity, error) {
	var search models.Activity

	return search, nil
}

func (a *ActivityRepository) GetIntervalLog(start time.Time, end time.Time, selectedTag string) ([]models.ActivityDayOverview, error) {
	var stmt string
	var err error
	var rows *sql.Rows
	var search []models.ActivityDayOverview

	if selectedTag != "" {
		stmt = `SELECT start_time, FLOOR(SUM(unixepoch(v1.end_time) - unixepoch(v1.start_time)))
			FROM (
				SELECT *
				FROM activities
					JOIN activities_tags ON fk_activity_id = activities.id
					JOIN tags ON fk_tag_id = tags.id
				WHERE activities.start_time >= $1
					AND activities.end_time < $2
					AND activities.end_time IS NOT NULL
					AND tags.name LIKE $3
			) v1
			GROUP BY start_time
			ORDER BY start_time ASC;`

		rows, err = a.Db.Query(stmt, start.Format(time.DateOnly), end.Format(time.DateOnly), selectedTag)

		fmt.Println(err)

	} else {
		stmt = `SELECT start_time, FLOOR(SUM(activities.end_time - activities.start_time))
			FROM activities
			WHERE activities.start_time >= $1
			AND activities.end_time < $2
			AND activities.end_time IS NOT NULL
			GROUP BY start_time
			ORDER BY start_time ASC;`

		rows, err = a.Db.Query(stmt, start.Format(time.DateOnly), end.Format(time.DateOnly))
	}

	if err != nil {
		return search, models.ErrDbOperation
	}

	for rows.Next() {
		var s models.ActivityDayOverview

		if err = rows.Scan(
			&s.Date,
			&s.TotalSec,
		); err != nil {
			return search, models.ErrDbOperation
		}

		search = append(search, s)
	}

	return search, nil
}
