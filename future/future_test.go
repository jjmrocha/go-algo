package future

import (
	"context"
	"errors"
	"testing"
	"time"
)

func asyncFuncMaker[T any](value T, err error) func(context.Context) (T, error) {
	return func(_ context.Context) (T, error) {
		time.Sleep(10 * time.Millisecond)
		return value, err
	}
}

func TestAsyncWithWait(t *testing.T) {
	// given
	ctx := t.Context()
	testError := errors.New("test error")

	tests := []struct {
		name        string
		value       any
		error       error
		expectValue any
		expectError error
	}{
		{
			name:        "success - int",
			value:       1,
			error:       nil,
			expectValue: 1,
			expectError: nil,
		},
		{
			name:        "success - string",
			value:       "hello",
			error:       nil,
			expectValue: "hello",
			expectError: nil,
		},
		{
			name:        "success - error",
			value:       nil,
			error:       testError,
			expectValue: nil,
			expectError: testError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// when
			got, err := AsyncWithContext(ctx, asyncFuncMaker(test.value, test.error)).
				Await(ctx)
			// then
			if test.expectError != nil {
				if !errors.Is(err, test.expectError) {
					t.Errorf("expect error %v, got %v", test.expectError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expect no error, got %v", got)
				}

				if test.expectValue != got {
					t.Errorf("expected %v, got %v", test.expectValue, got)
				}
			}
		})
	}
}

func TestAsyncWithTimeout(t *testing.T) {
	// given
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Millisecond)
	defer cancel()
	// when
	_, err := AsyncWithContext(ctx, asyncFuncMaker(1, nil)).
		Await(ctx)
	// then
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expect error, got nil")
	}
}

func TestAsyncWithCancel(t *testing.T) {
	// given
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	// when
	_, err := AsyncWithContext(ctx, asyncFuncMaker(1, nil)).
		Await(ctx)
	// then
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expect error, got nil")
	}
}

func TestAwaitCanOnlyBeCalledOnce(t *testing.T) {
	// given
	ctx := t.Context()
	f := AsyncWithContext(ctx, asyncFuncMaker(42, nil))
	f.Await(ctx) //nolint:errcheck // first call consumes the result
	// when
	_, err := f.Await(ctx)
	// then
	if !errors.Is(err, ErrAwaitAlreadyCalled) {
		t.Errorf("expected ErrAwaitAlreadyCalled, got %v", err)
	}
}
