package singleflight

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// when
	result := New[string, int]()
	// then
	assert.NotNil(t, result)
}

func TestDo(t *testing.T) {
	t.Run("calls provider and returns result", func(t *testing.T) {
		// given
		sf := New[string, int]()
		// when
		result, err := sf.Do("key", func() (int, error) {
			return 42, nil
		}).AwaitWithContext(t.Context())
		// then
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("provider error is returned", func(t *testing.T) {
		// given
		sf := New[string, int]()
		providerErr := errors.New("provider failed")
		// when
		result, err := sf.Do("key", func() (int, error) {
			return 0, providerErr
		}).AwaitWithContext(t.Context())
		// then
		assert.ErrorIs(t, err, providerErr)
		assert.Equal(t, 0, result)
	})

	t.Run("concurrent calls for same key share future", func(t *testing.T) {
		// given
		sf := New[string, int]()
		var calls atomic.Int32
		const goroutines = 50
		unblock := make(chan struct{})
		provider := func() (int, error) {
			calls.Add(1)
			<-unblock
			return 99, nil
		}
		results := make([]int, goroutines)
		errs := make([]error, goroutines)
		var ready, done sync.WaitGroup
		ctx := t.Context()
		// when: all goroutines call Do before the provider completes
		for i := range goroutines {
			ready.Add(1)
			done.Add(1)
			go func(i int) {
				defer done.Done()
				f := sf.Do("key", provider)
				ready.Done()
				results[i], errs[i] = f.AwaitWithContext(ctx)
			}(i)
		}
		ready.Wait()
		close(unblock)
		done.Wait()
		// then: provider called once, all callers received the same result
		assert.Equal(t, int32(1), calls.Load())
		for i := range goroutines {
			assert.NoError(t, errs[i])
			assert.Equal(t, 99, results[i])
		}
	})

	t.Run("completed computation allows new call", func(t *testing.T) {
		// given
		sf := New[string, int]()
		var calls atomic.Int32
		provider := func() (int, error) {
			return int(calls.Add(1)), nil
		}
		// when: two sequential calls, each after the previous completes
		result1, err1 := sf.Do("key", provider).AwaitWithContext(t.Context())
		result2, err2 := sf.Do("key", provider).AwaitWithContext(t.Context())
		// then: provider called twice, each got its own result
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, int32(2), calls.Load())
		assert.Equal(t, 1, result1)
		assert.Equal(t, 2, result2)
	})

	t.Run("error does not prevent subsequent calls", func(t *testing.T) {
		// given
		sf := New[string, int]()
		var calls atomic.Int32
		providerErr := errors.New("transient error")
		p := func() (int, error) {
			if calls.Add(1) == 1 {
				return 0, providerErr
			}
			return 42, nil
		}
		// when: first call fails, second retries
		_, firstErr := sf.Do("key", p).AwaitWithContext(t.Context())
		result, err := sf.Do("key", p).AwaitWithContext(t.Context())
		// then
		assert.ErrorIs(t, firstErr, providerErr)
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
		assert.Equal(t, int32(2), calls.Load())
	})

	t.Run("different keys are independent", func(t *testing.T) {
		// given
		sf := New[string, int]()
		var calls atomic.Int32
		provider := func() (int, error) {
			return int(calls.Add(1)), nil
		}
		// when
		result1, err1 := sf.Do("a", provider).AwaitWithContext(t.Context())
		result2, err2 := sf.Do("b", provider).AwaitWithContext(t.Context())
		// then
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, int32(2), calls.Load())
		assert.Equal(t, 1, result1)
		assert.Equal(t, 2, result2)
	})
}
