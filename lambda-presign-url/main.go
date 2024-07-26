package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	s3request "github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type RequestBody struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	Action string `json:"action"`
}

type ResponseBody struct {
	PresignedURL string `json:"presigned_url"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody RequestBody
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("Error parsing request body: %s", err),
		}, nil
	}

	bucket := requestBody.Bucket
	key := requestBody.Key
	action := requestBody.Action

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-south-1"
	}

	log.Printf("bucket: %s, key: %s, action: %s, region: %s", bucket, key, action, region)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Printf("failed to create session: %s", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error creating AWS session: %s", err),
		}, nil
	}

	svc := s3.New(sess)

	var req *s3request.Request

	switch action {
	case "upload":
		req, _ = svc.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			// ContentType: aws.String("application/pdf"),
		})
	case "read":
		req, _ = svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid action. Use 'upload' for PUT or 'read' for GET.",
		}, nil
	}

	urlStr, err := req.Presign(3 * time.Minute)
	if err != nil {
		log.Printf("failed to sign request: %s", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error signing request: %s", err),
		}, nil
	}

	log.Printf("presigned URL: %s", urlStr)

	responseBody, err := json.Marshal(ResponseBody{
		PresignedURL: urlStr,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error marshalling response body: %s", err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type":                     "text/plain",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Headers":     "*",
			"Access-Control-Allow-Methods":     "OPTIONS,POST,GET",
			"Access-Control-Allow-Credentials": "true",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}

// $env:GOOS = "linux"
// $env:GOARCH = "amd64"
// $env:CGO_ENABLED = "0"
// go build -o bootstrap main.go
// Compress-Archive -Path bootstrap -DestinationPath lambda-handler.zip
