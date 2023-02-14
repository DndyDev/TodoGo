package mysql

import (
	"database/sql"
	"errors"

	"dandydev.com/todogo/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (model *UserModel) Insert(nickname, lastName, email, password string) error {
	stmt := `INSERT INTO web_users (nickname, last_name,email,user_password) 
	VALUES(?,?,?,?)`

	_, err := model.DB.Exec(stmt, nickname, lastName, email, password)
	if err != nil {
		return err
	}

	return nil
}
func (model *UserModel) Get(nickname, password string) (*models.User, error) {
	stmt := `SELECT id,nickname, last_name,email,user_password,is_ban 
	FROM web_users 
	WHERE nickname = ? AND user_password = ? AND is_ban = 0`

	row := model.DB.QueryRow(stmt, nickname, password)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Nick, &user.LastName, &user.Email,
		&user.Password, &user.IsBan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return user, nil
}

func (model *UserModel) GetAll() ([]*models.User, error) {
	stmt := `SELECT * FROM web_users `

	rows, err := model.DB.Query(stmt)

	var users []*models.User

	for rows.Next() {
		user := &models.User{}
		err = rows.Scan(&user.ID, &user.Nick, &user.LastName, &user.Email,
			&user.Password, &user.IsBan)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
func (model *UserModel) Ban(id int) error {
	stmt := `UPDATE web_users SET is_ban = 1 WHERE id = ?`

	_, err := model.DB.Query(stmt, id)
	if err != nil {
		return err
	}
	return nil
}
func (model *UserModel) UnBan(id int) error {
	stmt := `UPDATE web_users SET is_ban = 0 WHERE id = ?`

	_, err := model.DB.Query(stmt, id)
	if err != nil {
		return err
	}
	return nil
}
