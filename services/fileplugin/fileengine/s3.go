package fileengine

import (
	"fmt"
	"im-server/commons/tools"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Config struct {
	AccessKey string `json:"access_key,omitempty"`
	SecretKey string `json:"secret_key,omitempty"`
	Endpoint  string `json:"endpoint,omitempty"`
	Region    string `json:"region,omitempty"`
	Bucket    string `json:"bucket,omitempty"`
}

type S3Storage struct {
	accessKey        string
	secretKey        string
	endpoint         string
	region           string
	bucket           string
	disableSSL       bool
	s3forcePathStyle bool
}

type Option func(*S3Storage)

func WithAccessKey(accessKey string) Option {
	return func(o *S3Storage) {
		o.accessKey = accessKey
	}
}
func WithSecretKey(secretKey string) Option {
	return func(o *S3Storage) {
		o.secretKey = secretKey
	}
}
func WithEndpoint(endpoint string) Option {
	return func(o *S3Storage) {
		o.endpoint = endpoint
	}
}
func WithRegion(region string) Option {
	return func(o *S3Storage) {
		o.region = region
	}
}
func WithBucket(bucket string) Option {
	return func(o *S3Storage) {
		o.bucket = bucket
	}
}
func WithConf(conf S3Config) Option {
	return func(o *S3Storage) {
		o.accessKey = conf.AccessKey
		o.secretKey = conf.SecretKey
		o.endpoint = conf.Endpoint
		o.region = conf.Region
		o.bucket = conf.Bucket
	}
}

func NewS3Storage(options ...Option) *S3Storage {
	s := &S3Storage{}
	for _, option := range options {
		option(s)
	}

	return s
}

func (s *S3Storage) putPreSignURL(fileType string, dir string) (url string, err error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(s.accessKey, s.secretKey, ""),
		Endpoint:         aws.String(s.endpoint),
		Region:           aws.String(s.region),
		DisableSSL:       aws.Bool(s.disableSSL),
		S3ForcePathStyle: aws.Bool(s.s3forcePathStyle), //virtual-host style方式
	})
	if err != nil {
		err = fmt.Errorf("error creating S3 session: %v", err)
		return "", err
	}
	svc := s3.New(sess)

	objectName := tools.GenerateUUIDShort22() + "." + fileType
	objectName = filepath.Join(dir, objectName)
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectName),
	})
	url, err = req.Presign(15 * time.Minute)

	return
}
func (s *S3Storage) PreSignedURL(fileType string, dir string) (url string, err error) {
	return s.putPreSignURL(fileType, dir)
}
