package main

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gabriel-vasile/mimetype"

	"github.com/axitdhola/zipfile-insights/lambda-zip-extractor/db"
	"github.com/axitdhola/zipfile-insights/lambda-zip-extractor/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func processS3Record(ctx context.Context, s3Client *s3.Client, database *db.Database, record events.S3EventRecord) error {
	sourceBucket := record.S3.Bucket.Name
	sourceKey := record.S3.Object.Key

	zipContent, userId, err := downloadFromS3(ctx, s3Client, sourceBucket, sourceKey)
	if err != nil {
		return fmt.Errorf("failed to download zip file: %v", err)
	}

	zipFile, err := database.InsertZipFile(sourceKey, userId)
	if err != nil {
		return fmt.Errorf("failed to insert zip file record: %v", err)
	}

	err = extractAndUpload(ctx, s3Client, database, zipContent, zipFile.ID)
	if err != nil {
		return fmt.Errorf("failed to extract and upload files: %v", err)
	}

	return nil
}

func extractAndUpload(ctx context.Context, client *s3.Client, database *db.Database, zipContent []byte, zipID int) error {
	zipReader, err := zip.NewReader(bytes.NewReader(zipContent), int64(len(zipContent)))
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		fileContent, err := readZipFile(file)
		if err != nil {
			log.Printf("Error reading file %s: %v", file.Name, err)
			continue
		}

		// utf8Content, err := convertToUTF8(fileContent)
		// if err != nil {
		// 	log.Printf("Error converting file %s to UTF-8: %v", file.Name, err)
		// 	continue
		// }

		s3Key := filepath.Join(strconv.Itoa(zipID), file.Name)
		err = uploadToS3(ctx, client, s3Key, fileContent)
		if err != nil {
			log.Printf("Error uploading file %s: %v", file.Name, err)
			continue
		}

		mimeType := mimetype.Detect(fileContent).String()

		var content string
		var extractionSuccessful bool

		switch {
		case mimeType == "application/pdf":
			content, err = extractPDFContent(fileContent)
			extractionSuccessful = (err == nil)
		case mimeType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
			content, err = extractDOCXContent(fileContent)
			extractionSuccessful = (err == nil)
		case mimeType == "text/plain" || strings.HasSuffix(strings.ToLower(file.Name), ".txt"):
			content, err = extractTextContent(fileContent)
			extractionSuccessful = true
		default:
			extractionSuccessful = false
		}

		if err != nil {
			log.Printf("Error extracting content from file %s: %v", file.Name, err)
			continue
		}

		var sanitizedContent string
		if extractionSuccessful {
			sanitizedContent = sanitizeUTF8(content)
		}

		extractedFile := models.ExtractedFile{
			Name:     file.Name,
			S3Key:    s3Key,
			ZipID:    zipID,
			Content:  sanitizedContent,
			FileSize: int64(len(fileContent)),
			MimeType: mimeType,
		}

		err = database.InsertExtractedFile(extractedFile)
		if err != nil {
			log.Printf("Error inserting extracted file record: %v", err)
		}
	}

	return nil
}

func readZipFile(file *zip.File) ([]byte, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return io.ReadAll(rc)
}

// func convertToUTF8(input []byte) ([]byte, error) {
// 	reader := transform.NewReader(bytes.NewReader(input), charmap.ISO8859_1.NewDecoder())
// 	return ioutil.ReadAll(reader)
// }
