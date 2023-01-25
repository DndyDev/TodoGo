// notes
package mysql

import (
	"database/sql"

	"dandydev.com/todogo/pkg/models"
)

type NoteModel struct {
	DB *sql.DB
}

func (m *NoteModel) Insert(title, content, expires string) (int, error) {
	return 0, nil
}
func (m *NoteModel) Get(id int) (*models.Note, error) {
	return nil, nil
}
func (m *NoteModel) Latest() ([]*models.Note, error) {
	return nil, nil
}
