package asyncevents

import (
	"context"
	"fmt"
	"sync"
)

// Runner abstracts [SaramaConsumerGroupManager], which is intended to be extended with
// additional capabilities and behaviors.
type Runner interface {
	Run(context.Context) error
}

// BlockingStartStopAdapter for [SaramaConsumerGroupManager] to use Start and Stop methods.
//
// This adapter provides a base for more specific adaptation to adjust behavior to their
// needs. For example [NonBlockingStartStopAdapter] fits well with Uber's fx.Lifecycle,
// while this adapter is more adaptable to platform's event.Runner interface.
type BlockingStartStopAdapter struct {
	Runner Runner

	cancelMu   sync.Mutex
	cancelFunc context.CancelFunc
}

func NewBlockingStartStopAdapter(runner Runner) *BlockingStartStopAdapter {
	return &BlockingStartStopAdapter{
		Runner: runner,
	}
}

func (a *BlockingStartStopAdapter) Start(ctx context.Context) error {
	cancelCtx, err := a.init(ctx)
	if err != nil {
		return err
	}
	return a.Runner.Run(cancelCtx)
}

func (a *BlockingStartStopAdapter) init(ctx context.Context) (context.Context, error) {
	a.cancelMu.Lock()
	defer a.cancelMu.Unlock()

	if a.cancelFunc != nil {
		return nil, fmt.Errorf("can't start consumer, it's already running")
	}
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	a.cancelFunc = cancelFunc
	return cancelCtx, nil
}

func (a *BlockingStartStopAdapter) Stop(_ context.Context) error {
	a.cancelMu.Lock()
	defer a.cancelMu.Unlock()

	if a.cancelFunc == nil {
		return fmt.Errorf("can't stop consumer, it's not running")
	}

	a.cancelFunc()
	a.cancelFunc = nil

	return nil
}

// NonBlockingStartStopAdapter for [SaramaConsumerGroupManager] for non-blocking Start and
// Stop methods.
//
// To facilitate error reporting during non-blocking operation, a callback can be provided,
// which if defined, will be called with errors that cause a [SaramaConsumerGroupManager]'s
// Run method to return. In addition, when the callback is defined, panics from within Run
// are recovered, converted to errors, and passed to the callback before being discarded.
type NonBlockingStartStopAdapter struct {
	*BlockingStartStopAdapter

	onError func(error)
}

func NewNonBlockingStartStopAdapter(consumer Runner, onError func(error)) *NonBlockingStartStopAdapter {
	blocking := NewBlockingStartStopAdapter(consumer)
	return &NonBlockingStartStopAdapter{
		BlockingStartStopAdapter: blocking,
		onError:                  onError,
	}
}

func (a *NonBlockingStartStopAdapter) Start(ctx context.Context) error {
	go a.start(ctx)
	return nil
}

func (a *NonBlockingStartStopAdapter) start(ctx context.Context) {
	defer a.maybeRecover()

	if err := a.BlockingStartStopAdapter.Start(ctx); err != nil {
		if a.onError != nil {
			a.onError(err)
		}
	}
}

// maybeRecover uses a callback, if defined, to process recovered panics.
//
// If the callback isn't defined, the panic will be re-raised.
func (a *NonBlockingStartStopAdapter) maybeRecover() {
	if r := recover(); r != nil {
		if a.onError != nil {
			a.onError(a.wrapPanic(r))
		} else {
			panic(r)
		}
	}
}

// wrapPanic converts a non-error value to an error for passing to an on-error callback.
//
// Existing error values are wrapped, to allow later unwrapping.
func (a *NonBlockingStartStopAdapter) wrapPanic(r any) error {
	if err, ok := r.(error); ok {
		return fmt.Errorf("consumer panicked: %w", err)
	}
	return fmt.Errorf("consumer panicked: %s", r)
}
