package utils

import "math"

func orient(a, b, c LocalPos) float32 {
	return (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X)
}

func onSegment(a, b, p LocalPos) bool {
	return float64(p.X) >= math.Min(float64(a.X), float64(b.X)) && float64(p.X) <= math.Max(float64(a.X), float64(b.X)) &&
		float64(p.Y) >= math.Min(float64(a.Y), float64(b.Y)) && float64(p.Y) <= math.Max(float64(a.Y), float64(b.Y))
}

func segmentsIntersect(p1, p2, q1, q2 LocalPos) bool {
	o1 := orient(p1, p2, q1)
	o2 := orient(p1, p2, q2)
	o3 := orient(q1, q2, p1)
	o4 := orient(q1, q2, p2)

	if (o1 > 0 && o2 < 0 || o1 < 0 && o2 > 0) && (o3 > 0 && o4 < 0 || o3 < 0 && o4 > 0) {
		return true
	}
	// collinear cases (touching)
	if o1 == 0 && onSegment(p1, p2, q1) {
		return true
	}
	if o2 == 0 && onSegment(p1, p2, q2) {
		return true
	}
	if o3 == 0 && onSegment(q1, q2, p1) {
		return true
	}
	if o4 == 0 && onSegment(q1, q2, p2) {
		return true
	}
	return false
}

func SegmentIntersectsRect(a, b LocalPos, r LocalRect) bool {
	// rectangle corners
	tl := r.NW
	tr := LocalPos{X: r.SE.X, Y: r.NW.Y}
	bl := LocalPos{X: r.NW.X, Y: r.SE.Y}
	br := r.SE

	// check segment vs each rect edge
	if segmentsIntersect(a, b, tl, tr) {
		return true
	} // top
	if segmentsIntersect(a, b, tr, br) {
		return true
	} // right
	if segmentsIntersect(a, b, br, bl) {
		return true
	} // bottom
	if segmentsIntersect(a, b, bl, tl) {
		return true
	} // left

	return false
}
