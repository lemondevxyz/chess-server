package board

// absolute value of i
func abs(i int) int {
	if i < 0 {
		return i * -1
	}

	return i
}

// Point represents the piece's position in the Board.
type Point struct {
	X int
	Y int
}

// Swap swaps X and Y
func Swap(p Point) Point {
	return Point{X: p.Y, Y: p.X}
}

// Forward allows dst to only move forward, where dst is greater than or equal to src. If src == dst, it returns false
func Forward(src, dst Point) bool {
	i, j := src.X-dst.X, src.Y-dst.Y
	if i >= 0 && j >= 0 {
		// we didnt move
		if i == 0 && j == i {
			return false
		}
		return true
	}

	return false
}

// Backward allows dst to only move backward, where dst is less than or equal to src. If src == dst, it returns false
func Backward(src, dst Point) bool {
	i, j := src.X-dst.X, src.Y-dst.Y
	if i <= 0 && j <= 0 {
		// we didnt move
		if i == 0 && j == i {
			return false
		}
		return true
	}

	return false
}

// Within allows dst to be inside area
func Within(area, src, dst Point) bool {
	i, j := abs(src.X-dst.X), abs(src.Y-dst.Y)
	if i == area.X && j == area.Y {
		return true
	}

	return false
}

// Horizontal allows dst to be one of:
// Up,
// Down
func Horizontal(src, dst Point) bool {
	if src.Y == dst.Y {
		if src.X != dst.X {
			return true
		}
	}

	return false
}

// Vertical allows src to be one of:
// Right,
// Left
func Vertical(src, dst Point) bool {
	if src.X == dst.X {
		if src.Y != dst.Y {
			return true
		}
	}

	return false
}

// Diagonal allows dst to be one of:
// Up Right,
// Up Left,
// Down Right,
// Down Left
func Diagonal(src, dst Point) bool {
	i, j := abs(src.X-dst.X), abs(src.Y-dst.Y)
	if i == j {
		return true
	}

	return false
}

// Square allows dst to be one of:
// Up Left,
// Up,
// Up Right,
// Right,
// Left,
// Down Left,
// Down,
// Down Right
func Square(src, dst Point) bool {
	corner := Point{
		X: 1,
		Y: 1,
	}
	area := Point{X: 1, Y: 0}

	if src.X == dst.X && src.Y == dst.Y {
		return false
	}

	return Within(corner, src, dst) || Within(area, src, dst) || Within(Swap(area), src, dst)
}
