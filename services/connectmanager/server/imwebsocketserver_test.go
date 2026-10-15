package server

import (
	"errors"
	"net"
	"sync"
	"testing"
	"time"
)

func TestImWebsocketChildStopConcurrentSafe(t *testing.T) {
	child := &ImWebsocketChild{stopChan: make(chan struct{})}
	child.isActive.Store(true)

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			child.Stop()
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("concurrent Stop calls did not return")
	}
	if child.isActive.Load() {
		t.Fatal("child is still active after Stop")
	}
	select {
	case <-child.stopChan:
	default:
		t.Fatal("stop channel was not closed")
	}

	child.Stop()
}

func TestWsHandleContextIsActiveHandlesNilChild(t *testing.T) {
	ctx := &WsHandleContextImpl{}
	if ctx.IsActive() {
		t.Fatal("context without a child must be inactive")
	}
}

func TestWsHandleContextWriteSkipsInactiveConnection(t *testing.T) {
	child := &ImWebsocketChild{stopChan: make(chan struct{})}
	ctx := &WsHandleContextImpl{wsChild: child}

	// The inactive check must happen before message conversion, locking, or
	// dereferencing the connection.
	ctx.Write(nil)
}

func TestIsExpectedClosedConnError(t *testing.T) {
	if !isExpectedClosedConnError(net.ErrClosed) {
		t.Fatal("net.ErrClosed must be treated as an expected disconnect")
	}
	if !isExpectedClosedConnError(errors.New("write tcp: use of closed network connection")) {
		t.Fatal("closed network connection must be treated as an expected disconnect")
	}
	if isExpectedClosedConnError(errors.New("i/o timeout")) {
		t.Fatal("an active connection timeout must remain observable")
	}
}
