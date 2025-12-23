package server

import (
	"math/rand"
	"time"

	"emitrrconnect4/internal/game" // Is path ko apne project ke hisaab se check kar lein
)

// NewBot creates a new bot player
func NewBot() *Player {
	return &Player{
		Username: "BOT",
		Symbol:   "O",
		IsBot:    true,
	}
}

// BotMove calculates the best move for the bot
func BotMove(board [][]string) int {
	rand.Seed(time.Now().UnixNano())

	// 1. Get Valid Columns (Checking which columns aren't full)
	var validCols []int
	for c := 0; c < 7; c++ {
		if board[0][c] == "" {
			validCols = append(validCols, c)
		}
	}

	if len(validCols) == 0 {
		return -1
	}

	maxScore := -200000
	var bestMoves []int

	// Center-out priority order (3 is middle)
	searchOrder := []int{3, 2, 4, 1, 5, 0, 6}

	for _, col := range searchOrder {
		// Column validity check
		if board[0][col] != "" {
			continue
		}

		score := 0

		// A. Can Bot win now?
		if canBotWinSim(board, col, "O") {
			return col
		}

		// B. Block Player?
		if canBotWinSim(board, col, "X") {
			score += 50000
		}

		// C. Strategic Scoring (Middle columns are better)
		score += centerWeightLocal(col)
		
		// D. Look-Ahead: Don't let the opponent win on the next turn
		tempBoard := copyLocalBoard(board)
		row := findLandingRow(tempBoard, col)
		if row != -1 {
			tempBoard[row][col] = "O"
			if opponentCanWinNextLocal(tempBoard, "X") {
				score -= 40000
			}
		}

		if score > maxScore {
			maxScore = score
			bestMoves = []int{col}
		} else if score == maxScore {
			bestMoves = append(bestMoves, col)
		}
	}

	if len(bestMoves) > 0 {
		return bestMoves[rand.Intn(len(bestMoves))]
	}
	return -1
}

/* --- Internal Helpers to Stop Errors --- */

func copyLocalBoard(board [][]string) [][]string {
	duplicate := make([][]string, len(board))
	for i := range board {
		duplicate[i] = make([]string, len(board[i]))
		copy(duplicate[i], board[i])
	}
	return duplicate
}

func findLandingRow(board [][]string, col int) int {
	for r := len(board) - 1; r >= 0; r-- {
		if board[r][col] == "" {
			return r
		}
	}
	return -1
}

func canBotWinSim(board [][]string, col int, symbol string) bool {
	temp := copyLocalBoard(board)
	row := findLandingRow(temp, col)
	if row == -1 {
		return false
	}
	temp[row][col] = symbol
	
	// Calling the CheckWin from game package
	res := game.CheckWin(temp, symbol)
	return res.Won
}

func opponentCanWinNextLocal(board [][]string, opponentSymbol string) bool {
	for c := 0; c < 7; c++ {
		if board[0][c] == "" && canBotWinSim(board, c, opponentSymbol) {
			return true
		}
	}
	return false
}

func centerWeightLocal(col int) int {
	weights := map[int]int{3: 100, 2: 60, 4: 60, 1: 20, 5: 20, 0: 0, 6: 0}
	return weights[col]
}