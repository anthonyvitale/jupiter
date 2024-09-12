package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonyvitale/jupiter/pkg/camera"
	"github.com/anthonyvitale/moneta"
	"github.com/aws/aws-sdk-go-v2/aws"
	signerv4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: os.Getenv("S3_ENDPOINT"),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     os.Getenv("S3_KEY"),
				SecretAccessKey: os.Getenv("S3_SECRET"),
			},
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	store, err := moneta.New(s3Client, "moneta", signerv4.SwapComputePayloadSHA256ForUnsignedPayloadMiddleware)
	if err != nil {
		log.Fatal(err)
	}

	err = store.Ping(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dir := filepath.Join(homeDir, "skycam", "images")

	// setup camera
	camera, err := camera.New(dir)
	if err != nil {
		log.Fatal(err)
	}

	// take the photo
	filePath, err := camera.Take()
	if err != nil {
		log.Fatal(err)
	}

	// read and upload to blob storage
	img, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	err = store.UploadImage(context.TODO(), strings.TrimPrefix(filePath, fmt.Sprintf("%s/", dir)), bytes.NewBuffer(img))
	if err != nil {
		log.Fatal(err)
	}
}
