package cache

import (
	"github.com/golang/freetype/raster"
)

// SpanCache stores spans created by Rasterizer
type SpanCache []raster.Span

func (sc SpanCache) Len() int {
	return len(sc)
}

func (sc SpanCache) Less(i, j int) bool {
	a, b := sc[i], sc[j]
	return a.Y < b.Y || a.X0 < b.X0 || a.X1 < b.X1
}

func (sc SpanCache) Swap(i, j int) {
	sc[i], sc[j] = sc[j], sc[i]
}

// Paint adds new spans to SpanCache
func (sc *SpanCache) Paint(ss []raster.Span, done bool) {
	*sc = append(*sc, ss...)
}

// UsePainter makes a single Paint call with all collected spans
func (sc SpanCache) UsePainter(p raster.Painter) {
	p.Paint(sc, true)
}

// OptimizeSpans connects bounds when X0 from current Span and X1 from previous one are equal.
// Y and Alpha must also be equal.
func OptimizeSpans(sc SpanCache) SpanCache {
	for i := 1; i < len(sc); i++ {
		ps, cs := sc[i-1], sc[i]
		if ps.Y == cs.Y && ps.Alpha == cs.Alpha && ps.X1 == cs.X0 {
			sc[i-1].X1 = cs.X1
			sc = sc[:i+copy(sc[i:], sc[i+1:])]
		}
	}
	return sc
}
