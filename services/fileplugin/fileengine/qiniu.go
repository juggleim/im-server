package fileengine

import (
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiNiuStorage struct {
	accessKey string
	secretKey string
	bucket    string
	domain    string
}

type QiNiuConfig struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	Domain    string `json:"domain"`
}

func NewQiNiu(config QiNiuConfig) *QiNiuStorage {
	return &QiNiuStorage{
		accessKey: config.AccessKey,
		secretKey: config.SecretKey,
		bucket:    config.Bucket,
		domain:    config.Domain,
	}
}

func (q *QiNiuStorage) UploadToken(fileType string) (uploadToken string, domain string) {
	putPolicy := storage.PutPolicy{
		Scope: q.bucket,
	}
	mac := qbox.NewMac(q.accessKey, q.secretKey)
	return putPolicy.UploadToken(mac), q.domain
}
