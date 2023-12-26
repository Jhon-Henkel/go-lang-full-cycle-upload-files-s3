package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"os"
	"sync"
)

var (
	s3Client  *s3.S3
	s3Bucket  string
	s3Id      string = "INSERT_YOUR_AWS_ID_HERE"
	s3Secret  string = "INSERT_YOUR_AWS_SECRET_HERE"
	s3Token   string = "INSERT_YOUR_AWS_TOKEN_HERE"
	waitGroup sync.WaitGroup
)

func init() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(s3Id, s3Secret, s3Token),
	})
	if err != nil {
		panic(err)
	}
	s3Client = s3.New(sess)
	s3Bucket = "example-bucket-go-expert"
}

func main() {
	dir, err := os.Open("./tmp")
	if err != nil {
		panic(err)
	}
	defer dir.Close()
	uploaderControl := make(chan struct{}, 100)
	for {
		files, err := dir.Readdir(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Failed to read directory: %s\n", err)
		}
		waitGroup.Add(1)
		uploaderControl <- struct{}{}
		go uploadFile(files[0].Name(), uploaderControl)
	}
	waitGroup.Wait()
}

func uploadFile(filename string, uploadControl <-chan struct{}) {
	defer waitGroup.Done()
	completeFileName := fmt.Sprintf("./tmp/%s", filename)
	fmt.Printf("Uploading file %s to bucket %s\n", completeFileName, s3Bucket)
	file, err := os.Open(completeFileName)
	if err != nil {
		fmt.Printf("Failed to open file %s\n", completeFileName)
		<-uploadControl
		return
	}
	defer file.Close()
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(filename),
		Body:   file,
	})
	<-uploadControl
	if err != nil {
		fmt.Printf("Failed to upload file %s\n", completeFileName)
		return
	}
	fmt.Printf("Successfully uploaded file %s\n", completeFileName)
}
