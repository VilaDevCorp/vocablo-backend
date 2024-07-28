package test

import (
	"context"
	"encoding/json"
	"testing"
	"vocablo/ent"
	"vocablo/ent/language"
	"vocablo/ent/user"
	"vocablo/ent/userword"
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

	_, err = client.UserWord.Query().Where(userword.TermEQ(testWordForm2.Term)).Only(ctx)
	assert.NoError(t, err)
}
