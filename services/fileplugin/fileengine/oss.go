package fileengine

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"im-server/commons/tools"
	"io"
	"path/filepath"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const (
	expireTime = int64(3600)
)

type OssConfig struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Endpoint  string `json:"endpoint"`
	Bucket    string `json:"bucket"`
}

type OssStorage struct {
	accessKeyId, accessKeySecret string
	endpoint, bucket             string
}

func NewOss(conf OssConfig) *OssStorage {
	return &OssStorage{
		accessKeyId:     conf.AccessKey,
		accessKeySecret: conf.SecretKey,
		endpoint:        conf.Endpoint,
		bucket:          conf.Bucket,
	}
}

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}
type PolicyToken struct {
	AccessKeyId string `json:"ossAccessKeyId"`
	Host        string `json:"host"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
}

func getGMTISO8601(expireEnd int64) string {
	return time.Unix(expireEnd, 0).UTC().Format("2006-01-02T15:04:05Z")
}
func (o *OssStorage) getPolicyToken() string {
	var uploadDir = "files/"

	now := time.Now().Unix()
	expireEnd := now + expireTime
	tokenExpire := getGMTISO8601(expireEnd)
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, uploadDir)
	config.Conditions = append(config.Conditions, condition)
	result, err := json.Marshal(config)
	if err != nil {
		fmt.Println("callback json err:", err)
		return ""
	}
	encodedResult := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(sha1.New, []byte(o.accessKeySecret))
	_, _ = io.WriteString(h, encodedResult)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))
	host := fmt.Sprintf("https://%s.%s", o.bucket, o.endpoint)
	policyToken := PolicyToken{
		AccessKeyId: o.accessKeyId,
		Host:        host,
		Signature:   signedStr,
		Policy:      encodedResult,
		Directory:   uploadDir,
	}
	response, err := json.Marshal(policyToken)
	if err != nil {
		fmt.Println("json err:", err)
		return ""
	}
	return string(response)
}

func (o *OssStorage) buildBucket() (bucket *oss.Bucket, err error) {
	var client *oss.Client
	client, err = oss.New(o.endpoint, o.accessKeyId, o.accessKeySecret)
	if err != nil {
		err = fmt.Errorf("oss err:%v", err)
		return
	}
	bucket, err = client.Bucket(o.bucket)

	return
}

func (o *OssStorage) getURL(fileType string, dir string) (signedURL string, err error) {
	bucket, err := o.buildBucket()
	if err != nil {
		err = fmt.Errorf("bucket err:%v", err)
		return
	}
	//options := []oss.Option{
	//	oss.ContentType(contentType(fileType)),
	//}
	objectName := tools.GenerateUUIDShort22() + "." + fileType
	objectName = filepath.Join(dir, objectName)

	signedURL, err = bucket.SignURL(objectName, oss.HTTPPut, expireTime)
	if err != nil {
		err = fmt.Errorf("sign url err:%v", err)
		return
	}
	return signedURL, nil
}

func (o *OssStorage) PreSignedURL(fileType string, dir string) (url string, err error) {
	return o.getURL(fileType, dir)
}
