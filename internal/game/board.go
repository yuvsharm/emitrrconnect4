package game

const Rows = 6
const Cols = 7

func NewBoard() [][]string {
	board := make([][]string, Rows)
	for i := range board {
		board[i] = make([]string, Cols)
	}
	return board
}

func Drop(board [][]string, col int, symbol string) bool {
	for r := Rows - 1; r >= 0; r-- {
		if board[r][col] == "" {
			board[r][col] = symbol
			return true
		}
	}
	return false
}
func IsBoardFull(board [][]string) bool {
	for r := 0; r < Rows; r++ {
		for c := 0; c < Cols; c++ {
			if board[r][c] == "" {
				return false
			}
		}
	}
	return true
}

