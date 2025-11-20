package fileengine

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
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
	Region    string `json:"region"`
	Domain    string `json:"domain"`
}

type OssStorage struct {
	accessKeyId, accessKeySecret string
	endpoint, bucket             string
	region                       string
	domain                       string
}

func NewOss(conf OssConfig) *OssStorage {
	return &OssStorage{
		accessKeyId:     conf.AccessKey,
		accessKeySecret: conf.SecretKey,
		endpoint:        conf.Endpoint,
		bucket:          conf.Bucket,
		region:          conf.Region,
		domain:          conf.Domain,
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

func (o *OssStorage) getURL(fileType string, dir string) (signedURL, objectKey, downUrl string, err error) {
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
	objectKey = objectName
	if o.domain != "" {
		downUrl = o.domain + "/" + objectKey
	}

	signedURL, err = bucket.SignURL(objectName, oss.HTTPPut, expireTime)
	if err != nil {
		err = fmt.Errorf("sign url err:%v", err)
		return
	}
	return
}

func (o *OssStorage) PreSignedURL(fileType string, dir string) (url, objectKey, downUrl string, err error) {
	return o.getURL(fileType, dir)
}

func (o *OssStorage) PostSign(fileType string, dir string) *pbobjs.PreSignResp {
	ret := &pbobjs.PreSignResp{}
	utcTime := time.Now().UTC()
	date := utcTime.Format("20060102")
	//policy
	expiration := utcTime.Add(1 * time.Hour)
	policyMap := map[string]any{
		"expiration": expiration.Format("2006-01-02T15:04:05.000Z"),
		"conditions": []any{
			map[string]string{"bucket": o.bucket},
			map[string]string{"x-oss-signature-version": "OSS4-HMAC-SHA256"},
			map[string]string{"x-oss-credential": fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request",
				o.accessKeyId, date, o.region, "oss")}, // 凭证
			map[string]string{"x-oss-date": utcTime.Format("20060102T150405Z")},
			// 其他条件
			[]any{"content-length-range", 1, 1024 * 1024 * 1024},
			// []any{"eq", "$success_action_status", "201"},
			// []any{"starts-with", "$key", "user/eric/"},
			// []any{"in", "$content-type", []string{"image/jpg", "image/png"}},
			// []any{"not-in", "$cache-control", []string{"no-cache"}},
		},
	}
	policy := tools.ToJson(policyMap)
	strToSign := base64.StdEncoding.EncodeToString([]byte(policy))
	ret.Policy = strToSign
	//signature
	h1 := tools.HmacSha256([]byte("aliyun_v4"+o.accessKeySecret), date)
	h2 := tools.HmacSha256(h1, o.region)
	h3 := tools.HmacSha256(h2, "oss")
	h4 := tools.HmacSha256(h3, "aliyun_v4_request")
	signature := hex.EncodeToString(tools.HmacSha256(h4, strToSign))
	ret.Signature = signature
	ret.SignVersion = "OSS4-HMAC-SHA256"
	ret.Date = utcTime.Format("20060102T150405Z")
	ret.Credential = fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request", o.accessKeyId, date, o.region, "oss")
	return ret
}
