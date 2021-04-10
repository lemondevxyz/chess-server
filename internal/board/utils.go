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

// GetInversePlayer returns the opposite player
func GetInversePlayer(player uint8) uint8 {
	if player == 1 {
		return 2
	} else if player == 2 {
		return 1
	}

	return 0
}

// BelongsTo returns if piece id to player
func BelongsTo(id int, player uint8) bool {
	if player == 1 {
		return id >= 16 && id < 32
	} else if player == 2 {
		return id < 16 && id >= 0
	}

	return false
}
