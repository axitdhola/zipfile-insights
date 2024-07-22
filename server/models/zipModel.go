package models

import "time"

type ZipModel struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	FileName  string    `json:"file_name"`
	CreatedAt time.Time `json:"created_at"`
}
