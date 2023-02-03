package mysql

import (
	"database/sql"
	"errors"

	"dandydev.com/todogo/pkg/models"
)

type ProjectModel struct {
	DB *sql.DB
}

func (model *ProjectModel) Insert(title string, user_id int) (int, error) {
	stmt := `INSERT INTO projects (title, web_user_id) 
	VALUES(?,?)`

	result, err := model.DB.Exec(stmt, title, user_id)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
func (model *ProjectModel) Get(id int) (*models.Project, error) {
	stmt := `SELECT id, title, web_user_id FROM projects 
	WHERE id = ?`

	row := model.DB.QueryRow(stmt, id)
	project := &models.Project{}
	err := row.Scan(&project.ID, &project.Title, &project.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return project, nil
}
