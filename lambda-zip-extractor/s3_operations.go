package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	destinationBucket = "zip-extracted-files"
)

func downloadFromS3(ctx context.Context, client *s3.Client, bucket, key string) ([]byte, int, error) {
	decodedKey, err := url.QueryUnescape(key)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode key: %v", err)
	}

	result, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &decodedKey,
	})
	if err != nil {
		return nil, 0, err
	}
	defer result.Body.Close()

	parts := strings.Split(decodedKey, "/")
	if len(parts) < 2 {
		return nil, 0, fmt.Errorf("invalid S3 key format")
	}
	userID, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse userID: %v", err)
	}

	content, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, userID, fmt.Errorf("failed to read content: %v", err)
	}
	return content, userID, nil
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
