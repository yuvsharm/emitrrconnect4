package server

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"emitrrconnect4/internal/game"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var activePlayers = make(map[string]*Player)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	var init struct {
		Username string `json:"username"`
	}
	if err := conn.ReadJSON(&init); err != nil {
		log.Println("Init read error:", err)
		return
	}

	player, exists := activePlayers[init.Username]
	if !exists {
		player = &Player{
			Username: init.Username,
			Conn:     conn,
			Symbol:   "X",
			IsBot:    false,
		}
		activePlayers[init.Username] = player
	} else {
		player.Conn = conn
	}

	log.Println("🔗 Player joined:", player.Username)

	gameInstance := MatchPlayer(player)
	if gameInstance == nil {
		return
	}

	// ✨ EXTRA ADDED: Assign symbols correctly for multiplayer
	if p1, ok := gameInstance.Player1.(*Player); ok {
		p1.Symbol = "X"
	}
	if p2, ok := gameInstance.Player2.(*Player); ok {
		p2.Symbol = "O" // Player 2 must be "O"
	}

	sendState(gameInstance)

	for {
		var msg struct {
			Column *int `json:"column,omitempty"`
			Undo   bool `json:"undo,omitempty"`
		}

		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("Client disconnected:", player.Username)
			delete(activePlayers, player.Username)
			return
		}

		if msg.Undo {
			if gameInstance.UndoLastTwoMoves() {
				sendState(gameInstance)
			}
			continue
		}

		// Turn validation: Ensure player only moves on their turn symbol
		if gameInstance.LastActionWasUndo || msg.Column == nil ||
			gameInstance.Winner != "" || gameInstance.IsDraw ||
			gameInstance.Turn != player.Symbol {
			continue
		}

		if !gameInstance.Play(*msg.Column) {
			continue
		}

		gameInstance.LastActionWasUndo = false
		sendState(gameInstance)

		// 🤖 INTELLIGENT BOT MOVE
		if gameInstance.Player2.IsBotPlayer() &&
			gameInstance.Turn == "O" &&
			gameInstance.Winner == "" &&
			!gameInstance.IsDraw {

			time.Sleep(700 * time.Millisecond) // Realism delay

			botCol := findSmartBotMove(gameInstance)
			if botCol != -1 {
				gameInstance.Play(botCol)
				log.Println("🤖 Smart BOT played column:", botCol)
				sendState(gameInstance)
			}
		}
	}
}

// ... (findSmartBotMove and all helper functions stay exactly as they were) ...

// ✨ UPDATED: sendState updated to handle opponent names without removing existing data
func sendState(g *game.Game) {
	var duration float64 = 0
	if !g.StartTime.IsZero() && !g.EndTime.IsZero() {
		duration = g.EndTime.Sub(g.StartTime).Seconds()
	}

	// Base state data preserved
	state := map[string]interface{}{
		"board":     g.Board,
		"player":    g.Turn,
		"winner":    g.Winner,
		"isDraw":    g.IsDraw,
		"winCells":  g.WinCells,
		"duration":  duration,
		"isBotGame": g.Player2.IsBotPlayer(),
	}

	// ✨ Added opponentName logic to sync UI between devices
	if p, ok := g.Player1.(*Player); ok && p.Conn != nil {
		state["opponentName"] = g.Player2.GetUsername()
		p.Conn.WriteJSON(state)
	}
	if p, ok := g.Player2.(*Player); ok && p.Conn != nil {
		state["opponentName"] = g.Player1.GetUsername()
		p.Conn.WriteJSON(state)
	}
}