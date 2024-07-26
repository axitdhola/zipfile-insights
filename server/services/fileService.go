package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/axitdhola/zipfile-insights/server/dao"
	"github.com/axitdhola/zipfile-insights/server/models"
)

type fileServiceImpl struct {
	fileDao dao.FileDao
}

type FileService interface {
	GetAllFiles(userId int) ([]models.FileModel, error)
	SerachFile(userId int, content string) ([]models.FileModel, error)
	GetPresignedUrl(req models.PresignedUrlReq) (models.PresignedUrlRes, error)
	RedirectToPresignedUrl(req models.PresignedUrlReq) (models.PresignedUrlRes, error)
}

func NewFileService(fileDao dao.FileDao) FileService {
	return &fileServiceImpl{fileDao: fileDao}
}

func (f *fileServiceImpl) GetAllFiles(userId int) ([]models.FileModel, error) {
	if userId == 0 {
		return nil, errors.New("invalid user id")
	}
	return f.fileDao.GetAllFiles(userId)
}

func (f *fileServiceImpl) SerachFile(userId int, content string) ([]models.FileModel, error) {
	if userId == 0 {
		return nil, errors.New("invalid user id")
	}

	if len(content) == 0 {
		return f.GetAllFiles(userId)
	}

	return f.fileDao.SerachFile(userId, content)
}

func (f *fileServiceImpl) GetPresignedUrl(req models.PresignedUrlReq) (models.PresignedUrlRes, error) {
	if req.Action != string(models.UPLOAD) && req.Action != string(models.DOWNLOAD) {
		return models.PresignedUrlRes{}, errors.New("invalid action")
	}

	apiGatewayUrl := os.Getenv("API_GATEWAY_URL")

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return models.PresignedUrlRes{}, err
	}

	resp, err := http.Post(apiGatewayUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return models.PresignedUrlRes{}, err
	}

	var presignedUrlRes models.PresignedUrlRes
	err = json.NewDecoder(resp.Body).Decode(&presignedUrlRes)
	if err != nil {
		return models.PresignedUrlRes{}, err
	}

	return presignedUrlRes, nil
}

func (f *fileServiceImpl) RedirectToPresignedUrl(req models.PresignedUrlReq) (models.PresignedUrlRes, error) {
	if req.Action != string(models.DOWNLOAD) {
		return models.PresignedUrlRes{}, errors.New("invalid action")
	}

	presignedUrl, err := f.GetPresignedUrl(req)
	if err != nil {
		return models.PresignedUrlRes{}, err
	}

	log.Println(presignedUrl.Url)

	return presignedUrl, nil
}
