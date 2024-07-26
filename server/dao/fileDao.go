package dao

import (
	"database/sql"
	"fmt"

	"github.com/axitdhola/zipfile-insights/server/models"
)

type FileDao interface {
	GetAllFiles(userId int) ([]models.FileModel, error)
	SerachFile(userId int, content string) ([]models.FileModel, error)
}

type fileDaoImpl struct {
	db *sql.DB
}

func NewFileDao(db *sql.DB) FileDao {
	return &fileDaoImpl{
		db: db,
	}
}

func (u *fileDaoImpl) GetAllFiles(userId int) ([]models.FileModel, error) {
	var files []models.FileModel

	query := "SELECT ef.Id, ef.zip_id, ef.name, ef.file_size, ef.mime_type, ef.s3_key, ef.created_at FROM extracted_files ef INNER JOIN zip_file zf ON ef.zip_id = zf.id WHERE zf.user_id = $1"
	res, err := u.db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %v", err)
	}
	defer res.Close()

	for res.Next() {
		var file models.FileModel

		err = res.Scan(&file.Id, &file.ZipId, &file.FileName, &file.Size, &file.Type, &file.S3Key, &file.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		files = append(files, file)
	}

	if err = res.Err(); err != nil {
		return nil, fmt.Errorf("result set iteration error: %v", err)
	}

	return files, nil
}

func (u *fileDaoImpl) SerachFile(userId int, content string) ([]models.FileModel, error) {
	// 	SELECT filename FROM files
	// WHERE searchable_content @@ plainto_tsquery('english', 'this')

	var files []models.FileModel

	res, err := u.db.Query("SELECT ef.Id, ef.zip_id, ef.name, ef.file_size, ef.mime_type, ef.s3_key, ef.created_at FROM extracted_files ef INNER JOIN zip_file zf ON ef.zip_id = zf.id WHERE zf.user_id = $1 AND ef.searchable_content @@ plainto_tsquery('english', $2)", userId, content)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var file models.FileModel

		err = res.Scan(&file.Id, &file.ZipId, &file.FileName, &file.Size, &file.Type, &file.S3Key, &file.CreatedAt)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}
