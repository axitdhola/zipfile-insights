package models

import "time"

type FileModel struct {
	Id        int       `json:"id"`
	ZipId     int       `json:"zip_id"`
	FileName  string    `json:"file_name"`
	Size      int       `json:"size"`
	Type      string    `json:"type"`
	S3Key     string    `json:"s3_key"`
	CreatedAt time.Time `json:"created_at"`
}

type SearchModel struct {
	UserId  int    `json:"user_id"`
	Content string `json:"content"`
}

type PresignedUrlReq struct {
	Key    string `json:"key"`
	Bucket string `json:"bucket"`
	Action string `json:"action"`
}

type PresignedUrlRes struct {
	Url string `json:"presigned_url"`
}

type Action string

const (
	UPLOAD   Action = "upload"
	DOWNLOAD Action = "read"
)
