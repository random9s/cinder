package manager

import (
	"net"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//BYTES
const (
	B  = 1         //BYTE
	KB = 1024 * B  //KILOBYTE
	MB = 1000 * KB //MEGABYTE
	GB = 1000 * MB //GIGABYTE
)

//NewS3 creates an s3 client in the given region
func NewS3(region string) *s3.S3 {
	var s = session.Must(session.NewSession(&aws.Config{
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   100,
				IdleConnTimeout:       30 * time.Second,
				TLSHandshakeTimeout:   5 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: time.Second * 30,
		},
		Region: aws.String(region),
		CredentialsChainVerboseErrors: aws.Bool(true),
		MaxRetries:                    aws.Int(3),
	}))

	return s3.New(s)
}

//NewDownloader creates a download manager from the S3 Client and is then customized to chunk
//downloads by 50MB/chunk
func NewDownloader(svc *s3.S3) *s3manager.Downloader {
	return s3manager.NewDownloaderWithClient(svc, func(d *s3manager.Downloader) {
		d.PartSize = 32 * MB
	})
}

//NewUploader creates an upload manager from the s3 client and is then customized to chunk
//uploads by 30MB/chunk
func NewUploader(svc *s3.S3) *s3manager.Uploader {
	return s3manager.NewUploaderWithClient(svc, func(d *s3manager.Uploader) {
		d.PartSize = 32 * MB
	})
}
