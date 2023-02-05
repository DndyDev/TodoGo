package mysql

import (
	"database/sql"
	"errors"

	"dandydev.com/todogo/pkg/models"
)

type StatusModel struct {
	DB *sql.DB
}

func (model *StatusModel) Get(id int) (*models.Status, error) {
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

func (model *StatusModel) GetAllStatus() ([]*models.Status, error) {
	stmt := `SELECT id, status_name	 FROM todo_status`

	rows, err := model.DB.Query(stmt)
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
