package test

import (
	"context"
	"encoding/json"
	"testing"
	"vocablo/customerrors"
	"vocablo/ent"
	"vocablo/ent/language"
	"vocablo/ent/user"
	"vocablo/svc/quiz"
	"vocablo/utils"

	"github.com/stretchr/testify/assert"
)

func SetupQuizTest(client *ent.Client, t *testing.T, ctx context.Context) {
	mainUser, err := client.User.Query().Where(user.UsernameEQ(testUserForm1.Username)).Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	client.UserWord.Create().SetTerm(testWordForm1.Term).SetLang(
		client.Language.Query().Where(language.CodeEqualFold(testWordForm1.Lang)).OnlyX(ctx)).
		SetDefinitions(testWordForm1.Definitions).SetUserID(mainUser.ID).SaveX(ctx)
	client.UserWord.Create().SetTerm(testWordForm2.Term).SetLang(
		client.Language.Query().Where(language.CodeEqualFold(testWordForm2.Lang)).OnlyX(ctx)).
		SetDefinitions(testWordForm2.Definitions).SetUserID(mainUser.ID).SaveX(ctx)
	client.UserWord.Create().SetTerm(otherUserWordForm.Term).SetLang(
		client.Language.Query().Where(language.CodeEqualFold(otherUserWordForm.Lang)).OnlyX(ctx)).
		SetDefinitions(otherUserWordForm.Definitions).SetUserID(mainUser.ID).SaveX(ctx)
}

func TestCreateQuiz(t *testing.T) {
	client, teardown, ctx := SetupTest(t, true, SetupQuizTest)
	defer teardown(t)
	mainUser, err := client.User.Query().Where(user.UsernameEQ(testUserForm1.Username)).Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	client.UserWord.Create().SetTerm(testWordForm3.Term).SetLang(
		client.Language.Query().Where(language.CodeEqualFold(testWordForm3.Lang)).OnlyX(ctx)).
		SetDefinitions(testWordForm3.Definitions).SetUserID(mainUser.ID).SaveX(ctx)

	body, err := json.Marshal(testCreateQuizForm)
	if err != nil {
		t.Fatal(err)
	}
	resp := testEnv.MakeAuthRequest("POST", "/api/quiz", utils.GetStringPointer(string(body)), ctx)
	var respObj utils.ResponseBody
	err = json.Unmarshal(resp.Body.Bytes(), &respObj)
	if err != nil {
		t.Fatal(err)
	}
	var respQuiz quiz.Quiz
	respQuizStr, err := json.Marshal(respObj.Data)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(respQuizStr, &respQuiz)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 200, resp.Code)
	for _, question := range respQuiz.Questions {
		questionTerm := question.Question
		switch questionTerm {
		case testWordForm1.Term:
			assert.Equal(t, testWordForm1.Definitions[0].Definition,
				question.Options[question.CorrectOptionPos])
		case testWordForm2.Term:
			assert.Equal(t, testWordForm2.Definitions[0].Definition,
				question.Options[question.CorrectOptionPos])
		case testWordForm3.Term:
			assert.Equal(t, testWordForm3.Definitions[0].Definition,
				question.Options[question.CorrectOptionPos])
		case otherUserWordForm.Term:
			assert.Equal(t, otherUserWordForm.Definitions[0].Definition,
				question.Options[question.CorrectOptionPos])
		}
		//We check that the options are different and not empty
		for _, option := range question.Options {
			assert.NotEqual(t, "", option)
			matches := 0
			for _, otherOption := range question.Options {
				if option == otherOption {
					matches++
				}
			}
			assert.Equal(t, 1, matches)
		}
	}
}

func TestCreateQuizNotEnoughWords(t *testing.T) {
	_, teardown, ctx := SetupTest(t, true, SetupQuizTest)
	defer teardown(t)

	body, err := json.Marshal(testCreateQuizForm)
	if err != nil {
		t.Fatal(err)
	}
	resp := testEnv.MakeAuthRequest("POST", "/api/quiz", utils.GetStringPointer(string(body)), ctx)
	var respObj utils.ResponseBody
	err = json.Unmarshal(resp.Body.Bytes(), &respObj)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, customerrors.NOT_ENOUGH_WORDS_FOR_QUIZ, *respObj.ErrorCode)
	assert.Equal(t, 409, resp.Code)
}

func TestAnswerQuiz(t *testing.T) {
	client, teardown, ctx := SetupTest(t, true, SetupQuizTest)
	defer teardown(t)
	mainUser, err := client.User.Query().Where(user.UsernameEQ(testUserForm1.Username)).Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	client.UserWord.Create().SetTerm(testWordForm3.Term).SetLang(
		client.Language.Query().Where(language.CodeEqualFold(testWordForm3.Lang)).OnlyX(ctx)).
		SetDefinitions(testWordForm3.Definitions).SetUserID(mainUser.ID).SaveX(ctx)

	body, err := json.Marshal(testCreateQuizForm)
	if err != nil {
		t.Fatal(err)
	}
	resp := testEnv.MakeAuthRequest("POST", "/api/quiz", utils.GetStringPointer(string(body)), ctx)
	var respObj utils.ResponseBody
	err = json.Unmarshal(resp.Body.Bytes(), &respObj)
	if err != nil {
		t.Fatal(err)
	}
	var respQuiz quiz.Quiz
	respQuizStr, err := json.Marshal(respObj.Data)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(respQuizStr, &respQuiz)
	if err != nil {
		t.Fatal(err)
	}

	//We answer the quiz with the correct answers
	for i, question := range respQuiz.Questions {
		respQuiz.Questions[i].AnswerPos = question.CorrectOptionPos
	}
	body, err = json.Marshal(respQuiz)
	if err != nil {
		t.Fatal(err)
	}
	resp = testEnv.MakeAuthRequest("POST", "/api/quiz/answer", utils.GetStringPointer(string(body)), ctx)
	err = json.Unmarshal(resp.Body.Bytes(), &respObj)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, 100, int(respObj.Data.(float64)))

	//We answer the quiz with the only half of the correct answers
	for i, question := range respQuiz.Questions {
		if i%2 == 0 {
			respQuiz.Questions[i].AnswerPos = question.CorrectOptionPos
		} else {
			respQuiz.Questions[i].AnswerPos = (question.CorrectOptionPos + 1) % 4
		}
	}
	body, err = json.Marshal(respQuiz)
	if err != nil {
		t.Fatal(err)
	}
	resp = testEnv.MakeAuthRequest("POST", "/api/quiz/answer", utils.GetStringPointer(string(body)), ctx)
	err = json.Unmarshal(resp.Body.Bytes(), &respObj)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, 50, int(respObj.Data.(float64)))

}
