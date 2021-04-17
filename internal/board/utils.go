package board

// utils.go is a file containing utility functions for board.

// GetKing returns the king number for player
func GetKing(p1 bool) int8 {
	if p1 {
		return 28
	} else {
		return 4
	}
}

// GetQueen returns the queen number
func GetQueen(p1 bool) int8 {
	if p1 {
		return 27
	} else {
		return 3
	}
}

// GetBishops returns bishop numbers
func GetBishops(p1 bool) [2]int8 {
	fir, sec := int8(2), int8(5)
	if p1 {
		return [2]int8{fir + 24, sec + 24}
	} else {
		return [2]int8{fir, sec}
	}
}

// GetKnights returns knight numbers
func GetKnights(p1 bool) [2]int8 {
	fir, sec := int8(1), int8(6)
	if p1 {
		return [2]int8{fir + 24, sec + 24}
	} else {
		return [2]int8{fir, sec}
	}
}

// GetRooks returns rook numbers
func GetRooks(p1 bool) [2]int8 {
	fir, sec := int8(0), int8(7)
	if p1 {
		return [2]int8{fir + 24, sec + 24}
	} else {
		return [2]int8{fir, sec}
	}
}

// GetInversePlayer returns the opposite player
func GetInversePlayer(p1 bool) bool {
	return !p1
}

// GetPawnRow returns the pawn row for the player
func GetPawnRow(p1 bool) int8 {
	if p1 {
		return 6
	} else {
		return 1
	}
}

// GetStartRow returns the start row(the row which has queen, bishop, ...) for the player
func GetStartRow(p1 bool) int8 {
	if p1 {
		return 7
	} else {
		return 0
	}

	return -1
}

// GetRange returns an array of possible ids for a player's pieces..
func GetRange(p1 bool) [16]int8 {
	start := int8(0)
	if p1 {
		start += 16
	}

	arr := [16]int8{}
	for i := int8(0); i < 16; i++ {
		arr[i] = i + start
	}

	return arr
}

// GetRangeStart returns range(piece ids) for the start row(row that contains bishops, knights, ...)
func GetRangeStart(p1 bool) [8]int8 {
	start := int8(0)
	if p1 {
		start += 24
	}

	arr := [8]int8{}
	for i := int8(0); i < int8(len(arr)); i++ {
		arr[i] = i + start
	}

	return arr
}

// GetRangeStart returns range(piece ids) for the pawn row
func GetRangePawn(p1 bool) [8]int8 {
	start := int8(8)
	if p1 {
		start += 8
	}

	arr := [8]int8{}
	for i := int8(0); i < 8; i++ {
		arr[i] = i + start
	}

	return arr
}

// IsIDValid returns true if the id is valid
func IsIDValid(id int8) bool {
	return id <= 31 && id >= 0
}

// GetEighthRank returns value of y whenever pawn needs to promote.
func GetEighthRank(p1 bool) int8 {
	if p1 {
		return 0
	} else {
		return 7
	}
}

// BelongsTo returns true if the id belongs the player number
func BelongsTo(id int8, p1 bool) bool {
	if p1 {
		return id >= 16 && id < 32
	} else {
		return id >= 0 && id < 16
	}
}
