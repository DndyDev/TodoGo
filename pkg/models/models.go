package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type Note struct {
	ID        int
	Title     string
	Content   string
	Created   time.Time
	Expires   time.Time
	ProjectID int
	StatusID  int
}
type User struct {
	ID       int
	Nick     string
	Password string
	LastName string
	Email    string
	IsBan    bool
}

type Project struct {
	ID     int
	Title  string
	UserID int
}
type Admin struct {
	ID       int
	Nick     string
	Password string
}

type Status struct {
	ID   int
	Name string
}
