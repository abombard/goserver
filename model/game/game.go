package game

import (
	"errors"
)

var ErrGameNotFound = errors.New("game not found")

type Game struct {
	StartDate   string `json:"startDate"`
	Type        string `json:"type"`
	PlayerCount int    `json:"playerCount"`
	RoundCount  int    `json:"roundCount"`
}

type GameDB struct {
	Game `gorm:"embedded"`
}
