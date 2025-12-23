package game

import (
	"log"
	"time"
)

type PlayerLike interface {
	GetUsername() string
	GetSymbol() string
	IsBotPlayer() bool
}

/* 🔥 Move struct */
type Move struct {
	Row    int
	Col    int
	Symbol string
}

type Game struct {
	Board             [][]string
	Turn              string
	Player1           PlayerLike
	Player2           PlayerLike
	Winner            string
	IsDraw            bool
	WinCells          [][2]int
	StartTime         time.Time
	EndTime           time.Time
	Moves             []Move
	LastActionWasUndo bool
}

func NewGame(p1, p2 PlayerLike) *Game {
	game := &Game{
		Board:     NewBoard(),
		Turn:      "X",
		Player1:   p1,
		Player2:   p2,
		StartTime: time.Now(),
		Moves:     []Move{},
	}

	log.Println("🎮 Game started")
	return game
}

/* ============================= */
/* ▶️ PLAYER MOVE               */
/* ============================= */

func (g *Game) Play(col int) bool {
	// Any real move cancels undo state
	g.LastActionWasUndo = false

	if g.Winner != "" || g.IsDraw {
		return false
	}

	if col < 0 || col >= Cols {
		log.Println("❌ Invalid column:", col)
		return false
	}

	current := g.getCurrentPlayer()

	// 🔽 DROP DISC
	row := -1
	for r := Rows - 1; r >= 0; r-- {
		if g.Board[r][col] == "" {
			g.Board[r][col] = current.GetSymbol()
			row = r
			break
		}
	}

	if row == -1 {
		log.Println("❌ Column full:", col)
		return false
	}

	// 🔥 SAVE MOVE
	g.Moves = append(g.Moves, Move{
		Row:    row,
		Col:    col,
		Symbol: current.GetSymbol(),
	})

	log.Println("✅ Disc placed:", current.GetSymbol(), "at", row, col)

	// 🏆 CHECK WIN
	if g.checkEnd(current) {
		return true
	}

	// 🔁 SWITCH TURN
	g.switchTurn()

	// 💡 NOTE: playBotIfNeeded() was removed from here.
	// The Bot logic is now controlled exclusively by ws.go 
	// to ensure the "Smart" version runs correctly.

	return true
}

/* ============================= */
/* 🧠 HELPERS                   */
/* ============================= */

func (g *Game) getCurrentPlayer() PlayerLike {
	if g.Turn == "X" {
		return g.Player1
	}
	return g.Player2
}

func (g *Game) switchTurn() {
	if g.Turn == "X" {
		g.Turn = "O"
	} else {
		g.Turn = "X"
	}
}

func (g *Game) checkEnd(player PlayerLike) bool {
	result := CheckWin(g.Board, player.GetSymbol())
	if result.Won {
		g.Winner = player.GetUsername()
		g.WinCells = result.WinCells
		g.EndTime = time.Now()
		return true
	}

	if IsBoardFull(g.Board) {
		g.IsDraw = true
		g.EndTime = time.Now()
		return true
	}

	return false
}

/* ============================= */
/* ⬅️ UNDO (PLAYER + BOT)       */
/* ============================= */

func (g *Game) UndoLastTwoMoves() bool {
	if len(g.Moves) < 2 {
		return false
	}

	// Remove two moves (Bot's move + Player's move)
	for i := 0; i < 2; i++ {
		last := g.Moves[len(g.Moves)-1]
		g.Moves = g.Moves[:len(g.Moves)-1]
		g.Board[last.Row][last.Col] = ""
	}

	g.Winner = ""
	g.IsDraw = false
	g.WinCells = nil
	g.EndTime = time.Time{}
	g.LastActionWasUndo = true

	// Restore turn back to "X" (assuming player always starts)
	g.Turn = "X"

	return true
}