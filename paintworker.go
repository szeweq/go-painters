package painters

import (
	"github.com/golang/freetype/raster"
	"sync"
)

type (
	PaintWorker struct {
		p    raster.Painter
		tnum int
		sc   chan PaintState
		wg   sync.WaitGroup
	}
	PaintState struct {
		ss   []raster.Span
		done bool
	}
)

func NewPaintWorker(p raster.Painter, n int) *PaintWorker {
	pw := &PaintWorker{p: p, tnum: n}
	pw.start()
	return pw
}

func (pw *PaintWorker) Paint(ss []raster.Span, done bool) {
	pw.sc <- PaintState{ss, done}
	if done {
		close(pw.sc)
		pw.wg.Wait()
	}
}
func (pw *PaintWorker) start() {
	for i := 0; i < pw.tnum; i++ {
		go pw.collect()
	}
}
func (pw *PaintWorker) collect() {
	pw.wg.Add(1)
	for ps := range pw.sc {
		pw.p.Paint(ps.ss, ps.done)
	}
	pw.wg.Done()
}
