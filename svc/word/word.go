package word

import (
	"context"
	"vocablo/customerrors"
	"vocablo/ent"
	"vocablo/ent/language"
)

type WordSvc interface {
	Create(ctx context.Context, form CreateForm) (*ent.Word, error)
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
