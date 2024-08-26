package fileengine

import (
	"context"
	"fmt"
	"im-server/commons/tools"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/s3utils"
)

type MinioConfig struct {
	AccessKey string `json:"access_key"` //8XmuYosPwFTqHHwdZf5o
	SecretKey string `json:"secret_key"` //q1Lj92oU4BDcVQvrgzS8qqmjEvNcanMKfufwTULC
	Endpoint  string `json:"endpoint"`
	UseSSL    bool   `json:"use_ssl"`
	Bucket    string `json:"bucket"`
}

type MinioStorage struct {
	config MinioConfig
}

func NewMinio(config MinioConfig) *MinioStorage {
	return &MinioStorage{config: config}
}

func (m *MinioStorage) putPreSignURL(fileType string, dir string) (url string, err error) {
	client, err := minio.New(m.config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.config.AccessKey, m.config.SecretKey, ""),
		Secure: m.config.UseSSL,
	})
	if err != nil {
		return "", err
	}
	buckets, err := client.ListBuckets(context.Background())
	fmt.Println(buckets, err)

	if err != nil {
		return "", err
	}

	expiry := 15 * time.Minute
	objectName := tools.GenerateUUIDShort22() + "." + fileType

	fullPath := filepath.Join(dir, objectName)

	if err = s3utils.CheckValidObjectName(fullPath); err != nil {
		return
	}

	preSignedURL, err := client.PresignedPutObject(context.Background(), m.config.Bucket, fullPath, expiry)
	if err != nil {
		return "", err
	}
	url = preSignedURL.String()
	return
}

func (m *MinioStorage) PreSignedURL(fileType string, dir string) (url string, err error) {
	return m.putPreSignURL(fileType, dir)
}
