package mysql

import (
	"database/sql"
	"errors"

	"dandydev.com/todogo/pkg/models"
)

type AdminModel struct {
	DB *sql.DB
}

func (model *AdminModel) Get(nickname, password string) (*models.Admin, error) {
	stmt := `SELECT id, nickname, admin_password FROM web_admins 
	WHERE nickname = ? AND admin_password = ?`

	row := model.DB.QueryRow(stmt, nickname, password)
	admin := &models.Admin{}
	err := row.Scan(&admin.ID, &admin.Nick, &admin.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return admin, nil
}
