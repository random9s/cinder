package manager

import (
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/random9s/cinder/sys/mem"
)

//FileManager handles it all
type FileManager struct {
	svc *s3.S3
	d   *s3manager.Downloader
	u   *s3manager.Uploader
}

//NewWithClient creates a file manager from an s3 client
func NewWithClient(svc *s3.S3) *FileManager {
	return &FileManager{
		svc: svc,
		d:   NewDownloader(svc),
		u:   NewUploader(svc),
	}
}

//NewWithRecord creates a file manager from the s3 event record
func NewWithRecord(record events.S3EventRecord) *FileManager {
	svc := NewS3(record.AWSRegion)
	return &FileManager{
		svc: svc,
		d:   NewDownloader(svc),
		u:   NewUploader(svc),
	}
}

//PrimeFile returns file with buffer of size ContentLength
func (fm *FileManager) PrimeFile(bucket, key string) (*mem.File, error) {
	result, err := fm.svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("manager.PrimeFile: %s", err)
	}

	var len = *result.ContentLength
	fp := mem.WithFileSize(len)
	return fp, nil
}

//DownloadToFile downloads from S3 to memfile
func (fm *FileManager) DownloadToFile(bucket, key string, size int64) (*mem.File, error) {
	var fp *mem.File
	var err error
	if size == 0 {
		fp, err = fm.PrimeFile(bucket, key)
		if err != nil {
			return nil, fmt.Errorf("manager.DownloadToFile: %s", err)
		}
	} else if size > 0 {
		fp = mem.WithFileSize(size)
	} else {
		return nil, errors.New("manager.DownloadToFile: negative size")
	}

	_, err = fm.d.Download(fp, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return fp, err
}

//UploadFile uploads files to S3
func (fm *FileManager) UploadFile(bucket, key string, fp *mem.File) error {
	_, err := fm.u.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   fp,
	})

	return err
}
