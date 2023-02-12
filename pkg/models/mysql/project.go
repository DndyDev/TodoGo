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
	WHERE id = ? AND is_delete = 0`

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

func (model *ProjectModel) GetUserProjects(id int) ([]*models.Project, error) {
	stmt := `SELECT id, title, web_user_id FROM projects 
	WHERE web_user_id = ?`

	rows, err := model.DB.Query(stmt, id)
	var projects []*models.Project
	for rows.Next() {
		project := &models.Project{}
		err = rows.Scan(&project.ID, &project.Title, &project.UserID)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (model *ProjectModel) Put(id int, title string) error {
	stmt := `UPDATE projects SET title = ? WHERE id = ?`

	_, err := model.DB.Exec(stmt, title, id)

	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	return nil

}

func (model *ProjectModel) Delete(id int) error {

	stmt := `UPDATE projects SET is_delete = ? WHERE id = ?`
	_, err := model.DB.Exec(stmt, 1, id)

	if err != nil {
		return err
	}

	return nil

}
