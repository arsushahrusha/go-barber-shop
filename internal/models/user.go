package models

import "time"

type User struct {
	ID			string		`db:"id" json:"id"`
	Login 		string		`db:"login" json:"login"`
	Password 	string		`db:"password" json:"-"`
	CreatedAt 	time.Time	`db:"created_at" json:"created_at"`
}