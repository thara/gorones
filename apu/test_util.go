package apu

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
)

func assertRecv[T any](t *testing.T, ch <-chan T, msgAndArgs ...interface{}) {
	t.Helper()

	select {
	case <-time.After(300 * time.Millisecond):
		assert.Fail(t, "should receive", msgAndArgs...)
	case <-ch:
	}
}

func assertRecvValue[T any](t *testing.T, ch <-chan T, expected T, msgAndArgs ...interface{}) {
	t.Helper()

	select {
	case <-time.After(300 * time.Millisecond):
		assert.Fail(t, "should receive", msgAndArgs...)
	case v := <-ch:
		assert.EqualValues(t, expected, v)
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

func clock[T constraints.Unsigned](t *testing.T, d *divider[T], n int) {
	t.Helper()

	for i := 0; i < n; i++ {
		d.clock()
	}
}
