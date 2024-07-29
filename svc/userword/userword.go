package userword

import (
	"context"
	"fmt"
	"vocablo/customerrors"
	"vocablo/ent"
	"vocablo/ent/language"
	"vocablo/ent/userword"
	"vocablo/utils"

	"github.com/google/uuid"
)

type UserWordSvc interface {
	Create(ctx context.Context, form CreateForm) (*ent.UserWord, error)
	Update(ctx context.Context, form UpdateForm) (*ent.UserWord, error)
	Search(ctx context.Context, form SearchForm) (*utils.Page[*ent.UserWord], error)
	Delete(ctx context.Context, id string) error
}

type UserWordSvcImpl struct {
	DB *ent.Client
}

func (s *UserWordSvcImpl) Create(ctx context.Context, form CreateForm) (*ent.UserWord, error) {
	if form.Term == "" || form.Lang == "" || form.Definitions == nil || len(form.Definitions) == 0 {
		return nil, customerrors.EmptyFormFieldsError{}
	}
	lang, err := s.DB.Language.Query().Where(language.CodeEQ(form.Lang)).Only(ctx)
	if err != nil {
		return nil, err
	}
	userID := ctx.Value(utils.UserIdKey).(uuid.UUID)
	userWord, err := s.DB.UserWord.Create().SetTerm(form.Term).SetLangID(lang.ID).SetUserID(userID).Save(ctx)
	if err != nil {
		return nil, err
	}
	return userWord, nil
}

func (s *UserWordSvcImpl) Update(ctx context.Context, form UpdateForm) (*ent.UserWord, error) {
	if form.ID == "" {
		return nil, customerrors.EmptyFormFieldsError{}
	}
	uuidId, err := uuid.Parse(form.ID)

	if err != nil {
		return nil, err
	}
	userWord, err := s.DB.UserWord.Query().Where(userword.ID(uuidId)).WithUser().Only(ctx)
	if err != nil {
		return nil, err
	}
	userIdKey := ctx.Value(utils.UserIdKey)
	fmt.Println(userIdKey)
	if userWord.Edges.User.ID != ctx.Value(utils.UserIdKey).(uuid.UUID) {
		return nil, customerrors.NotAllowedResourceError{}
	}

	updateBuilder := s.DB.UserWord.UpdateOneID(uuidId)

	if form.Term != nil {
		if (*form.Term) == "" {
			return nil, customerrors.EmptyFormFieldsError{}
		}
		updateBuilder.SetTerm(*form.Term)
	}
	if form.Definitions != nil {
		if len(*form.Definitions) == 0 {
			return nil, customerrors.EmptyFormFieldsError{}
		}
		updateBuilder.SetDefinitions(*form.Definitions)
	}

	updatedUserWord, err := updateBuilder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return updatedUserWord, nil
}

func (s *UserWordSvcImpl) Search(ctx context.Context, form SearchForm) (*utils.Page[*ent.UserWord], error) {
	if form.Page <= 0 {
		form.Page = 0
	}
	if form.PageSize <= 0 {
		form.PageSize = 10
	}
	query := s.DB.UserWord.Query()
	if form.Term != nil && *form.Term != "" {
		query = query.Where(userword.TermContainsFold(*form.Term))
	}
	if form.Lang != nil {
		query = query.Where(userword.HasLangWith(language.CodeEqualFold(*form.Lang)))
	}
	total, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}
	page := utils.Page[*ent.UserWord]{PageNumber: form.Page}
	if total > (form.Page+1)*form.PageSize {
		page.HasNext = true
	}
	userWords, err := query.Offset(form.Page * form.PageSize).Limit(form.PageSize).All(ctx)
	if err != nil {
		return nil, err
	}
	page.Content = userWords
	return &page, nil
}

func (s *UserWordSvcImpl) Delete(ctx context.Context, id string) error {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	userWord, err := s.DB.UserWord.Query().Where(userword.ID(uuidId)).WithUser().Only(ctx)
	if err != nil {
		return err
	}
	if userWord.Edges.User.ID != ctx.Value(utils.UserIdKey).(uuid.UUID) {
		return customerrors.NotAllowedResourceError{}
	}
	err = s.DB.UserWord.DeleteOneID(uuidId).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
