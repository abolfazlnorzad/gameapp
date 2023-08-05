package entity

type Game struct {
	ID          uint
	QuestionIDs []uint
	PlayerIDs   []uint
	CategoryID  uint
}
