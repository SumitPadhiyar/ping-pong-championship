package models

type Player struct {
	ID              int
	Name            string
	PlayerRemoteURL string
}

type Game struct {
	ID             int
	FirstPlayer    *GamePlayerInfo
	SecondPlayer   *GamePlayerInfo
	WinnerPlayerID int
}

type GamePlayerInfo struct {
	Player       *Player
	Score        int
	DefenceArray []int
}
