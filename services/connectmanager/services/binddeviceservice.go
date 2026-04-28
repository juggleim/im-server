package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/connectmanager/dbs"
	"strings"
	"time"
)

var bindDeviceCache *caches.LruCache

func init() {
	bindDeviceCache = caches.NewLruCacheWithAddReadTimeout("binddevice_cache", 100000, nil, 10*time.Minute, 10*time.Minute)
}

func CheckBindDevice(appkey, userId, deviceId string) bool {
	key := strings.Join([]string{appkey, userId, deviceId}, "_")
	if obj, exist := bindDeviceCache.Get(key); exist {
		return obj.(bool)
	}

	l := userLocks.GetLocks(appkey, userId)
	l.Lock()
	defer l.Unlock()

	if obj, exist := bindDeviceCache.Get(key); exist {
		return obj.(bool)
	}

	dao := dbs.BindDeviceDao{}
	device, err := dao.FindByUserAndDevice(appkey, userId, deviceId)
	if err != nil || device == nil {
		bindDeviceCache.Add(key, false)
		return false
	}
	bindDeviceCache.Add(key, true)
	return true
}

func AddBindDevice(ctx context.Context, device *pbobjs.BindDevice) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)

	key := strings.Join([]string{appkey, userId, device.DeviceId}, "_")
	if !bindDeviceCache.Contains(key) {
		l := userLocks.GetLocks(appkey, userId)
		l.Lock()
		defer l.Unlock()
		if !bindDeviceCache.Contains(key) {
			dao := dbs.BindDeviceDao{}
			err := dao.Upsert(dbs.BindDeviceDao{
				AppKey:        appkey,
				UserId:        userId,
				DeviceId:      device.DeviceId,
				Platform:      device.Platform,
				DeviceCompany: device.DeviceCompany,
				DeviceModel:   device.DeviceModel,
			})
			if err != nil {
				logs.WithContext(ctx).Error(err.Error())
			}
			bindDeviceCache.Add(key, true)
		}
	}

	return errs.IMErrorCode_SUCCESS
}

func DelBindDevice(ctx context.Context, device *pbobjs.BindDevice) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	key := strings.Join([]string{appkey, userId, device.DeviceId}, "_")

	l := userLocks.GetLocks(appkey, userId)
	l.Lock()
	defer l.Unlock()

	dao := dbs.BindDeviceDao{}
	err := dao.DelByUserAndDevice(appkey, userId, device.DeviceId)
	if err != nil {
		logs.WithContext(ctx).Error(err.Error())
	}
	bindDeviceCache.Add(key, false)
	return errs.IMErrorCode_SUCCESS
}

func QryBindDevices(ctx context.Context) (errs.IMErrorCode, *pbobjs.BindDevices) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	bindDevices := &pbobjs.BindDevices{
		Devices: make([]*pbobjs.BindDevice, 0),
	}
	dao := dbs.BindDeviceDao{}
	devices, err := dao.FindByUserId(appkey, userId)
	if err == nil {
		for _, device := range devices {
			bindDevices.Devices = append(bindDevices.Devices, &pbobjs.BindDevice{
				DeviceId:      device.DeviceId,
				DeviceCompany: device.DeviceCompany,
				DeviceModel:   device.DeviceModel,
				Platform:      device.Platform,
				CreatedTime:   device.CreatedTime.UnixMilli(),
			})
		}
	}
	return errs.IMErrorCode_SUCCESS, bindDevices
}
