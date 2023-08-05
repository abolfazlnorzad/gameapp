package entity

type Question struct {
	ID                uint
	Text              string
	Correct           PossibleAnswerChoice
	PossibleAnswerIDs []uint
	CategoryID        uint
	Difficulty        QuestionDifficulty
}

type QuestionDifficulty uint8

const (
	QuestionDifficultyEasy QuestionDifficulty = iota + 1
	QuestionDifficultyMedium
	QuestionDifficultyHard
)

func (q QuestionDifficulty) IsValid() bool {
	if q > QuestionDifficultyHard || q < QuestionDifficultyEasy {
		return false
	}

	return true
}

type PossibleAnswer struct {
	ID    uint
	Text  string
	Index PossibleAnswerChoice
}

type PossibleAnswerChoice uint8

const (
	PossibleAnswerA PossibleAnswerChoice = iota + 1
	PossibleAnswerB
	PossibleAnswerC
	PossibleAnswerD
)

func (p PossibleAnswerChoice) IsValid() bool {
	if p < PossibleAnswerA || p > PossibleAnswerD {
		return false
	}
	return true
}
