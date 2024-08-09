package userword

import (
	"context"
	"fmt"
	"vocablo/customerrors"
	"vocablo/ent"
	"vocablo/ent/language"
	"vocablo/ent/user"
	"vocablo/ent/userword"
	"vocablo/utils"

	"github.com/google/uuid"
)

type UserWordSvc interface {
	Create(ctx context.Context, form CreateForm) (*ent.UserWord, error)
	Update(ctx context.Context, form UpdateForm) (*ent.UserWord, error)
	Get(ctx context.Context, id string) (*ent.UserWord, error)
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
	userWord, err := s.DB.UserWord.Create().SetTerm(form.Term).SetDefinitions(form.Definitions).SetLangID(lang.ID).SetUserID(userID).Save(ctx)
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

func (s *UserWordSvcImpl) Get(ctx context.Context, id string) (*ent.UserWord, error) {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	userWord, err := s.DB.UserWord.Query().Where(userword.ID(uuidId)).WithUser().Only(ctx)
	if err != nil {
		if _, ok := err.(*ent.NotFoundError); ok {
			return nil, customerrors.NotFoundError{}
		}
		return nil, err
	}

	if userWord.Edges.User.ID != ctx.Value(utils.UserIdKey).(uuid.UUID) {
		return nil, customerrors.NotAllowedResourceError{}
	}
	return userWord, nil
}

func (s *UserWordSvcImpl) Search(ctx context.Context, form SearchForm) (*utils.Page[*ent.UserWord], error) {
	if form.Page <= 0 {
		form.Page = 0
	}
	if form.PageSize <= 0 {
		form.PageSize = 10
	}
	//we only allow the user to search his own words
	query := s.DB.UserWord.Query().Where(userword.HasUserWith(user.IDEQ(ctx.Value(utils.UserIdKey).(uuid.UUID))))
	if form.Term != nil && *form.Term != "" {
		query = query.Where(userword.TermContainsFold(*form.Term))
	}
	if form.Lang != nil {
		query = query.Where(userword.HasLangWith(language.CodeEqualFold(*form.Lang)))
	}
	if form.Learned != nil {
		if *form.Learned {
			query = query.Where(userword.LearningProgressGTE(100))
		} else {
			query = query.Where(userword.LearningProgressLT(100))
		}
	}
	total, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}
	page := utils.Page[*ent.UserWord]{PageNumber: form.Page}
	page.NElements = total

	if form.Count {
		return &page, nil
	}

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
