package main

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	destinationBucket = "zip-extracted-files"
)

func handleRequest(ctx context.Context, s3Event events.S3Event) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	for _, record := range s3Event.Records {
		log.Println("File name:", record.S3.Object.Key)
		log.Println("record:", record)
		sourceBucket := record.S3.Bucket.Name
		sourceKey := record.S3.Object.Key

		decodedKey, err := url.QueryUnescape(sourceKey)
		if err != nil {
			return fmt.Errorf("failed to decode key: %v", err)
		}

		zipContent, err := downloadFromS3(ctx, s3Client, sourceBucket, decodedKey)
		if err != nil {
			return fmt.Errorf("failed to download zip file: %v", err)
		}

		err = extractAndUpload(ctx, s3Client, zipContent)
		if err != nil {
			return fmt.Errorf("failed to extract and upload files: %v", err)
		}
	}

	return nil
}

func downloadFromS3(ctx context.Context, client *s3.Client, bucket, key string) ([]byte, error) {
	result, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

func extractAndUpload(ctx context.Context, client *s3.Client, zipContent []byte) error {
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

		log.Println("File name:", file.Name, "File size:", len(fileContent))

		err = uploadToS3(ctx, client, file.Name, fileContent)
		if err != nil {
			log.Printf("Error uploading file %s: %v", file.Name, err)
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

func uploadToS3(ctx context.Context, client *s3.Client, fileName string, content []byte) error {
	bucket := destinationBucket
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &fileName,
		Body:   bytes.NewReader(content),
	})
	return err
}

func main() {
	lambda.Start(handleRequest)
}

// $env:GOOS = "linux"
// $env:GOARCH = "amd64"
// $env:CGO_ENABLED = "0"
// go build -o bootstrap main.go
// Compress-Archive -Path bootstrap -DestinationPath lambda-handler.zip
