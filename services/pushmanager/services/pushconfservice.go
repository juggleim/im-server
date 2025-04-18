package services

import (
	"context"
	"encoding/json"
	"im-server/commons/caches"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/pushmanager/services/fcmpush"
	"im-server/services/pushmanager/services/hwpush"
	"im-server/services/pushmanager/services/jpush"
	"im-server/services/pushmanager/services/oppopush"
	"im-server/services/pushmanager/services/vivopush"
	"im-server/services/pushmanager/services/xiaomipush"
	"im-server/services/pushmanager/storages/dbs"
	"strings"
	"time"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

var (
	pushConfCache *caches.LruCache
	pushConfLocks *tools.SegmentatedLocks
)

func init() {
	pushConfCache = caches.NewLruCacheWithAddReadTimeout("pushconf_cache", 10000, nil, 5*time.Minute, 5*time.Minute)
	pushConfLocks = tools.NewSegmentatedLocks(128)
}

type PushConfWraper struct {
	IosPushConf     *IosPushConf
	AndroidPushConf *AndroidPushConf
}
type IosPushConf struct {
	Package        string
	ApnsClient     *apns2.Client
	ApnsVoipClient *apns2.Client
}
type AndroidPushConf struct {
	Package          string
	HwPushClient     *hwpush.HwPushClient
	XiaomiPushClient *xiaomipush.XiaomiPushClient
	OppoPushClient   *oppopush.OppoPushClient
	VivoPushClient   *vivopush.VivoPushClient
	JpushClient      *jpush.JpushClient
	FcmPushClient    *fcmpush.FcmPushClient
}

func GetIosPushConf(ctx context.Context, appkey, packageName string) *IosPushConf {
	key := strings.Join([]string{appkey, packageName}, "_")
	if obj, exist := pushConfCache.Get(key); exist {
		pushConf := obj.(*PushConfWraper)
		if pushConf.IosPushConf != nil {
			if pushConf.IosPushConf == noExistIosPushConf {
				return nil
			}
			return pushConf.IosPushConf
		}
	}
	lock := pushConfLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	var pushConf *PushConfWraper
	if obj, exist := pushConfCache.Get(key); exist {
		pushConf = obj.(*PushConfWraper)
		if pushConf.IosPushConf != nil {
			if pushConf.IosPushConf == noExistIosPushConf {
				return nil
			}
			return pushConf.IosPushConf
		}
	}
	//load from db
	//load iOS certificate
	iosPushConf := initIosPushConf(ctx, appkey, packageName)
	if pushConf != nil {
		pushConf.IosPushConf = iosPushConf
	} else {
		pushConf = &PushConfWraper{}
		pushConf.IosPushConf = iosPushConf
		pushConfCache.Add(key, pushConf)
	}
	if pushConf.IosPushConf == noExistIosPushConf {
		return nil
	}
	return pushConf.IosPushConf
}

var noExistIosPushConf *IosPushConf = &IosPushConf{}

func initIosPushConf(ctx context.Context, appkey, packageName string) *IosPushConf {
	iosDao := dbs.IosCertificateDao{}
	iosDb, err := iosDao.FindByPackage(appkey, packageName)
	if err == nil {
		iosPushConf := &IosPushConf{
			Package: packageName,
		}
		if len(iosDb.Certificate) > 0 {
			cert, err := certificate.FromP12Bytes(iosDb.Certificate, iosDb.CertPwd)
			if err == nil {
				if iosDb.IsProduct > 0 {
					iosPushConf.ApnsClient = apns2.NewClient(cert).Production()
				} else {
					iosPushConf.ApnsClient = apns2.NewClient(cert).Development()
				}
				return iosPushConf
			} else {
				logs.WithContext(ctx).Errorf("init ios certificate failed. %v", err)
			}
		}
		//voip cert
		if len(iosDb.VoipCert) > 0 {
			voipCert, err := certificate.FromP12Bytes(iosDb.VoipCert, iosDb.VoipCertPwd)
			if err == nil {
				if iosDb.IsProduct > 0 {
					iosPushConf.ApnsVoipClient = apns2.NewClient(voipCert).Production()
				} else {
					iosPushConf.ApnsVoipClient = apns2.NewClient(voipCert).Development()
				}
			} else {
				logs.WithContext(ctx).Errorf("init ios voip certificate failed. %v", err)
			}
		}
		return iosPushConf
	} else {
		logs.WithContext(ctx).Errorf("qry ios push conf failed.app_key:%s\tpackage:%s\terr:%v", appkey, packageName, err)
	}
	return noExistIosPushConf
}

func GetAndroidPushConf(ctx context.Context, appkey, packageName string) *AndroidPushConf {
	key := strings.Join([]string{appkey, packageName}, "_")
	if obj, exist := pushConfCache.Get(key); exist {
		pushConf := obj.(*PushConfWraper)
		if pushConf.AndroidPushConf != nil {
			if pushConf.AndroidPushConf == noExistAndroidPushConf {
				return nil
			}
			return pushConf.AndroidPushConf
		}
	}
	lock := pushConfLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	var pushConf *PushConfWraper
	if obj, exist := pushConfCache.Get(key); exist {
		pushConf = obj.(*PushConfWraper)
		if pushConf.AndroidPushConf != nil {
			if pushConf.AndroidPushConf == noExistAndroidPushConf {
				return nil
			}
			return pushConf.AndroidPushConf
		}
	}
	//load from db
	androidPushConf := initAndroidPushConf(ctx, appkey, packageName)
	if pushConf != nil {
		pushConf.AndroidPushConf = androidPushConf
	} else {
		pushConf = &PushConfWraper{}
		pushConf.AndroidPushConf = androidPushConf
		pushConfCache.Add(key, pushConf)
	}
	if pushConf.AndroidPushConf == noExistAndroidPushConf {
		return nil
	}
	return pushConf.AndroidPushConf
}

var noExistAndroidPushConf *AndroidPushConf = &AndroidPushConf{}

func initAndroidPushConf(ctx context.Context, appkey, packageName string) *AndroidPushConf {
	androidPushDao := dbs.AndroidPushConfDao{}
	pushConfs, err := androidPushDao.FindByPackage(appkey, packageName)
	if err == nil && len(pushConfs) > 0 {
		androidPushConf := &AndroidPushConf{
			Package: packageName,
		}
		for _, dbPushConf := range pushConfs {
			switch dbPushConf.PushChannel {
			case string(commonservices.PushChannel_Huawei):
				var pushConf = &commonservices.HuaweiPushConf{}
				err = json.Unmarshal([]byte(dbPushConf.PushConf), pushConf)
				if err == nil && pushConf.Valid() {
					hwPushClient, err := hwpush.NewHwPushClient(pushConf.AppId, pushConf.AppSecret)
					if err == nil {
						androidPushConf.HwPushClient = hwPushClient
					} else {
						logs.WithContext(ctx).Errorf("init huawei push client failed. %v", err)
					}
				} else {
					logs.WithContext(ctx).Errorf("huawei push conf is illegal %v", err)
				}
			case string(commonservices.PushChannel_Xiaomi):
				var pushConf = &commonservices.XiaomiPushConf{}
				err = json.Unmarshal([]byte(dbPushConf.PushConf), pushConf)
				if err == nil && pushConf.Valid() {
					androidPushConf.XiaomiPushClient = xiaomipush.NewXiaomiPushClient(pushConf.AppSecret)
				} else {
					logs.WithContext(ctx).Errorf("xiaomi push conf is illegal %v", err)
				}
			case string(commonservices.PushChannel_OPPO):
				var pushConf = &commonservices.OppoPushConf{}
				err = json.Unmarshal([]byte(dbPushConf.PushConf), pushConf)
				if err == nil && pushConf.Valid() {
					androidPushConf.OppoPushClient = oppopush.NewOppoPushClient(pushConf.AppKey, pushConf.MasterSecret)
				} else {
					logs.WithContext(ctx).Errorf("oppo push conf is illegal %v", err)
				}
			case string(commonservices.PushChannel_VIVO):
				var pushConf = &commonservices.VivoPushConf{}
				err = json.Unmarshal([]byte(dbPushConf.PushConf), pushConf)
				if err == nil && pushConf.Valid() {
					androidPushConf.VivoPushClient = vivopush.NewVivoPushClient(pushConf.AppId, pushConf.AppKey, pushConf.AppSecret)
				} else {
					logs.WithContext(ctx).Errorf("vivo push conf is illegal %v", err)
				}
			case string(commonservices.PushChannel_Jpush):
				var pushConf = &commonservices.JPushConf{}
				err = json.Unmarshal([]byte(dbPushConf.PushConf), pushConf)
				if err == nil && pushConf.Valid() {
					androidPushConf.JpushClient = jpush.NewJpushClient(pushConf.AppKey, pushConf.MasterSecret)
				} else {
					logs.WithContext(ctx).Errorf("jiguang push conf is illegal %v", err)
				}
			case string(commonservices.PushChannel_FCM):
				if len(dbPushConf.PushExt) > 0 {
					fcmClient, err := fcmpush.NewFcmPushClient(dbPushConf.PushExt)
					if err == nil {
						androidPushConf.FcmPushClient = fcmClient
					} else {
						logs.WithContext(ctx).Errorf("fcm conf is illegal %v", err)
					}
				}
			default:
			}
		}
		return androidPushConf
	}
	return noExistAndroidPushConf
}
