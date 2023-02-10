// notes
package mysql

import (
	"database/sql"
	"errors"

	"dandydev.com/todogo/pkg/models"
)

type NoteModel struct {
	DB *sql.DB
}

func (model *NoteModel) Latest() ([]*models.Note, error) {
	stmt := `SELECT id, title, content, created, expires FROM notes 
	WHERE  is_delete = 0 AND 
	expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10 `

	rows, err := model.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*models.Note

	for rows.Next() {
		note := &models.Note{}
		err = rows.Scan(&note.ID, &note.Title, &note.Content, &note.Created,
			&note.Expires)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (model *NoteModel) Insert(title, content, expires,
	projectId, statusId string) (int, error) {
	stmt := `INSERT INTO notes (title, content, created, expires,project_id,status_id) 
	VALUES(?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(),INTERVAL ? DAY),?,?)`

	result, err := model.DB.Exec(stmt, title, content, expires,
		projectId, statusId)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}
func (model *NoteModel) Get(id int) (*models.Note, error) {
	stmt := `SELECT id, title, content, created, expires, project_id, status_id FROM notes 
	WHERE expires > UTC_TIMESTAMP() AND id = ? AND is_delete = 0`

	row := model.DB.QueryRow(stmt, id)
	note := &models.Note{}
	err := row.Scan(&note.ID, &note.Title, &note.Content,
		&note.Created, &note.Expires, &note.ProjectID, &note.StatusID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return note, nil
}

func (model *NoteModel) GetProjectNotes(projectId int) ([]*models.Note, error) {
	stmt := `SELECT id, title, content,created, expires, project_id, status_id FROM notes
	WHERE project_id = ? AND is_delete = 0 AND expires > UTC_TIMESTAMP()`

	rows, err := model.DB.Query(stmt, projectId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notes []*models.Note

	for rows.Next() {
		note := &models.Note{}
		err = rows.Scan(&note.ID, &note.Title, &note.Content, &note.Created,
			&note.Expires, &note.ProjectID, &note.StatusID)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil

}

func (model *NoteModel) Put(title, content, expires,
	projectId, statusId string) (int, error) {
	stmt := `UPDATE notes title = ?, content = ?,expires = ?,
	project_id = ?,status_id = ?)`

	result, err := model.DB.Exec(stmt, title, content, expires,
		projectId, statusId)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil

}

func (model *NoteModel) Delete(id int) error {

	stmt := `UPDATE notes SET is_delete = ? WHERE id = ?`
	_, err := model.DB.Exec(stmt, 1, id)

	if err != nil {
		return err
	}

	return nil

}
