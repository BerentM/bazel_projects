package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// AwsHelper wapper around AWS S3 operations
type AwsHelper struct {
	Region             string `yaml:"Region"`
	Bucket             string `yaml:"Bucket"`
	CredentialsProfile string `yaml:"CredentialsProfile"`
	Config             aws.Config
}

// New initialize AWS session
func (ah *AwsHelper) New() {
	ah.Region = "eu-central-1"
	ah.Bucket = "dtmx-images-poc"
	ah.CredentialsProfile = "dtmx-images"
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(ah.CredentialsProfile),
	)
	if err != nil {
		exitErrorf("Unable to load credentials, %v", err)
	}
	ah.Config = cfg
}

// CheckBuckets list all buckets
func (ah *AwsHelper) CheckBuckets() {
	// Create S3 service client
	svc := s3.NewFromConfig(ah.Config)

	result, err := svc.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			*b.Name, b.CreationDate)
	}
}

func (ah *AwsHelper) checkIfFileExists(svc *s3.Client, uniqueID string) (bool, error) {
	input := &s3.HeadObjectInput{
		Bucket: &ah.Bucket,
		Key:    &uniqueID,
	}
	_, err := svc.HeadObject(context.TODO(), input)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// Upload the object to S3 using the unique identifier as the key
func (ah *AwsHelper) Upload(byteFile []byte, uniqueID string) {
	svc := s3.NewFromConfig(ah.Config)
	exist, err := ah.checkIfFileExists(svc, uniqueID)
	if exist {
		fmt.Println("File already exist in S3")
		return
	}

	// Create an uploader with S3 client and default options
	result, err := svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &ah.Bucket,
		Key:    &uniqueID,
		Body:   bytes.NewReader(byteFile),
	})
	if err != nil {
		log.Printf("Couldn't upload file to %v:%v. Here's why: %v\n",
			ah.Bucket, uniqueID, err)
	}

	// Perform an upload.
	fmt.Println(result)
}

// // Download get ObjectData from S3
// func (ah *AwsHelper) Download(uniqueID string) {
// 	svc := s3.New(ah.Session)
// 	output, err := svc.GetObject(&s3.GetObjectInput{
// 		Bucket: ah.Bucket,
// 		Key:    uniqueID,
// 	})
// 	if err != nil {
// 		exitErrorf("Unable to download file, %v", err)
// 	}

// 	// Process the retrieved object
// 	fmt.Println("Retrieved object:", output.Body)
// }
