package entity

type Game struct {
	ID          uint
	QuestionIDs []uint
	PlayerIDs   []uint
	Category    Category
}
