package word

import (
	"context"
	"encoding/json"
	"net/http"
	"vocablo/apischema"
	"vocablo/customerrors"
	"vocablo/ent"
	"vocablo/ent/language"
	"vocablo/ent/word"
	"vocablo/utils"
)

type WordSvc interface {
	Create(ctx context.Context, form CreateForm) (*ent.Word, error)
	CreateBulk(ctx context.Context, forms []CreateForm) []*ent.Word
	Search(ctx context.Context, lang string, term string) (result *utils.Page[*ent.Word], err error)
}

type WordSvcImpl struct {
	DB *ent.Client
}

func (s *WordSvcImpl) Create(ctx context.Context, form CreateForm) (*ent.Word, error) {
	if form.Term == "" || form.Lang == "" || form.Definitions == nil || len(form.Definitions) == 0 {
		return nil, customerrors.EmptyFormFieldsError{}
	}
	lang, err := s.DB.Language.Query().Where(language.CodeEQ(form.Lang)).Only(ctx)
	if err != nil {
		return nil, err
	}
	word, err := s.DB.Word.Create().SetTerm(form.Term).SetLangID(lang.ID).Save(ctx)
	if err != nil {
		return nil, err
	}
	return word, nil
}

func (s *WordSvcImpl) CreateBulk(ctx context.Context, forms []CreateForm) []*ent.Word {

	if len(forms) == 0 {
		return nil
	}

	lang, err := s.DB.Language.Query().Where(language.CodeEQ(forms[0].Lang)).Only(ctx)
	if err != nil {
		return nil
	}

	var builders []*ent.WordCreate
	for _, form := range forms {
		if form.Term == "" || form.Lang == "" || form.Definitions == nil || len(form.Definitions) == 0 {
			continue
		}
		builders = append(builders, s.DB.Word.Create().SetTerm(form.Term).SetLangID(lang.ID).SetDefinitions(form.Definitions))

	}
	words, err := s.DB.Word.CreateBulk(builders...).Save(ctx)
	if err != nil {
		return nil
	}
	return words
}

func (s *WordSvcImpl) Search(ctx context.Context, lang string, term string) (result *utils.Page[*ent.Word], err error) {
	if term == "" || lang == "" {
		return nil, customerrors.EmptyFormFieldsError{}
	}
	query := s.DB.Word.Query()
	query = query.Where(word.TermContainsFold(term))
	query = query.Where(word.HasLangWith(language.CodeEqualFold(lang)))
	words, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	if len(words) == 0 {
		words, err = s.searchInApi(ctx, term, lang)
		if err != nil {
			return nil, err
		}
	}
	var page utils.Page[*ent.Word] = utils.Page[*ent.Word]{PageNumber: 0, Content: words, HasNext: false}
	return &page, nil
}

func (s *WordSvcImpl) searchInApi(ctx context.Context, term string, lang string) ([]*ent.Word, error) {
	if term == "" || lang == "" {
		return nil, customerrors.EmptyFormFieldsError{}
	}
	resp, err := http.Get("https://api.dictionaryapi.dev/api/v2/entries/en/" + term)
	if err != nil {
		return nil, err
	}
	var respObj apischema.ApiResponse
	err = json.NewDecoder(resp.Body).Decode(&respObj)
	if err != nil {
		return nil, err
	}
	forms := ConvertApiResponseToWordForms(respObj)
	createdWords := s.CreateBulk(ctx, forms)
	return createdWords, nil
}
