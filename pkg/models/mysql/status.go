package mysql

import (
	"database/sql"
	"errors"

	"dandydev.com/todogo/pkg/models"
)

type StatusModel struct {
	DB *sql.DB
}

func (model *ProjectModel) Get(id int) (*models.Project, error) {
	stmt := `SELECT id, status_name FROM todo_status 
	WHERE id = ?`

	row := model.DB.QueryRow(stmt, id)
	status := &models.Status{}
	err := row.Scan(&status.ID, &status.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return status, nil
}

func (model *ProjectModel) GetAllStatus(id int) ([]*models.Project, error) {
	stmt := `SELECT id, status_name	 FROM todo_status`

	rows, err := model.DB.Query(stmt, id)
	var statuses []*models.Status
	for rows.Next() {
		status := &models.Status{}
		err = rows.Scan(&status.ID, &status.Name)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return statuses, nil
}
