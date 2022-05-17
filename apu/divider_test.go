package apu

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_divider(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	d := runDivider(ctx, 5)

	for i := 0; i < 5; i++ {
		t.Logf("clock %d", i)
		d.clock()

		select {
		case <-d.output():
			assert.Failf(t, "should not output", "i=%d", i)
		default:
		}
	}

	d.clock()

	select {
	case <-d.output():
	default:
		assert.Fail(t, "should output")
	}
}
