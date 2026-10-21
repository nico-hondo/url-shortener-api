package model

import "time"

type Users struct {
	Id          int        `db:"id"`
	Password    string     `db:"password"`
	First_Name  string     `db:"first_name"`
	Last_Name   string     `db:"last_name"`
	Email       string     `db:"email"`
	Verified_At *time.Time `db:"verified_at"`
	Created_At  time.Time  `db:"created_at"`
	Updated_At  time.Time  `db:"updated_at"`
}
