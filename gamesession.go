package main

import (
	"sync"
	"crypto/rand"
	"encoding/hex"
)

type GameSession struct {
	PredictionTree *Tree 
	GameBoard *BoardInstance 
	SessionID string 	
}

func newGameSession() *GameSession {

	gb   := NewBoardInstance()
	tree := NewTree(*gb)
	sID := generateSessionID()

	return &GameSession{
		PredictionTree: tree,
		GameBoard: gb,
		SessionID: sID,
	}

}

func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes) // genereate 16 random bytes
	return hex.EncodeToString(bytes) //show the bytes in hex string format 
}



// ------ Game Sessions START------------------------------
/*
	- 

*/
type GameSessions struct {
	GameSessionsRunning map[string] *GameSession 
	mu sync.Mutex
}

// constuctor for the GameSessions map 
func NewGameSessions() *GameSessions {
	return &GameSessions{
		GameSessionsRunning: make(map[string]*GameSession),
	}
}

// each new player/browser instance gets a new session 
func (gs *GameSessions) CreateNewSession() string {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	gamesesh := newGameSession()
	gs.GameSessionsRunning[gamesesh.SessionID] = gamesesh
	return gamesesh.SessionID
}

// returns a session pointer to the desired session. It's used to support the end point 
// which returns the progress /staut of the move tree building function that each game session has 
func (gs *GameSessions) GetSession(seshId string) (*GameSession, bool) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	sesh, exists := gs.GameSessionsRunning[seshId]
	if !exists {
		return nil, false
	}
	return sesh, true

}
