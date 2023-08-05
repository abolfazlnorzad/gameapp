package entity

type Player struct {
	ID      uint
	UserID  uint
	GameID  uint
	Score   uint
	Answers []uint
}

type PlayerAnswer struct {
	ID         uint
	PlayerID   uint
	QuestionID uint
	Choice     PossibleAnswerChoice
}
