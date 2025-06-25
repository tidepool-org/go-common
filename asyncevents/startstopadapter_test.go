package asyncevents

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestBlockingStartStopAdapter_StartWhenAlreadyStarted(t *testing.T) {
	ctx := context.Background()
	a := NewBlockingStartStopAdapter(newTestConsumer(t, nil))
	if err := a.Start(ctx); err != nil {
		t.Fatalf("expected nil error, got %s", err)
	}

	got := a.Start(ctx)
	if got == nil {
		t.Errorf("expected error, got nil")
	} else if !strings.Contains(got.Error(), "it's already running") {
		t.Errorf("expected error to contain \"it's already running\", got %s", got)
	}
}

func TestBlockingStartStopAdapter_StopWhenNotStarted(t *testing.T) {
	ctx := context.Background()
	a := NewBlockingStartStopAdapter(newTestConsumer(t, nil))

	got := a.Stop(ctx)
	if got == nil {
		t.Errorf("expected error, got nil")
	} else if !strings.Contains(got.Error(), "it's not running") {
		t.Errorf("expected error to contain \"it's not running\", got %s", got)
	}
}

func TestBlockingStartStopAdapter_StopsWhenCanceled(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	t.Cleanup(cancel)
	run, started := runUntilCanceled(t)
	a := NewBlockingStartStopAdapter(newTestConsumer(t, run))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.Start(ctx); err != nil {
			t.Errorf("expected nil error from Start, got %s", err)
		}
	}()
	<-started

	got := a.Stop(ctx)
	if got != nil {
		t.Errorf("expected nil error from Stop, got %s", got)
	}
	wg.Wait()
}

func TestBlockingStartStopAdapter_ReturnsErrorWhenDeadlineExceeded(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Millisecond)
	t.Cleanup(cancel)
	run, started := runUntilCanceled(t)
	a := NewBlockingStartStopAdapter(newTestConsumer(t, run))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := a.Start(ctx)
		if err == nil || !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("expected deadline exceeded error, got %v", err)
		}
	}()
	<-started
	wg.Wait()
}

func TestNonBlockingStartStopAdapter_StartDoesntBlock(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	t.Cleanup(cancel)
	run, _ := runUntilCanceled(t)
	a := NewNonBlockingStartStopAdapter(newTestConsumer(t, run), nil)

	start := time.Now()
	if err := a.Start(ctx); err != nil {
		t.Errorf("expected nil error, got %s", err)
	}
	if time.Since(start) > 10*time.Millisecond {
		t.Errorf("expected start to return immediately, but it took longer than expected")
	}
}

func TestNonBlockingStartStopAdapter_StopsWhenCalled(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	t.Cleanup(cancel)
	run, started := runUntilCanceled(t)
	a := NewNonBlockingStartStopAdapter(newTestConsumer(t, run), nil)

	if err := a.Start(ctx); err != nil {
		t.Errorf("expected nil error, got %s", err)
	}
	<-started

	if err := a.Stop(ctx); err != nil {
		t.Errorf("expected nil error, got %s", err)
	}
}

func TestNonBlockingStartStopAdapter_UsesCallback(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	t.Cleanup(cancel)
	run, done := runReturningErrorAfter(t, time.Millisecond)
	cbErrs := []error{}
	cb := func(err error) {
		cbErrs = append(cbErrs, err)
	}
	a := NewNonBlockingStartStopAdapter(newTestConsumer(t, run), cb)

	if err := a.Start(ctx); err != nil {
		t.Errorf("expected nil error, got %s", err)
	}
	<-done
	if len(cbErrs) < 1 {
		t.Errorf("expected the callback to be called, but it wasn't")
		err := cbErrs[0]
		if !strings.Contains(err.Error(), "blowing up") {
			t.Errorf("expected callback error to contain \"blowing up\", got: %s", err)
		}
	}
}

func TestNonBlockingStartStopAdapter_RecoversPanicsWithCallback(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	t.Cleanup(cancel)
	run, done := runPanickingAfter(t, time.Millisecond)
	cbErrs := []error{}
	cb := func(err error) {
		cbErrs = append(cbErrs, err)
	}
	a := NewNonBlockingStartStopAdapter(newTestConsumer(t, run), cb)

	if err := a.Start(ctx); err != nil {
		t.Errorf("expected nil error, got %s", err)
	}
	<-done
	// It can take a little while for the panic to propagate even after done is closed.
	startedWaiting := time.Now()
	for len(cbErrs) == 0 && time.Since(startedWaiting) < time.Second {
		time.Sleep(time.Microsecond)
	}
	if len(cbErrs) < 1 {
		t.Errorf("expected the callback to be called, but it wasn't")
		err := cbErrs[0]
		if !strings.Contains(err.Error(), "panic in the disco") {
			t.Errorf("expected callback error to contain \"panic in the disco\", got: %s", err)
		}
	}
}

func TestNonBlockingStartStopAdapter_RecoversPanicsWithCallbackIncludingError(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	t.Cleanup(cancel)
	run, done := runPanickingAfterWithError(t, time.Millisecond)
	cbErrs := []error{}
	cb := func(err error) {
		cbErrs = append(cbErrs, err)
	}
	a := NewNonBlockingStartStopAdapter(newTestConsumer(t, run), cb)

	if err := a.Start(ctx); err != nil {
		t.Errorf("expected nil error, got %s", err)
	}
	<-done
	// It can take a little while for the panic to propagate even after done is closed.
	startedWaiting := time.Now()
	for len(cbErrs) == 0 && time.Since(startedWaiting) < time.Second {
		time.Sleep(time.Microsecond)
	}
	if len(cbErrs) < 1 {
		t.Errorf("expected the callback to be called, but it wasn't")
		err := cbErrs[0]
		if !strings.Contains(err.Error(), "dignified panic") {
			t.Errorf("expected callback error to contain \"dignified panic\", got: %s", err)
		}
	}
}

func TestNonBlockingStartStopAdapter_DoesntRecoversPanicsWithoutCallback(t *testing.T) {
	// This can't really be tested. Why? Because the panic occurs in a separate goroutine,
	// and there's no way for the test to recover that panic. It will just bubble up until
	// it crashes the entire test.
}

type testConsumer struct {
	run func(context.Context) error
	t   testing.TB
}

func newTestConsumer(t testing.TB, run func(context.Context) error) *testConsumer {
	if run == nil {
		run = runReturningImmediately
	}
	return &testConsumer{
		run: run,
		t:   t,
	}
}

func (c *testConsumer) Run(ctx context.Context) error {
	c.t.Helper()
	if c.run != nil {
		return c.run(ctx)
	}
	return nil
}

func runUntilCanceled(t testing.TB) (func(ctx context.Context) error, <-chan struct{}) {
	t.Helper()
	started := make(chan struct{})
	return func(ctx context.Context) error {
		t.Helper()
		close(started)
		<-ctx.Done()
		if err := ctx.Err(); !errors.Is(err, context.Canceled) {
			return err
		}
		return nil
	}, started
}

func runReturningImmediately(ctx context.Context) error {
	return nil
}

func runReturningErrorAfter(t testing.TB, d time.Duration) (func(ctx context.Context) error, <-chan struct{}) {
	t.Helper()
	started := make(chan struct{})
	return func(ctx context.Context) error {
		t.Helper()
		defer close(started)
		select {
		case <-ctx.Done():
			if err := ctx.Err(); !errors.Is(err, context.Canceled) {
				return err
			}
		case <-time.After(d):
			return fmt.Errorf("blowing up")
		}
		return nil
	}, started
}

func runPanickingAfter(t testing.TB, d time.Duration) (func(ctx context.Context) error, <-chan struct{}) {
	t.Helper()
	started := make(chan struct{})
	return func(ctx context.Context) error {
		t.Helper()
		defer close(started)
		select {
		case <-ctx.Done():
			if err := ctx.Err(); !errors.Is(err, context.Canceled) {
				return err
			}
		case <-time.After(d):
			panic("panic in the disco!")
		}
		return nil
	}, started
}

func runPanickingAfterWithError(t testing.TB, d time.Duration) (func(ctx context.Context) error, <-chan struct{}) {
	t.Helper()
	started := make(chan struct{})
	return func(ctx context.Context) error {
		t.Helper()
		defer close(started)
		select {
		case <-ctx.Done():
			if err := ctx.Err(); !errors.Is(err, context.Canceled) {
				return err
			}
		case <-time.After(d):
			panic(fmt.Errorf("a dignified panic"))
		}
		return nil
	}, started
}
