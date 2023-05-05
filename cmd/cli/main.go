package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/anthonyvitale/jupiter/pkg/moneta"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           "https://nyc3.digitaloceanspaces.com",
			SigningRegion: "us-east-2",
		}, nil
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     os.Getenv("SPACES_KEY"),
				SecretAccessKey: os.Getenv("SPACES_SECRET"),
			},
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	spacesClient := s3.NewFromConfig(cfg)

	store, err := moneta.New(spacesClient, "moneta")
	if err != nil {
		log.Fatal(err)
	}

	err = store.Ping(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create directory to save picture to
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now().UTC()
	localDir := filepath.Join(homeDir, "skycam", "images")

	imgPath := now.Format("2006/01/02")

	err = os.MkdirAll(filepath.Join(localDir, imgPath), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	filePath := fmt.Sprintf("%s.jpeg", filepath.Join(localDir, imgPath, now.Format("150405Z")))
	log.Printf("taking img, saving to %s", filePath)
	cmd := exec.Command("libcamera-jpeg", "-o", filePath)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	slurp, _ := io.ReadAll(stderr)
	fmt.Printf("%s\n", slurp)

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	img, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	err = store.UploadImage(context.TODO(), strings.TrimPrefix(filePath, fmt.Sprintf("%s/", localDir)), bytes.NewBuffer(img))
	if err != nil {
		log.Fatal(err)
	}
}
