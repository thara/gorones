package apu

import (
	"context"
	"testing"
	"time"
)

func Test_divider(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	d := runDivider(ctx, uint(5))
	ch := d.output()

	for i := 0; i < 5; i++ {
		t.Logf("clock %d", i)
		d.clock()

		assertNotRecv(t, ch, "i=%d", i)
	}

	d.clock()
	assertRecv(t, ch)
}
