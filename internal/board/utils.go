package board

// utils.go is a file containing utility functions for board that do not depend on it.
// GetKing returns the king number for player
func GetKing(player uint8) int {
	if player == 1 {
		return 28
	} else if player == 2 {
		return 4
	}

	return -1
}

// GetQueen returns the queen number
func GetQueen(player uint8) int {
	if player == 1 {
		return 27
	} else if player == 2 {
		return 3
	}

	return -1
}

// GetBishops returns bishop numbers
func GetBishops(player uint8) (arr [2]int) {
	fir, sec := 2, 5
	if player == 1 {
		arr = [2]int{fir + 24, sec + 24}
	} else if player == 2 {
		arr = [2]int{fir, sec}
	}

	return arr
}

// GetKnights returns knight numbers
func GetKnights(player uint8) (arr [2]int) {
	fir, sec := 1, 6
	if player == 1 {
		arr = [2]int{fir + 24, sec + 24}
	} else if player == 2 {
		arr = [2]int{fir, sec}
	}

	return arr
}

// GetRooks returns rook numbers
func GetRooks(player uint8) (arr [2]int) {
	fir, sec := 0, 7
	if player == 1 {
		arr = [2]int{fir + 24, sec + 24}
	} else if player == 2 {
		arr = [2]int{fir, sec}
	}

	return arr
}

// GetInversePlayer returns the opposite player
func GetInversePlayer(player uint8) uint8 {
	if player == 1 {
		return 2
	} else if player == 2 {
		return 1
	}

	return 0
}

// GetPawnRow returns the pawn row for the player
func GetPawnRow(player uint8) int8 {
	if player == 1 {
		return 6
	} else if player == 2 {
		return 1
	}

	return -1
}

// GetStartRow returns the start row(the row which has queen, bishop, ...) for the player
func GetStartRow(player uint8) int8 {
	if player == 1 {
		return 7
	} else if player == 2 {
		return 0
	}

	return -1
}

// GetRange returns an array of possible ids for a player's pieces..
func GetRange(player uint8) [16]int {
	start := 0
	if player == 1 {
		start += 16
	}

	arr := [16]int{}
	for i := 0; i < 16; i++ {
		arr[i] = i + start
	}

	return arr
}

// GetRangeStart returns range(piece ids) for the start row(row that contains bishops, knights, ...)
func GetRangeStart(player uint8) [8]int {
	start := 0
	if player == 1 {
		start += 24
	}

	arr := [8]int{}
	for i := 0; i < len(arr); i++ {
		arr[i] = i + start
	}

	return arr
}

// GetRangeStart returns range(piece ids) for the pawn row
func GetRangePawn(player uint8) [8]int {
	start := 8
	if player == 1 {
		start += 8
	}

	arr := [8]int{}
	for i := 0; i < len(arr); i++ {
		arr[i] = i + start
	}

	return arr
}

// BelongsTo returns true if the id belongs the player number
func BelongsTo(id int8, player uint8) bool {
	if player == 1 {
		return id >= 16 && id < 32
	} else if player == 2 {
		return id >= 0 && id < 16
	}

	return false
}
