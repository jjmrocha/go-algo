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
		time.Sleep(10 * time.Millisecond)
		return value, err
	}
}

func TestAsyncWithWait(t *testing.T) {
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
			result, err := f.Await(ctx)
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
	// given
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Millisecond)
	defer cancel()
	// when
	_, result := AsyncWithContext(ctx, asyncFuncMaker(1, nil)).
		Await(ctx)
	// then
	assert.ErrorIs(t, result, context.DeadlineExceeded)
}

func TestAsyncWithCancel(t *testing.T) {
	// given
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	// when
	_, result := AsyncWithContext(ctx, asyncFuncMaker(1, nil)).
		Await(ctx)
	// then
	assert.ErrorIs(t, result, context.Canceled)
}

func TestAwaitCanBeCalledMultipleTimes(t *testing.T) {
	// given
	ctx := t.Context()
	f := AsyncWithContext(ctx, asyncFuncMaker(42, nil))
	// when: await the same future twice sequentially
	result1, err1 := f.Await(ctx)
	result2, err2 := f.Await(ctx)
	// then: both calls return the same result
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, 42, result1)
	assert.Equal(t, 42, result2)
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
	assert.ErrorIs(t, err, context.Canceled)
	// then: second Await with a fresh context still receives the result
	result, err := f.Await(t.Context())
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
			results[idx], errs[idx] = f.Await(ctx)
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
		// given: computation finishes well within the timeout
		f := Async(func() (int, error) {
			time.Sleep(10 * time.Millisecond)
			return 42, nil
		})
		// when
		result, err := f.AwaitWithTimeout(200 * time.Millisecond)
		// then
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("timeout exceeded", func(t *testing.T) {
		// given: computation takes longer than the timeout
		f := Async(func() (int, error) {
			time.Sleep(200 * time.Millisecond)
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
		time.Sleep(20 * time.Millisecond)
		return expected, nil
	})
	// when
	result1, err1 := f.Await(t.Context())
	result2, err2 := f.Await(t.Context())
	// then
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, expected, result1)
	assert.Equal(t, expected, result2)
}
