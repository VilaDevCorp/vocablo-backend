package quiz

import "github.com/google/uuid"

type Quiz struct {
	Questions []QuizQuestion `json:"questions"`
	Score     int            `json:"score"`
}

type QuizQuestion struct {
	UserWordID       uuid.UUID `json:"userWordID"`
	Question         string    `json:"question"`
	Options          []string  `json:"options"`
	CorrectOptionPos int       `json:"correctOptionPos"`
	AnswerPos        int       `json:"answerPos"`
}
