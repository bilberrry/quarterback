package storage

import (
	"os"
	"path"

	logger "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	Base
	bucket string
	path   string
	client *s3manager.Uploader
}

func (ctx *S3) open() (err error) {
	viper := ctx.viper

	viper.SetDefault("region", "us-east-1")

	config := aws.NewConfig()

	endpoint := ctx.viper.GetString("endpoint")

	if len(endpoint) > 0 {
		config.Endpoint = aws.String(endpoint)
	}

	config.Credentials = credentials.NewStaticCredentials(
		viper.GetString("access_key_id"),
		viper.GetString("secret_access_key"),
		viper.GetString("token"),
	)

	config.Region = aws.String(ctx.viper.GetString("region"))
	config.MaxRetries = aws.Int(ctx.viper.GetInt("max_retries"))

	ctx.bucket = ctx.viper.GetString("bucket")
	ctx.path = ctx.viper.GetString("path")

	sess := session.Must(session.NewSession(config))
	ctx.client = s3manager.NewUploader(sess)

	return
}

func (ctx *S3) close() {}

func (ctx *S3) send(fileName string) (err error) {
	f, err := os.Open(ctx.archivePath)

	if err != nil {
		return err
	}

	remotePath := path.Join(ctx.path, fileName)

	input := &s3manager.UploadInput{
		Bucket: aws.String(ctx.bucket),
		Key:    aws.String(remotePath),
		Body:   f,
	}

	logger.Info("=> Uploading to S3...")

	_, err = ctx.client.Upload(input)
	if err != nil {
		logger.Error("Failed to upload file: %v", err)
		return
	}

	logger.Info("=> Stored successfully")
	return nil
}

func (ctx *S3) delete(fileName string) (err error) {
	remotePath := path.Join(ctx.path, fileName)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(ctx.bucket),
		Key:    aws.String(remotePath),
	}

	_, err = ctx.client.S3.DeleteObject(input)
	return
}
