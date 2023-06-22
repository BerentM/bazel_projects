package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// S3Client wapper around AWS S3 operations
type S3Client struct {
	Region string `yaml:"Region"`
	Bucket string `yaml:"Bucket"`
	Client *s3.Client
}

// NewS3Client create new S3 client
func NewS3Client(credProfile string) *S3Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(credProfile),
	)
	if err != nil {
		exitErrorf("Unable to load credentials, %v", err)
	}

	return &S3Client{
		Region: "eu-central-1",
		Bucket: "dtmx-images-poc",
		Client: s3.NewFromConfig(cfg),
	}
}

// CheckBuckets list all buckets
func (c *S3Client) CheckBuckets() {
	result, err := c.Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")
	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			*b.Name, b.CreationDate)
	}
}

func (c *S3Client) checkIfFileExists(uniqueID string) (bool, error) {
	input := &s3.HeadObjectInput{
		Bucket: &c.Bucket,
		Key:    &uniqueID,
	}
	_, err := c.Client.HeadObject(context.TODO(), input)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// Upload the object to S3 using the unique identifier as the key
func (c *S3Client) Upload(byteFile []byte, uniqueID string) {
	exist, err := c.checkIfFileExists(uniqueID)
	if exist {
		fmt.Println("File already exist in S3")
		return
	}

	// Create an uploader with S3 client and default options
	result, err := c.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &c.Bucket,
		Key:    &uniqueID,
		Body:   bytes.NewReader(byteFile),
	})
	if err != nil {
		log.Printf("Couldn't upload file to %v:%v. Here's why: %v\n",
			c.Bucket, uniqueID, err)
	}

	// Perform an upload.
	fmt.Println(result)
}

// Download get ObjectData from S3
func (c *S3Client) Download(uniqueID string) {
	output, err := c.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &c.Bucket,
		Key:    &uniqueID,
	})
	if err != nil {
		exitErrorf("Unable to download file, %v", err)
	}

	// Process the retrieved object
	fmt.Println("Retrieved object:", output.Body)
}
