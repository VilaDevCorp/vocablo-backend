package quiz

import (
	"context"
	"math"
	"math/rand"
	"slices"
	"vocablo/customerrors"
	"vocablo/ent"
	"vocablo/ent/user"
	"vocablo/ent/userword"
	"vocablo/utils"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

type QuizSvc interface {
	Create(ctx context.Context, form CreateForm) (*Quiz, error)
	Answer(ctx context.Context, filledQuiz Quiz) (int, error)
}

type QuizSvcImpl struct {
	DB *ent.Client
}

func (s *QuizSvcImpl) Create(ctx context.Context, form CreateForm) (*Quiz, error) {
	nQuestions := 10
	if form.NQuestions > 0 {
		nQuestions = form.NQuestions
	}

	//We collect double the number of questions to have more options to put in the quiz (avoiding the already learned)
	userWords, err := s.DB.UserWord.Query().Limit(nQuestions * 2).
		Where(userword.And(userword.LearningProgressLT(100),
			userword.HasUserWith(user.IDEQ(ctx.Value(utils.UserIdKey).(uuid.UUID))))).Order(sql.OrderByRand()).All(ctx)

	//If the user has less than 4 words, we return an error (we need at least 4 words to create a quiz, 1 correct and 3 incorrect)
	if len(userWords) < 4 {
		return nil, customerrors.NotEnoughWordsForQuizError{}
	}

	//If we dont have enough words to fill the demand, we just return the ones we have
	nQuestions = int(math.Min(float64(nQuestions), float64(len(userWords))))

	questions := make([]QuizQuestion, 0, nQuestions)
	for i := 0; i < nQuestions; i++ {
		userWord := userWords[i]
		options := make([]string, 4)
		//We add one of the definitions of the word as the correct option (randomly to explore all the options)
		nDefinitions := len(userWord.Definitions)
		correctWordDefinitionPosition := rand.Intn(nDefinitions)
		correctOption := userWord.Definitions[correctWordDefinitionPosition]
		//We calculate the position of the correct option
		correctOptionPosition := rand.Intn(4)
		options[correctOptionPosition] = correctOption.Definition

		var alreadyUsedPositions []int
		alreadyUsedPositions = append(alreadyUsedPositions, i)

		//We search the other 3 options randomly in the other words definitions
		for j := 0; j < 3; j++ {
			//we select a random word, different from the other options already selected
			randomWordPosition := rand.Intn(len(userWords))
			for slices.Contains(alreadyUsedPositions, randomWordPosition) {
				randomWordPosition = rand.Intn(len(userWords))
			}
			randomWord := userWords[randomWordPosition]

			//We add the position to the already used positions
			alreadyUsedPositions = append(alreadyUsedPositions, randomWordPosition)

			//We select a random definition from the word
			randomDefinitionPosition := rand.Intn(len(randomWord.Definitions))
			for i, option := range options {
				if option == "" {
					options[i] = randomWord.Definitions[randomDefinitionPosition].Definition
					break
				}
			}
		}
		question := QuizQuestion{
			UserWordID:       userWord.ID,
			Question:         userWord.Term,
			Options:          options,
			CorrectOptionPos: correctOptionPosition,
		}
		questions = append(questions, question)
	}
	if err != nil {
		return nil, err
	}

	return &Quiz{Questions: questions}, nil
}

func (s *QuizSvcImpl) Answer(ctx context.Context, filledQuiz Quiz) (int, error) {
	totalScore := 0.0
	questionValue := 100.0 / float64(len(filledQuiz.Questions))
	clientTx, err := s.DB.Tx(ctx)
	if err != nil {
		return 0, err
	}
	for _, question := range filledQuiz.Questions {
		//If the answer is correct, we add the value of the question to the total score and we add 10 to the learning progress of the word
		if question.AnswerPos != nil && question.Options[*question.AnswerPos] == question.Options[question.CorrectOptionPos] {
			totalScore += questionValue
			err = clientTx.UserWord.UpdateOneID(question.UserWordID).AddLearningProgress(10).Exec(ctx)
			if err != nil {
				clientTx.Rollback()
				return 0, err
			}
		}
	}
	err = clientTx.Commit()
	if err != nil {
		return 0, err
	}
	return int(totalScore), nil
}
