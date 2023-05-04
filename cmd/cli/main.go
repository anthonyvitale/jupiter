package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/anthonyvitale/jupiter/pkg/moneta"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	log.Println("hello there")

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
		panic(err)
	}

	spacesClient := s3.NewFromConfig(cfg)

	store, err := moneta.New(spacesClient, "moneta")
	if err != nil {
		panic(err)
	}

	err = store.Ping(context.TODO())
	if err != nil {
		panic(err)
	}

	log.Println("trying to upload file")
	err = store.Upload(context.TODO(), "2023/my_file.txt", bytes.NewBufferString("hello there"))
	if err != nil {
		panic(err)
	}

}
