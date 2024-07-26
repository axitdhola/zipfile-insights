package db

import (
	"database/sql"
	"log"
	"os"

	models "github.com/axitdhola/zipfile-insights/lambda-zip-extractor/models"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	db_string := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", db_string)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) InsertZipFile(name string, userID int) (models.ZipFile, error) {
	var zipFile models.ZipFile
	err := d.db.QueryRow(`
        INSERT INTO zip_file (name, user_id, s3_key)
        VALUES ($1, $2, $3)
        RETURNING id, name, user_id, s3_key, created_at
    `, name, userID, name).Scan(
		&zipFile.ID, &zipFile.Name, &zipFile.UserID, &zipFile.S3Key, &zipFile.CreatedAt,
	)
	return zipFile, err
}

func (d *Database) InsertExtractedFile(file models.ExtractedFile) error {
	_, err := d.db.Exec(`
        INSERT INTO extracted_files (name, s3_key, zip_id, content, searchable_content, file_size, mime_type)
        VALUES ($1, $2, $3, $4, to_tsvector('english',$4), $5, $6)
    `, file.Name, file.S3Key, file.ZipID, file.Content, file.FileSize, file.MimeType)
	if err != nil {
		log.Printf("%v db.go InsertExtractedFile error: %v", file.Name, err)
		return err
	}
	return err
}
