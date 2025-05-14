package services

import (
	"context"
	"encoding/json"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/fileplugin/dbs"
	"im-server/services/fileplugin/fileengine"
	"time"
)

var fileConfCache *caches.LruCache
var fileConfLocks *tools.SegmentatedLocks
var notExistFileConf *FileConfItem

type FileConfItem struct {
	AppKey     string `json:"app_key,omitempty"`
	FileEngine string `json:"file_engine"`
	//QiniuConfig *QiniuFileConfig `json:"qiniu,omitempty"`

	QiNiu *fileengine.QiNiuStorage
	Oss   *fileengine.OssStorage
	Minio *fileengine.MinioStorage
	S3    *fileengine.S3Storage
}

func init() {
	fileConfCache = caches.NewLruCacheWithAddReadTimeout("fileconf_cache", 10000, nil, 5*time.Minute, 10*time.Minute)
	fileConfLocks = tools.NewSegmentatedLocks(128)
	notExistFileConf = &FileConfItem{}
}

func GetFileConf(ctx context.Context, appkey string) *FileConfItem {
	if obj, exist := fileConfCache.Get(appkey); exist {
		return obj.(*FileConfItem)
	} else {
		lock := fileConfLocks.GetLocks(appkey)
		lock.Lock()
		defer lock.Unlock()

		if obj, exist := fileConfCache.Get(appkey); exist {
			return obj.(*FileConfItem)
		} else { //load from db
			fileConf, err := loadFileConfFromDb(appkey)
			if err != nil {
				fileConf = notExistFileConf
			}
			fileConfCache.Add(appkey, fileConf)
			return fileConf
		}
	}
}

func loadFileConfFromDb(appkey string) (*FileConfItem, error) {
	dao := dbs.FileConfDao{}
	conf, err := dao.FindEnableFileConf(appkey)
	if err != nil {
		return nil, err
	}
	fileConf := &FileConfItem{
		AppKey:     appkey,
		FileEngine: conf.Channel,
	}
	var confData = make(map[string]interface{})
	_ = json.Unmarshal([]byte(conf.Conf), &confData)

	switch conf.Channel {
	case fileengine.ChannelQiNiu:
		c := tools.MapToStruct[fileengine.QiNiuConfig](confData)
		fileConf.QiNiu = fileengine.NewQiNiu(c)
	case fileengine.ChannelMinio:
		c := tools.MapToStruct[fileengine.MinioConfig](confData)
		fileConf.Minio = fileengine.NewMinio(c)
	case fileengine.ChannelOss:
		c := tools.MapToStruct[fileengine.OssConfig](confData)
		fileConf.Oss = fileengine.NewOss(c)
	case fileengine.ChannelAws:
		c := tools.MapToStruct[fileengine.S3Config](confData)
		fileConf.S3 = fileengine.NewS3Storage(fileengine.WithConf(c))
	}
	return fileConf, nil
}

func GetFileCred(ctx context.Context, req *pbobjs.QryFileCredReq) (errs.IMErrorCode, *pbobjs.QryFileCredResp) {
	appKey := bases.GetAppKeyFromCtx(ctx)
	fileConf := GetFileConf(ctx, appKey)
	if fileConf == nil || fileConf == notExistFileConf {
		return errs.IMErrorCode_OTHER_NOOSS, nil
	}

	dir := fileTypeToDir(req.FileType)

	switch fileConf.FileEngine {
	case fileengine.ChannelQiNiu:
		if fileConf.QiNiu == nil {
			return errs.IMErrorCode_OTHER_NOOSS, nil
		}
		uploadToken, domain := fileConf.QiNiu.UploadToken(req.Ext)
		return errs.IMErrorCode_SUCCESS, &pbobjs.QryFileCredResp{
			OssType: pbobjs.OssType_QiNiu,
			OssOf: &pbobjs.QryFileCredResp_QiniuCred{
				QiniuCred: &pbobjs.QiNiuCredResp{
					Domain: domain,
					Token:  uploadToken,
				},
			},
		}
	case fileengine.ChannelMinio:
		if fileConf.Minio == nil {
			return errs.IMErrorCode_OTHER_NOOSS, nil
		}
		signedURL, err := fileConf.Minio.PreSignedURL(req.Ext, dir)
		if err != nil {
			return errs.IMErrorCode_OTHER_SIGNERR, nil
		}
		return errs.IMErrorCode_SUCCESS, &pbobjs.QryFileCredResp{
			OssType: pbobjs.OssType_Minio,
			OssOf:   &pbobjs.QryFileCredResp_PreSignResp{PreSignResp: &pbobjs.PreSignResp{Url: signedURL}},
		}
	case fileengine.ChannelAws:
		if fileConf.S3 == nil {
			return errs.IMErrorCode_OTHER_NOOSS, nil
		}

		signedURL, err := fileConf.S3.PreSignedURL(req.Ext, dir)
		if err != nil {
			return errs.IMErrorCode_OTHER_SIGNERR, nil
		}
		return errs.IMErrorCode_SUCCESS, &pbobjs.QryFileCredResp{
			OssType: pbobjs.OssType_S3,
			OssOf:   &pbobjs.QryFileCredResp_PreSignResp{PreSignResp: &pbobjs.PreSignResp{Url: signedURL}},
		}
	case fileengine.ChannelOss:
		if fileConf.Oss == nil {
			return errs.IMErrorCode_OTHER_NOOSS, nil
		}
		signedURL, err := fileConf.Oss.PreSignedURL(req.Ext, dir)
		if err != nil {
			return errs.IMErrorCode_OTHER_SIGNERR, nil
		}
		resp := fileConf.Oss.PostSign(req.Ext, dir)
		return errs.IMErrorCode_SUCCESS, &pbobjs.QryFileCredResp{
			OssType: pbobjs.OssType_Oss,
			OssOf: &pbobjs.QryFileCredResp_PreSignResp{PreSignResp: &pbobjs.PreSignResp{
				Url:         signedURL,
				ObjKey:      resp.ObjKey,
				Policy:      resp.Policy,
				SignVersion: resp.SignVersion,
				Credential:  resp.Credential,
				Date:        resp.Date,
				Signature:   resp.Signature,
			}},
		}
	default:
		return errs.IMErrorCode_OTHER_NOOSS, nil
	}

}

func fileTypeToDir(fileType pbobjs.FileType) string {
	switch fileType {
	case pbobjs.FileType_Image:
		return "images"
	case pbobjs.FileType_Video:
		return "videos"
	case pbobjs.FileType_Audio:
		return "audios"
	case pbobjs.FileType_File:
		return "files"
	case pbobjs.FileType_Log:
		return "logs"
	default:
		return "files"
	}
}
