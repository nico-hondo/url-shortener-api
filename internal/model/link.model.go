package model

import "time"

type Link struct {
	User_Id      int       `db:"user_id"`
	Original_Url string    `db:"original_url"`
	Slug         string    `db:"slug"`
	Created_At   time.Time `db:"created_at"`
}
