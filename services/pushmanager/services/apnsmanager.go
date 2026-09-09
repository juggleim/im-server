package services

import (
	"context"
	"errors"
	"sync"

	"im-server/services/pushmanager/services/apnspush"
	"im-server/services/pushmanager/storages/dbs"

	"github.com/sideshow/apns2"
	"gorm.io/gorm"
)

var iosManagerState struct {
	sync.RWMutex
	manager *apnspush.Manager
}

// StartupIosPush starts P8 lifecycle management without reading credentials or
// accessing the database. Legacy P12 clients remain on pushconfservice's path.
func StartupIosPush() {
	iosManagerState.Lock()
	defer iosManagerState.Unlock()
	if iosManagerState.manager == nil {
		iosManagerState.manager = apnspush.NewManager(loadIosPushConfig)
	}
}

func ShutdownIosPush() {
	iosManagerState.Lock()
	defer iosManagerState.Unlock()
	if iosManagerState.manager != nil {
		iosManagerState.manager.Close()
		iosManagerState.manager = nil
	}
}

func SendIosPush(ctx context.Context, appkey, packageName string, notification *apns2.Notification, isVoip bool, ordinaryToken, voipToken string) (*apns2.Response, error) {
	iosManagerState.RLock()
	manager := iosManagerState.manager
	iosManagerState.RUnlock()
	if manager == nil {
		return nil, apnspush.ErrClosed
	}
	return manager.Send(ctx, appkey, packageName, notification, isVoip, ordinaryToken, voipToken)
}

func loadIosPushConfig(ctx context.Context, appkey, packageName string) (*apnspush.Config, error) {
	row, err := (dbs.IosCertificateDao{}).FindByPackageWithContext(ctx, appkey, packageName)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apnspush.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if row.AuthType != "p8" {
		return nil, apnspush.ErrNotFound
	}
	return &apnspush.Config{
		AppKey: row.AppKey, Package: row.Package,
		IsProduct: row.IsProduct, ConfigVersion: row.ConfigVersion,
		P8KeyID: row.P8KeyID, P8TeamID: row.P8TeamID,
		P8PrivateKey: row.P8PrivateKey,
	}, nil
}
