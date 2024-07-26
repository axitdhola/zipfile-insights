package models

import "time"

type ZipFile struct {
	ID        int
	Name      string
	UserID    int
	S3Key     string
	CreatedAt time.Time
}

type ExtractedFile struct {
	ID        int
	Name      string
	S3Key     string
	ZipID     int
	Content   string
	FileSize  int64
	MimeType  string
	CreatedAt time.Time
}
