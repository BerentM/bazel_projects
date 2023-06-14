package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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
	Session            *session.Session
}

// New initialize AWS session
func (ah *AwsHelper) New() {
	ah.Region = "eu-central-1"
	ah.Bucket = "dtmx-images-poc"
	ah.CredentialsProfile = "dtmx-images"
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: ah.CredentialsProfile,
		Config: aws.Config{
			Region: aws.String(ah.Region),
		},
	})

	_, err = sess.Config.Credentials.Get()
	if err != nil {
		exitErrorf("Unable to load credentials, %v", err)
	}
	ah.Session = sess
}

// CheckBuckets list all buckets
func (ah *AwsHelper) CheckBuckets() {
	// Create S3 service client
	svc := s3.New(ah.Session)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}

func (ah *AwsHelper) checkIfFileExists(svc *s3.S3, uniqueID string) (bool, error) {
	_, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(ah.Bucket),
		Key:    aws.String(uniqueID),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound": // s3.ErrCodeNoSuchKey does not work, aws is missing this error code so we hardwire a string
				return false, nil
			default:
				return false, err
			}
		}
		return false, err
	}
	return true, nil
}

// Upload the object to S3 using the unique identifier as the key
func (ah *AwsHelper) Upload(byteFile []byte, uniqueID string) {
	svc := s3.New(ah.Session)
	exist, err := ah.checkIfFileExists(svc, uniqueID)
	if exist {
		fmt.Println("File already exist in S3")
		return
	}

	// Create an uploader with S3 client and default options
	uploader := s3manager.NewUploaderWithClient(svc)
	upParams := &s3manager.UploadInput{
		Bucket: aws.String(ah.Bucket),
		Key:    aws.String(uniqueID),
		Body:   bytes.NewReader(byteFile),
	}

	// Perform an upload.
	result, err := uploader.Upload(upParams)
	if err != nil {
		exitErrorf("Unable to upload file, %v", err)
	}
	fmt.Println(result.Location)
}

// Download get ObjectData from S3
func (ah *AwsHelper) Download(uniqueID string) {
	svc := s3.New(ah.Session)
	output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(ah.Bucket),
		Key:    aws.String(uniqueID),
	})
	if err != nil {
		exitErrorf("Unable to download file, %v", err)
	}

	// Process the retrieved object
	fmt.Println("Retrieved object:", output.Body)
}