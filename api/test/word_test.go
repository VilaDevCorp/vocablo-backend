package test

import (
	"encoding/json"
	"testing"
	"vocablo/ent"
	"vocablo/utils"

	"github.com/stretchr/testify/assert"
)

func TestSearchWord(t *testing.T) {
	client, teardown, ctx := SetupTest(t, true, SetupUserWordTest)
	defer teardown(t)

	wordsInDb, err := client.Word.Query().All(ctx)
	if err != nil {
		t.Errorf("Error querying words: %s", err)
	}
	//We check that no words are in the db
	assert.Equal(t, 0, len(wordsInDb))

	resp := testEnv.MakeAuthRequest("GET", "/api/word/en/"+WORD_TO_SEARCH, nil, ctx)
	assert.Equal(t, 200, resp.Code)

	var respBody utils.ResponseBody
	var respData []ent.Word
	err = json.Unmarshal(resp.Body.Bytes(), &respBody)
	if err != nil {
		t.Fatal(err)
	}
	respDataStr, err := json.Marshal(respBody.Data)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(respDataStr, &respData)
	if err != nil {
		t.Fatal(err)
	}

	//We check that the response is not empty
	assert.NotEqual(t, 0, len(respData))
	wordsInDb, err = client.Word.Query().All(ctx)
	if err != nil {
		t.Errorf("Error querying words: %s", err)
	}
	//We check that the words are now in the db
	wordsInDbLength := len(wordsInDb)
	assert.NotEqual(t, 0, wordsInDbLength)

	//We check that the words are the same after another search
	resp = testEnv.MakeAuthRequest("GET", "/api/word/en/"+WORD_TO_SEARCH, nil, ctx)
	wordsInDb, err = client.Word.Query().All(ctx)
	if err != nil {
		t.Errorf("Error querying words: %s", err)
	}
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, wordsInDbLength, len(wordsInDb))

}
