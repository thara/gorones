package apu

import "context"

type pulse struct {
	clk         chan interface{}
	envelopeClk chan interface{}
	sweepClk    chan interface{}
	lcClk       chan interface{}

	upd chan regUpdate
	out chan uint8
}

type regUpdate struct {
	reg   uint8
	value uint8
}

func runPulse(ctx context.Context, cmp sweepComplement) *pulse {
	clk := make(chan interface{})
	envClk := make(chan interface{})
	sweepClk := make(chan interface{})
	lcClk := make(chan interface{})
	upd := make(chan regUpdate)
	out := make(chan uint8)

	go func() {
		defer close(clk)
		defer close(envClk)
		defer close(sweepClk)
		defer close(lcClk)
		defer close(upd)
		defer close(out)

		var (
			dutyCycle   uint8
			dutySeq     uint8
			timerPeriod uint16
			sweepMuted  bool
			t           uint16
		)
		var (
			env   = runEnvelope(ctx)
			sweep = runSweep(ctx, cmp)
			lc    = runLengthCounter(ctx)
			seq   = runSequencer(ctx)
		)
		for {
			select {
			case <-ctx.Done():

			case u := <-upd:
				switch u.reg {
				case 0:
					dutyCycle = u.value >> 6
				case 1:
					sweep.update(u.value)
				case 2:
					timerPeriod = (timerPeriod & 0xFF00) | uint16(u.value)
				case 3:
					lc.load(u.value)

					env.start()

					timerPeriod = (timerPeriod & 0x00FF) | uint16(u.value)<<8
					sweep.input(timerPeriod)

					seq.restart(0)
				}

			case <-clk:
				if 0 < t {
					t--
				} else {
					t = timerPeriod
				}

			case <-envClk:
				env.clock()
			case <-sweepClk:
				sweep.clock()
			case <-lcClk:
				lc.clock()

			case v := <-sweep.output():
				sweepMuted = v == 0
				seq.clock()

			case v := <-seq.output():
				dutySeq = v

			case v := <-env.output():
				d := dutyTable[dutyCycle][dutySeq]
				if lc.counter() == 0 || sweepMuted || d == 0 || t == 8 {
					continue
				}
				out <- v
			}
		}
	}()

	return &pulse{clk, envClk, sweepClk, lcClk, upd, out}
}

func (p *pulse) clock()               { p.clk <- struct{}{} }
func (p *pulse) clockEnvelope()       { p.envelopeClk <- struct{}{} }
func (p *pulse) clockSweep()          { p.sweepClk <- struct{}{} }
func (p *pulse) clockLengthCounter()  { p.lcClk <- struct{}{} }
func (p *pulse) update(reg, v uint8)  { p.upd <- regUpdate{reg, v} }
func (p *pulse) output() <-chan uint8 { return p.out }

var dutyTable [4][8]uint8 = [4][8]uint8{
	{0, 1, 0, 0, 0, 0, 0, 0}, // 12.5%
	{0, 1, 1, 0, 0, 0, 0, 0}, // 25%
	{0, 1, 1, 1, 1, 0, 0, 0}, // 50%
	{1, 0, 0, 1, 1, 1, 1, 1}, // 25% negated
}
