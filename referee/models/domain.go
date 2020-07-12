package models

type Player struct {
	ID              int
	Name            string
	PlayerRemoteURL string
}

type Game struct {
	ID           int
	FirstPlayer  *GamePlayerInfo
	SecondPlayer *GamePlayerInfo
	WinnerPlayer *GamePlayerInfo
}

type GamePlayerInfo struct {
	Player     *Player
	Score      int
	DefenceMap map[int]struct{}
}
