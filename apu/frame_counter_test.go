package apu

import (
	"context"
	"testing"
	"time"
)

func Test_frameCounter_step4(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	d := runDivider(ctx, 2)

	fc := runFrameCounter(ctx, d)

	clock(t, d, 3)
	assertRecv(t, fc.quarterFrame(), "step 1")
	assertNotRecv(t, fc.halfFrame(), "step 1")

	clock(t, d, 3)
	assertRecv(t, fc.quarterFrame(), "step 2")
	assertRecv(t, fc.halfFrame(), "step 2")

	clock(t, d, 3)
	assertRecv(t, fc.quarterFrame(), "step 3")
	assertNotRecv(t, fc.halfFrame(), "step 3")

	clock(t, d, 3)
	assertRecv(t, fc.quarterFrame(), "step 4")
	assertRecv(t, fc.halfFrame(), "step 4")
}

func Test_frameCounter_step4_frameInterrupt(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	d := runDivider(ctx, 2)

	fc := runFrameCounter(ctx, d)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-fc.quarterFrame():
			case <-fc.halfFrame():
			}
		}
	}()

	clock(t, d, 3)
	assertNotRecv(t, fc.frameInterrupt(), "step 1")
	clock(t, d, 3)
	assertNotRecv(t, fc.frameInterrupt(), "step 2")
	clock(t, d, 3)
	assertNotRecv(t, fc.frameInterrupt(), "step 3")
	clock(t, d, 3)
	assertRecvValue(t, fc.frameInterrupt(), true, "step 4")

	fc.update(0b01000000)
	assertRecvValue(t, fc.frameInterrupt(), false, "step 4")

	clock(t, d, 3)
	assertNotRecv(t, fc.frameInterrupt(), "step 1")
	clock(t, d, 3)
	assertNotRecv(t, fc.frameInterrupt(), "step 2")
	clock(t, d, 3)
	assertNotRecv(t, fc.frameInterrupt(), "step 3")
	clock(t, d, 3)
	assertNotRecv(t, fc.frameInterrupt(), "step 4")
}

func Test_frameCounter_step5(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	d := runDivider(ctx, 2)

	fc := runFrameCounter(ctx, d)
	fc.update(0b10000000)

	clock(t, d, 3)
	assertRecv(t, fc.quarterFrame(), "step 1")
	assertNotRecv(t, fc.halfFrame(), "step 1")

	clock(t, d, 3)
	assertRecv(t, fc.quarterFrame(), "step 2")
	assertRecv(t, fc.halfFrame(), "step 2")

	clock(t, d, 3)
	assertRecv(t, fc.quarterFrame(), "step 3")
	assertNotRecv(t, fc.halfFrame(), "step 3")

	clock(t, d, 3)
	assertNotRecv(t, fc.quarterFrame(), "step 4")
	assertNotRecv(t, fc.halfFrame(), "step 4")

	clock(t, d, 3)
	assertRecv(t, fc.quarterFrame(), "step 5")
	assertRecv(t, fc.halfFrame(), "step 5")

	clock(t, d, 3)
	assertRecv(t, fc.quarterFrame(), "step 1")
	assertNotRecv(t, fc.halfFrame(), "step 1")
}
