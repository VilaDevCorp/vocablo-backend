package activity

// import (
// 	"appname/ent"
// 	"appname/ent/activity"
// 	"appname/ent/predicate"
// 	"appname/ent/user"
// 	"appname/utils"
// 	"context"
// 	"math"

// 	"github.com/google/uuid"
// )

// type Svc interface {
// 	Create(ctx context.Context, form CreateForm) (*ent.Activity, error)
// 	Update(ctx context.Context, form UpdateForm) (*ent.Activity, error)
// 	Search(ctx context.Context, form SearchForm) (*utils.Page, error)
// 	Get(ctx context.Context, activityId uuid.UUID) (*ent.Activity, error)
// 	Delete(ctx context.Context, activityId []uuid.UUID) error
// }

// type Store struct {
// 	DB *ent.Client
// }

// func (s *Store) Create(ctx context.Context, form CreateForm) (*ent.Activity, error) {
// 	return s.DB.Activity.Create().SetName(form.Name).SetDescription(form.Description).SetIcon(form.Icon).SetSize(form.Size).SetUserID(form.UserId).Save(ctx)
// }

// func (s *Store) Update(ctx context.Context, form UpdateForm) (*ent.Activity, error) {
// 	update := s.DB.Activity.UpdateOneID(form.Id)
// 	if form.Name != nil {
// 		update.SetName(*form.Name)
// 	}
// 	if form.Description != nil {
// 		update.SetDescription(*form.Description)
// 	}
// 	if form.Icon != nil {
// 		update.SetIcon(*form.Icon)
// 	}
// 	if form.Size != nil {
// 		update.SetSize(*form.Size)
// 	}

// 	return update.Save(ctx)
// }

// func (s *Store) Get(ctx context.Context, activityId uuid.UUID) (*ent.Activity, error) {
// 	return s.DB.Activity.Get(ctx, activityId)
// }

// func (s *Store) Search(ctx context.Context, form SearchForm) (*utils.Page, error) {
// 	query := s.DB.Activity.Query()
// 	var conditions []predicate.Activity

// 	offset := 0
// 	limit := 1000

// 	if form.PageSize > 0 {
// 		offset = form.PageSize * form.Page
// 		limit = form.PageSize
// 	}

// 	if form.Name != nil {
// 		conditions = append(conditions, activity.NameContainsFold(*form.Name))
// 	}
// 	if form.Size != nil {
// 		conditions = append(conditions, activity.SizeEQ(*form.Size))
// 	}
// 	if form.UserId != nil {
// 		conditions = append(conditions, activity.HasUserWith(user.IDEQ(*form.UserId)))
// 	}
// 	totalRows, err := query.Where(activity.And(conditions...)).Count(ctx)
// 	content, err := query.Where(activity.And(conditions...)).Offset(offset).Limit(limit).All(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	totalPages := int(math.Ceil(float64(totalRows) / float64(form.PageSize)))
// 	page := utils.Page{Content: content, TotalRows: totalRows, TotalPages: totalPages}
// 	return &page, err
// }

// func (s *Store) Delete(ctx context.Context, activityIds []uuid.UUID) error {
// 	_, err := s.DB.Activity.Delete().Where(activity.IDIn(activityIds...)).Exec(ctx)
// 	return err
// }
