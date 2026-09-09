package push

import (
	"context"
	"errors"
	"testing"

	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/pushmanager/services"

	"github.com/sideshow/apns2"
)

type startupCheckingRegister struct {
	gmicro.IActorRegister
	t *testing.T
}

func (r startupCheckingRegister) RegisterActor(_ string, _ func() actorsystem.IUntypedActor) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := services.SendIosPush(ctx, "app", "com.example", &apns2.Notification{}, false, "ab", "")
	if !errors.Is(err, context.Canceled) {
		r.t.Fatalf("APNs must be initialized before exposing actors, got %v", err)
	}
}

func TestAPNSReadyBeforeActorsAreRegistered(t *testing.T) {
	services.ShutdownIosPush()
	defer services.ShutdownIosPush()
	manager := &PushManager{}
	manager.RegisterActors(startupCheckingRegister{t: t})
	manager.Startup(nil)
}
