package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Client *s3.Client
	Bucket string
}

func NewS3Client(bucket string) *S3Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return &S3Client{
		Client: s3.NewFromConfig(cfg),
		Bucket: bucket,
	}
}

func (s *S3Client) GetObject(key string) ([]byte, error) {
	output, err := s.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Printf("Error getting object from S3: %v", err)
		return nil, fmt.Errorf("failed to get object: %v", err)
	}
	defer output.Body.Close()

	data, readErr := ioutil.ReadAll(output.Body)
	if readErr != nil {
		log.Printf("Error reading object body: %v", readErr)
		return nil, fmt.Errorf("failed to read object body: %v", readErr)
	}

	log.Printf("Successfully retrieved object from S3: %s", key)
	return data, nil
}

func (s *S3Client) PutObject(key string, content []byte) error {
	_, err := s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   ioutil.NopCloser(bytes.NewReader(content)),
	})
	return err
}

func (s *S3Client) UploadFile(key string, body io.Reader) (string, error) {
	_, err := s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   ioutil.NopCloser(body),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.Bucket, key), nil
}
