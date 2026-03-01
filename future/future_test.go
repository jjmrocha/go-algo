package future

import (
	"context"
	"errors"
	"sync"
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
					t.Errorf("expect no error, got %v", err)
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
		t.Errorf("expected DeadlineExceeded, got %v", err)
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
		t.Errorf("expected Canceled, got %v", err)
	}
}

func TestAwaitCanBeCalledMultipleTimes(t *testing.T) {
	// given
	ctx := t.Context()
	f := AsyncWithContext(ctx, asyncFuncMaker(42, nil))
	// when: await the same future twice sequentially
	v1, err1 := f.Await(ctx)
	v2, err2 := f.Await(ctx)
	// then: both calls return the same result
	if err1 != nil {
		t.Fatalf("first Await: unexpected error %v", err1)
	}

	if err2 != nil {
		t.Fatalf("second Await: unexpected error %v", err2)
	}

	if v1 != 42 || v2 != 42 {
		t.Fatalf("expected 42/42, got %d/%d", v1, v2)
	}
}

func TestAwaitContextCancelledThenRetried(t *testing.T) {
	// given: a future that resolves after a short delay
	f := Async(func() (int, error) {
		time.Sleep(50 * time.Millisecond)
		return 7, nil
	})
	// when: first Await is cancelled immediately
	cancelledCtx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err := f.Await(cancelledCtx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected Canceled, got %v", err)
	}
	// then: second Await with a fresh context still receives the result
	got, err := f.Await(t.Context())
	if err != nil {
		t.Fatalf("retry Await: unexpected error %v", err)
	}

	if got != 7 {
		t.Fatalf("expected 7, got %d", got)
	}
}

func TestAwaitConcurrent(t *testing.T) {
	// given: many goroutines all awaiting the same future
	const n = 50
	ctx := t.Context()
	f := AsyncWithContext(ctx, asyncFuncMaker(99, nil))

	var wg sync.WaitGroup
	results := make([]int, n)
	errs := make([]error, n)

	for i := range n {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx], errs[idx] = f.Await(ctx)
		}(i)
	}
	wg.Wait()

	// then: every goroutine received the same result
	for i := range n {
		if errs[i] != nil {
			t.Errorf("goroutine %d: unexpected error %v", i, errs[i])
		}

		if results[i] != 99 {
			t.Errorf("goroutine %d: expected 99, got %d", i, results[i])
		}
	}
}

func TestAwaitWithTimeout_Success(t *testing.T) {
	// given: computation finishes well within the timeout
	f := Async(func() (int, error) {
		time.Sleep(10 * time.Millisecond)
		return 42, nil
	})
	// when
	got, err := f.AwaitWithTimeout(200 * time.Millisecond)
	// then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}
}

func TestAwaitWithTimeout_Timeout(t *testing.T) {
	// given: computation takes longer than the timeout
	f := Async(func() (int, error) {
		time.Sleep(200 * time.Millisecond)
		return 42, nil
	})
	// when
	_, err := f.AwaitWithTimeout(5 * time.Millisecond)
	// then
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}

func TestMultipleConsumers(t *testing.T) {
	const want = "done"
	// given
	f := Async(func() (string, error) {
		time.Sleep(20 * time.Millisecond)
		return want, nil
	})
	// when
	r1, err1 := f.Await(t.Context())
	r2, err2 := f.Await(t.Context())
	// then
	if err1 != nil {
		t.Fatalf("first Await: unexpected error %v", err1)
	}

	if err2 != nil {
		t.Fatalf("second Await: unexpected error %v", err2)
	}

	if r1 != want || r2 != want {
		t.Fatalf("expected '%s'/'%s', got '%s'/'%s'", want, want, r1, r2)
	}
}
