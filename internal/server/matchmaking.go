package server

import (
	"sync"
	"time"

	"emitrrconnect4/internal/game"
)

var (
	waitingPlayer *Player
	mutex         sync.Mutex
	// Har player ke liye ek channel taaki unhe game mil sake
	gameChannels = make(map[string]chan *game.Game)
)

func MatchPlayer(p *Player) *game.Game {
	mutex.Lock()

	// 1. Check karein ki kya koi pehle se intezar kar raha hai
	if waitingPlayer != nil && waitingPlayer.Username != p.Username {
		opponent := waitingPlayer
		waitingPlayer = nil // Queue khali karo
		
		// Naya game banao dono players ke liye
		newGame := game.NewGame(opponent, p)
		
		// Pehle wale player ko signal bhejo ki game mil gaya hai
		if ch, ok := gameChannels[opponent.Username]; ok {
			ch <- newGame
		}
		
		mutex.Unlock()
		return newGame
	}

	// 2. Agar koi waiting nahi hai, toh waiting list mein jao
	waitingPlayer = p
	myChan := make(chan *game.Game, 1)
	gameChannels[p.Username] = myChan
	mutex.Unlock()

	// 3. Background mein 10 sec timer chalao (Sirf bot ke liye)
	go func(player *Player) {
		time.Sleep(10 * time.Second)
		mutex.Lock()
		// Agar 10 sec baad bhi yehi player waiting hai
		if waitingPlayer == player {
			waitingPlayer = nil
			bot := NewBot()
			newGame := game.NewGame(player, bot)
			if ch, ok := gameChannels[player.Username]; ok {
				ch <- newGame
			}
		}
		mutex.Unlock()
	}(p)

	// 4. Game ka intezar karein (Yeh blocking hai taaki game return ho sake)
	select {
	case g := <-myChan:
		mutex.Lock()
		delete(gameChannels, p.Username)
		mutex.Unlock()
		return g
	case <-time.After(12 * time.Second): // Safety Timeout
		return nil
	}
}