package services

import (
	"context"
	"errors"
	"sync"
	"testing"

	"im-server/services/pushmanager/services/apnspush"

	"github.com/sideshow/apns2"
)

func TestAPNSLifecycle(t *testing.T) {
	ShutdownIosPush()
	if _, err := SendIosPush(context.Background(), "app", "com.example", &apns2.Notification{}, false, "ab", ""); !errors.Is(err, apnspush.ErrClosed) {
		t.Fatalf("send before startup: %v", err)
	}
	StartupIosPush()
	first := iosManagerState.manager
	StartupIosPush()
	if first != iosManagerState.manager {
		t.Fatal("repeated startup replaced manager")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := SendIosPush(ctx, "app", "com.example", &apns2.Notification{}, false, "ab", ""); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation was not preserved: %v", err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := SendIosPush(ctx, "app", "com.example", &apns2.Notification{}, false, "ab", "")
			if !errors.Is(err, context.Canceled) && !errors.Is(err, apnspush.ErrClosed) {
				t.Errorf("unexpected concurrent shutdown error: %v", err)
			}
		}()
	}
	ShutdownIosPush()
	ShutdownIosPush()
	wg.Wait()
	if _, err := SendIosPush(ctx, "app", "com.example", &apns2.Notification{}, false, "ab", ""); !errors.Is(err, apnspush.ErrClosed) {
		t.Fatalf("send after shutdown: %v", err)
	}
}
