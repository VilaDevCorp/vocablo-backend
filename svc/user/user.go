package user

import (
	"appname/ent"
	"appname/ent/predicate"
	"appname/ent/user"
	"appname/utils"
	"context"
	"math"

	"github.com/google/uuid"
)

type UserSvc interface {
	Create(ctx context.Context, form CreateForm) (*ent.User, error)
	Search(ctx context.Context, form SearchForm) (*utils.Page, error)
	Get(ctx context.Context, userId uuid.UUID) (*ent.User, error)
	GetByUsername(ctx context.Context, username string) (*ent.User, error)
	Delete(ctx context.Context, userId uuid.UUID) error
}

type UserSvcImpl struct {
	DB *ent.Client
}

func (s *UserSvcImpl) Create(ctx context.Context, form CreateForm) (*ent.User, error) {
	return s.DB.User.Create().SetUsername(form.Username).SetPassword(form.Password).SetEmail(form.Email).Save(ctx)
}

func (s *UserSvcImpl) Get(ctx context.Context, userId uuid.UUID) (*ent.User, error) {
	return s.DB.User.Get(ctx, userId)
}

func (s *UserSvcImpl) GetByUsername(ctx context.Context, username string) (*ent.User, error) {
	return s.DB.User.Query().Where(user.UsernameEQ(username)).Only(ctx)
}

func (s *UserSvcImpl) Search(ctx context.Context, form SearchForm) (*utils.Page, error) {
	query := s.DB.User.Query()
	var conditions []predicate.User

	offset := 0
	limit := 1000

	if form.PageSize > 0 {
		offset = form.PageSize * form.Page
		limit = form.PageSize
	}
	if form.Name != nil {
		conditions = append(conditions, user.UsernameContains(*form.Name))
	}

	totalRows, err := query.Where(user.And(conditions...)).Count(ctx)
	if err != nil {
		return nil, err
	}
	content, err := query.Where(user.And(conditions...)).Offset(offset).Limit(limit).All(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(form.PageSize)))
	page := utils.Page{Content: content, TotalRows: totalRows, TotalPages: totalPages}
	return &page, err
}

func (s *UserSvcImpl) Delete(ctx context.Context, userId uuid.UUID) error {
	err := s.DB.User.DeleteOneID(userId).Exec(ctx)
	return err
}
