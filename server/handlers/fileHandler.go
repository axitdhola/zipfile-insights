package handlers

import (
	"net/http"
	"strconv"

	"github.com/axitdhola/zipfile-insights/server/models"
	"github.com/axitdhola/zipfile-insights/server/services"
	"github.com/axitdhola/zipfile-insights/server/util"
	"github.com/gin-gonic/gin"
)

type FileHandler interface {
	GetAllFiles(c *gin.Context)
	SerachFile(c *gin.Context)
	GetPresignedUrl(c *gin.Context)
	RedirectToPresignedUrl(c *gin.Context)
}

type fileHandler struct {
	fileService services.FileService
}

func NewFileHandler(fileService services.FileService) FileHandler {
	return &fileHandler{fileService: fileService}
}

func (f *fileHandler) GetAllFiles(c *gin.Context) {
	id := c.Param("uid")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := f.fileService.GetAllFiles(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (f *fileHandler) SerachFile(c *gin.Context) {
	var serachInput models.SearchModel

	if err := c.ShouldBindJSON(&serachInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	res, err := f.fileService.SerachFile(serachInput.UserId, serachInput.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (f *fileHandler) GetPresignedUrl(c *gin.Context) {
	var presignedUrlReq models.PresignedUrlReq
	if err := c.ShouldBindJSON(&presignedUrlReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	presignedUrlReq.Bucket = util.GetBucketName(presignedUrlReq.Action)

	res, err := f.fileService.GetPresignedUrl(presignedUrlReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (f *fileHandler) RedirectToPresignedUrl(c *gin.Context) {
	var presignedUrlReq models.PresignedUrlReq
	if err := c.ShouldBindJSON(&presignedUrlReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	presignedUrlReq.Bucket = util.GetBucketName(presignedUrlReq.Action)

	res, err := f.fileService.RedirectToPresignedUrl(presignedUrlReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
	// c.Header("Access-Control-Allow-Credentials", "true")
	// c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	// c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.JSON(http.StatusOK, res)
}
