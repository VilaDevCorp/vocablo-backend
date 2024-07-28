package test

import (
	"context"
	"encoding/json"
	"testing"
	"vocablo/ent"
	"vocablo/ent/language"
	"vocablo/ent/user"
	entuserword "vocablo/ent/userword"
	"vocablo/svc/userword"
	"vocablo/utils"

	"github.com/stretchr/testify/assert"
)

func SetupUserWordTest(client *ent.Client, t *testing.T, ctx context.Context) {
	client.Language.Create().SetCode("en").SaveX(ctx)
	client.Language.Create().SetCode("es").SaveX(ctx)
	mainUser, err := client.User.Query().Where(user.UsernameEQ(testUserForm.Username)).Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	client.UserWord.Create().SetTerm(testWordForm1.Term).SetLang(
		client.Language.Query().Where(language.CodeEqualFold(testWordForm1.Lang)).OnlyX(ctx)).
		SetDefinitions(testWordForm1.Definitions).SetUserID(mainUser.ID).SaveX(ctx)
}

func TestCreateWord(t *testing.T) {
	client, teardown, ctx := SetupTest(t, true, SetupUserWordTest)
	defer teardown(t)

	body, err := json.Marshal(testWordForm2)
	if err != nil {
		t.Fatal(err)
	}
	resp := testEnv.MakeAuthRequest("POST", "/api/userword", utils.GetStringPointer(string(body)), ctx)

	assert.Equal(t, 200, resp.Code)

	_, err = client.UserWord.Query().Where(entuserword.TermEQ(testWordForm2.Term)).Only(ctx)
	assert.NoError(t, err)
}

func TestCreateWordEmptyFields(t *testing.T) {
	_, teardown, ctx := SetupTest(t, true, SetupUserWordTest)
	defer teardown(t)

	testWordFormEmpty := testWordForm2
	testWordFormEmpty.Term = ""
	body, err := json.Marshal(testWordFormEmpty)
	if err != nil {
		t.Fatal(err)
	}
	resp := testEnv.MakeAuthRequest("POST", "/api/userword", utils.GetStringPointer(string(body)), ctx)

	assert.Equal(t, 400, resp.Code)
}

func TestUpdateWord(t *testing.T) {
	client, teardown, ctx := SetupTest(t, true, SetupUserWordTest)
	defer teardown(t)

	userWordCreated, err := client.UserWord.Query().Where(entuserword.TermEQ(testWordForm1.Term)).Only(ctx)

	if err != nil {
		t.Fatal(err)
	}

	updateForm := userword.UpdateForm{ID: userWordCreated.ID.String(), Term: utils.GetStringPointer(UPDATED_TERM)}
	body, err := json.Marshal(updateForm)
	if err != nil {
		t.Fatal(err)
	}
	resp := testEnv.MakeAuthRequest("PUT", "/api/userword", utils.GetStringPointer(string(body)), ctx)

	assert.Equal(t, 200, resp.Code)

	updatedUserWord, err := client.UserWord.Query().Where(entuserword.IDEQ(userWordCreated.ID)).Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, UPDATED_TERM, updatedUserWord.Term)
	assert.Equal(t, testWordForm1.Definitions[0].Definition, updatedUserWord.Definitions[0].Definition)
}

func TestUpdateWordEmptyId(t *testing.T) {
	_, teardown, ctx := SetupTest(t, true, SetupUserWordTest)
	defer teardown(t)

	updateForm := userword.UpdateForm{ID: "", Term: utils.GetStringPointer(UPDATED_TERM)}
	body, err := json.Marshal(updateForm)
	if err != nil {
		t.Fatal(err)
	}
	resp := testEnv.MakeAuthRequest("PUT", "/api/userword", utils.GetStringPointer(string(body)), ctx)

	assert.Equal(t, 400, resp.Code)
}

func TestSearchWord(t *testing.T) {
	client, teardown, ctx := SetupTest(t, true, SetupUserWordTest)
	defer teardown(t)

	_, err := client.UserWord.Create().SetTerm(testWordForm2.Term).SetLang(
		client.Language.Query().Where(language.CodeEqualFold(testWordForm2.Lang)).OnlyX(ctx)).
		SetDefinitions(testWordForm2.Definitions).Save(ctx)

	if err != nil {
		t.Fatal(err)
	}
	searchForm := userword.SearchForm{Term: utils.GetStringPointer(testWordForm1.Term)}
	body, err := json.Marshal(searchForm)
	if err != nil {
		t.Fatal(err)
	}
	resp := testEnv.MakeAuthRequest("POST", "/api/userword/search", utils.GetStringPointer(string(body)), ctx)
	assert.Equal(t, 200, resp.Code)
	var respBodyObj utils.ResponseBody
	var respData utils.Page[ent.UserWord]
	err = json.Unmarshal(resp.Body.Bytes(), &respBodyObj)
	if err != nil {
		t.Fatal(err)
	}
	respDataStr, err := json.Marshal(respBodyObj.Data)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(respDataStr, &respData)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(respData.Content))
	assert.Equal(t, testWordForm1.Term, respData.Content[0].Term)
	assert.Equal(t, testWordForm1.Definitions[0].Definition, respData.Content[0].Definitions[0].Definition)

	searchForm2 := userword.SearchForm{Term: utils.GetStringPointer(testWordForm2.Term)}
	body, err = json.Marshal(searchForm2)
	if err != nil {
		t.Fatal(err)
	}
	resp = testEnv.MakeAuthRequest("POST", "/api/userword/search", utils.GetStringPointer(string(body)), ctx)
	assert.Equal(t, 200, resp.Code)
	err = json.Unmarshal(resp.Body.Bytes(), &respBodyObj)
	if err != nil {
		t.Fatal(err)
	}
	respDataStr, err = json.Marshal(respBodyObj.Data)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(respDataStr, &respData)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(respData.Content))
	assert.Equal(t, testWordForm2.Term, respData.Content[0].Term)
	assert.Equal(t, testWordForm2.Definitions[0].Definition, respData.Content[0].Definitions[0].Definition)

	searchForm3 := userword.SearchForm{Lang: utils.GetStringPointer("en"), Page: 0, PageSize: 1}
	body, err = json.Marshal(searchForm3)
	if err != nil {
		t.Fatal(err)
	}
	resp = testEnv.MakeAuthRequest("POST", "/api/userword/search", utils.GetStringPointer(string(body)), ctx)
	assert.Equal(t, 200, resp.Code)
	err = json.Unmarshal(resp.Body.Bytes(), &respBodyObj)
	if err != nil {
		t.Fatal(err)
	}
	respDataStr, err = json.Marshal(respBodyObj.Data)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(respDataStr, &respData)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(respData.Content))
	assert.Equal(t, respData.PageNumber, 0)
	assert.True(t, respData.HasNext)

	searchForm4 := userword.SearchForm{Lang: utils.GetStringPointer("es"), Page: 0, PageSize: 1}
	body, err = json.Marshal(searchForm4)
	if err != nil {
		t.Fatal(err)
	}
	resp = testEnv.MakeAuthRequest("POST", "/api/userword/search", utils.GetStringPointer(string(body)), ctx)
	assert.Equal(t, 200, resp.Code)
	err = json.Unmarshal(resp.Body.Bytes(), &respBodyObj)
	if err != nil {
		t.Fatal(err)
	}
	respDataStr, err = json.Marshal(respBodyObj.Data)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(respDataStr, &respData)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, len(respData.Content))
}

func TestDeleteWord(t *testing.T) {
	client, teardown, ctx := SetupTest(t, true, SetupUserWordTest)
	defer teardown(t)

	userWordCreated, err := client.UserWord.Query().Where(entuserword.TermEQ(testWordForm1.Term)).Only(ctx)

	if err != nil {
		t.Fatal(err)
	}

	resp := testEnv.MakeAuthRequest("DELETE", "/api/userword/"+userWordCreated.ID.String(), nil, ctx)

	assert.Equal(t, 200, resp.Code)

	_, err = client.UserWord.Query().Where(entuserword.IDEQ(userWordCreated.ID)).Only(ctx)
	assert.Error(t, err)
}
