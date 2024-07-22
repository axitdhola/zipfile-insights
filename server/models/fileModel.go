package models

import "time"

type FileModel struct {
	Id        int       `json:"id"`
	ZipId     int       `json:"zip_id"`
	FileName  string    `json:"file_name"`
	Size      int       `json:"size"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
