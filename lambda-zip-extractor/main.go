package main

import (
	"context"
	"fmt"
	"log"

	"github.com/axitdhola/zipfile-insights/lambda-zip-extractor/db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func handleRequest(ctx context.Context, s3Event events.S3Event) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	database, err := db.NewDatabase()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer database.Close()

	for _, record := range s3Event.Records {
		err := processS3Record(ctx, s3Client, database, record)
		if err != nil {
			log.Printf("Error processing record: %v", err)
		}
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}

// $env:GOOS = "linux"
// $env:GOARCH = "amd64"
// $env:CGO_ENABLED = "0"
// go build -o bootstrap .
// Compress-Archive -Path bootstrap -DestinationPath lambda-handler.zip
