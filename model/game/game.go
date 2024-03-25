package game

import (
	"errors"
	"time"
)

var ErrGameNotFound = errors.New("game not found")

type Game struct {
	Id          string    `json:"id"`
	StartDate   time.Time `json:"startDate"`
	Type        string    `json:"type"`
	PlayerCount int       `json:"playerCount"`
	RoundCount  int       `json:"roundCount"`
}

type GameDB struct {
	Game `gorm:"embedded"`
}
