package game

type WinResult struct {
	Won      bool
	WinCells [][2]int
}

func CheckWin(b [][]string, s string) WinResult {
	for r := 0; r < Rows; r++ {
		for c := 0; c < Cols; c++ {

			// Horizontal →
			if c+3 < Cols &&
				b[r][c] == s &&
				b[r][c+1] == s &&
				b[r][c+2] == s &&
				b[r][c+3] == s {
				return WinResult{
					Won: true,
					WinCells: [][2]int{
						{r, c}, {r, c + 1}, {r, c + 2}, {r, c + 3},
					},
				}
			}

			// Vertical ↓
			if r+3 < Rows &&
				b[r][c] == s &&
				b[r+1][c] == s &&
				b[r+2][c] == s &&
				b[r+3][c] == s {
				return WinResult{
					Won: true,
					WinCells: [][2]int{
						{r, c}, {r + 1, c}, {r + 2, c}, {r + 3, c},
					},
				}
			}

			// Diagonal ↘
			if r+3 < Rows && c+3 < Cols &&
				b[r][c] == s &&
				b[r+1][c+1] == s &&
				b[r+2][c+2] == s &&
				b[r+3][c+3] == s {
				return WinResult{
					Won: true,
					WinCells: [][2]int{
						{r, c}, {r + 1, c + 1}, {r + 2, c + 2}, {r + 3, c + 3},
					},
				}
			}

			// Diagonal ↗
			if r-3 >= 0 && c+3 < Cols &&
				b[r][c] == s &&
				b[r-1][c+1] == s &&
				b[r-2][c+2] == s &&
				b[r-3][c+3] == s {
				return WinResult{
					Won: true,
					WinCells: [][2]int{
						{r, c}, {r - 1, c + 1}, {r - 2, c + 2}, {r - 3, c + 3},
					},
				}
			}
		}
	}

	return WinResult{Won: false}
}
