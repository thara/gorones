package apu

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func assertRecv[T any](t *testing.T, ch <-chan T, msgAndArgs ...interface{}) {
	t.Helper()

	select {
	case <-time.After(300 * time.Millisecond):
		assert.Fail(t, "should receive", msgAndArgs...)
	case <-ch:
	}
}

func assertNotRecv[T any](t *testing.T, ch <-chan T, msgAndArgs ...interface{}) {
	t.Helper()

	select {
	case <-time.After(200 * time.Millisecond):
	case <-ch:
		assert.Fail(t, "should not receive", msgAndArgs...)
	}
}
