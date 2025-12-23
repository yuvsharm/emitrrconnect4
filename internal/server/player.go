package server

import "github.com/gorilla/websocket"

type Player struct {
	Username string
	Conn     *websocket.Conn
	Symbol   string
	IsBot    bool
}

func (p *Player) GetUsername() string {
	return p.Username
}

func (p *Player) GetSymbol() string {
	return p.Symbol
}

func (p *Player) IsBotPlayer() bool {
	return p.IsBot
}

