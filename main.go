package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	uploadRequest bool
	bucket        string
	key           string
	file          string
	partSize      int64
	concurrent    int
)

func upload(sess *session.Session) {
	fileIO, err := os.Open(file)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = partSize * 1024 * 1024
		u.Concurrency = concurrent
	})

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   fileIO,
	})
	if err != nil {
		log.Fatalf("failed to upload: %s", err)
	}

	fmt.Printf("upload result: %#v\n", result)
}

func download(sess *session.Session) {
	fileIO, err := os.Create(file)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	downloader := s3manager.NewDownloader(sess, func(d *s3manager.Downloader) {
		d.PartSize = partSize * 1024 * 1024
		d.Concurrency = concurrent
	})

	n, err := downloader.Download(fileIO,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		log.Fatalf("failed to download: %s", err)
	}

	fmt.Printf("successfully downloaded %d bytes\n", n)
}

func main() {

	flag.BoolVar(&uploadRequest, "upload", false, "whether upload/download object")
	flag.StringVar(&bucket, "bucket", "", "bucket to upload/download object")
	flag.StringVar(&key, "key", "", "key to upload/download")
	flag.StringVar(&file, "file", "", "local file path to upload/download")
	flag.Int64Var(&partSize, "part-size", 5, "part size in MB to upload/download")
	flag.IntVar(&concurrent, "concurrent", 10, "number of parts to upload/download in parallel")
	flag.Parse()

	fmt.Printf("%s\n", os.Args)

	region := os.Getenv("REGION")
	accessServer := os.Getenv("ACCESS_SERVER")
	accessKey := os.Getenv("ACCESS_KEY")
	accessSecret := os.Getenv("ACCESS_SECRET")

	if region == "" || accessServer == "" || accessKey == "" || accessSecret == "" {
		log.Fatalf("You must define environment variable: REGION, ACCESS_SERVER, ACCESS_KEY and ACCESS_SECRET!")
	}

	if bucket == "" || key == "" || file == "" {
		log.Fatalf("You must define bucket, key and file!")
	}

	cfg := &aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(accessServer),
		Credentials:      credentials.NewStaticCredentials(accessKey, accessSecret, ""),
		DisableSSL:       aws.Bool(false),
		S3ForcePathStyle: aws.Bool(true),
	}

	sess := session.New(cfg)

	if uploadRequest {
		upload(sess)
	} else {
		download(sess)
	}
}
