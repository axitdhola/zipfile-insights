package util

import (
	"os"

	"github.com/axitdhola/zipfile-insights/server/models"
)

func GetBucketName(action string) string {
	if action == string(models.UPLOAD) {
		return os.Getenv("ZIP_BUCKET")
	} else {
		return os.Getenv("FILE_BUCKET")
	}
}
