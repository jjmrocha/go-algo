package future

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func asyncFuncMaker[T any](value T, err error) func(context.Context) (T, error) {
	return func(_ context.Context) (T, error) {
		return value, err
	}
}

func TestAwaitWithContext(t *testing.T) {
	// given
	ctx := t.Context()
	testError := errors.New("test error")

	tests := []struct {
		name          string
		value         any
		error         error
		expectedValue any
		expectedError error
	}{
		{
			name:          "success - int",
			value:         1,
			error:         nil,
			expectedValue: 1,
			expectedError: nil,
		},
		{
			name:          "success - string",
			value:         "hello",
			error:         nil,
			expectedValue: "hello",
			expectedError: nil,
		},
		{
			name:          "success - error",
			value:         nil,
			error:         testError,
			expectedValue: nil,
			expectedError: testError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			f := AsyncWithContext(ctx, asyncFuncMaker(tt.value, tt.error))
			// when
			result, err := f.AwaitWithContext(ctx)
			// then
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedValue, result)
			}
		})
	}
}

func TestAsyncWithTimeout(t *testing.T) {
	// given: a provider that never completes, so the deadline is what fires
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Millisecond)
	defer cancel()
	block := make(chan struct{})
	defer close(block)
	f := AsyncWithContext(ctx, func(context.Context) (int, error) {
		<-block
		return 1, nil
	})
	// when
	_, result := f.AwaitWithContext(ctx)
	// then
	assert.ErrorIs(t, result, context.DeadlineExceeded)
}

func TestAsyncWithCancel(t *testing.T) {
	// given: ctx is cancelled up front and the provider never completes
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	block := make(chan struct{})
	defer close(block)
	f := AsyncWithContext(ctx, func(context.Context) (int, error) {
		<-block
		return 1, nil
	})
	// when
	_, result := f.AwaitWithContext(ctx)
	// then
	assert.ErrorIs(t, result, context.Canceled)
}

func TestAwaitCanBeCalledMultipleTimes(t *testing.T) {
	// given
	ctx := t.Context()
	f := AsyncWithContext(ctx, asyncFuncMaker(42, nil))
	// when: await the same future twice sequentially
	result1, err1 := f.AwaitWithContext(ctx)
	result2, err2 := f.AwaitWithContext(ctx)
	// then: both calls return the same result
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, 42, result1)
	assert.Equal(t, 42, result2)
}

func TestAwaitContextCancelledThenRetried(t *testing.T) {
	// given: a future held unresolved until release is closed
	release := make(chan struct{})
	f := Async(func() (int, error) {
		<-release
		return 7, nil
	})
	// when: the first await is cancelled while the future is still unresolved
	cancelledCtx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err := f.AwaitWithContext(cancelledCtx)
	assert.ErrorIs(t, err, context.Canceled)
	// then: once resolved, a second await with a live context receives the result
	close(release)
	result, err := f.AwaitWithContext(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, 7, result)
}

func TestAwaitConcurrent(t *testing.T) {
	// given: many goroutines all awaiting the same future
	const n = 50
	ctx := t.Context()
	f := AsyncWithContext(ctx, asyncFuncMaker(99, nil))
	var wg sync.WaitGroup
	results := make([]int, n)
	errs := make([]error, n)
	// when
	for i := range n {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx], errs[idx] = f.AwaitWithContext(ctx)
		}(i)
	}
	wg.Wait()
	// then: every goroutine received the same result
	for i := range n {
		assert.NoError(t, errs[i], "goroutine %d", i)
		assert.Equal(t, 99, results[i], "goroutine %d", i)
	}
}

func TestAwaitWithTimeout(t *testing.T) {
	t.Run("success within timeout", func(t *testing.T) {
		// given: the computation resolves on its own
		f := Async(func() (int, error) {
			return 42, nil
		})
		// when: a generous timeout that will not be reached
		result, err := f.AwaitWithTimeout(time.Second)
		// then
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("timeout exceeded", func(t *testing.T) {
		// given: a computation that never completes
		block := make(chan struct{})
		defer close(block)
		f := Async(func() (int, error) {
			<-block
			return 42, nil
		})
		// when
		_, result := f.AwaitWithTimeout(5 * time.Millisecond)
		// then
		assert.ErrorIs(t, result, context.DeadlineExceeded)
	})
}

func TestMultipleConsumers(t *testing.T) {
	// given
	expected := "done"
	f := Async(func() (string, error) {
		return expected, nil
	})
	// when
	result1, err1 := f.AwaitWithContext(t.Context())
	result2, err2 := f.AwaitWithContext(t.Context())
	// then
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, expected, result1)
	assert.Equal(t, expected, result2)
}

func TestAwait(t *testing.T) {
	t.Run("returns result when computation completes", func(t *testing.T) {
		// given
		f := Async(func() (int, error) {
			return 42, nil
		})
		// when
		result, err := f.Await()
		// then
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("returns error when computation fails", func(t *testing.T) {
		// given
		expectedErr := errors.New("compute failed")
		f := Async(func() (int, error) {
			return 0, expectedErr
		})
		// when
		result, err := f.Await()
		// then
		assert.ErrorIs(t, err, expectedErr)
		assert.Equal(t, 0, result)
	})
}

func TestDone(t *testing.T) {
	t.Run("false before computation completes", func(t *testing.T) {
		// given
		started := make(chan struct{})
		block := make(chan struct{})
		f := Async(func() (int, error) {
			close(started)
			<-block
			return 1, nil
		})
		<-started
		// when
		result := f.Done()
		close(block)
		// then
		assert.False(t, result)
	})

	t.Run("true after computation completes", func(t *testing.T) {
		// given
		f := Async(func() (int, error) {
			return 1, nil
		})
		_, _ = f.Await()
		// when
		result := f.Done()
		// then
		assert.True(t, result)
	})
}

func TestResolve(t *testing.T) {
	testError := errors.New("test error")

	t.Run("unblocks await with the resolved value", func(t *testing.T) {
		// given
		f := New[int]()
		assert.False(t, f.Done())
		// when
		f.Resolve(42, nil)
		// then
		result, err := f.Await()
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
		assert.True(t, f.Done())
	})

	t.Run("resolves with an error", func(t *testing.T) {
		// given
		f := New[string]()
		// when
		f.Resolve("", testError)
		// then
		result, err := f.Await()
		assert.ErrorIs(t, err, testError)
		assert.Empty(t, result)
	})

	t.Run("second resolve is a no-op", func(t *testing.T) {
		// given
		f := New[int]()
		f.Resolve(1, nil)
		// when
		f.Resolve(2, testError)
		// then
		result, err := f.Await()
		assert.NoError(t, err)
		assert.Equal(t, 1, result)
	})

	t.Run("concurrent resolves keep a single result", func(t *testing.T) {
		// given
		const goroutines = 100
		f := New[int]()
		var wg sync.WaitGroup
		// when
		for i := range goroutines {
			wg.Add(1)
			go func(v int) {
				defer wg.Done()
				f.Resolve(v, nil)
			}(i)
		}
		wg.Wait()
		// then
		result, err := f.Await()
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, result, 0)
		assert.Less(t, result, goroutines)
	})
}
