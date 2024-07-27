package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		CommonMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("Username").Unique().NotEmpty().StructTag(`json:"username"`),
		field.String("Email").Unique().NotEmpty().StructTag(`json:"email"`),
		field.String("Password").NotEmpty().StructTag(`json:"-"`),
		field.Bool("Validated").StorageKey("validated").Default(false).StructTag(`json:"validated"`),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("codes", VerificationCode.Type),
	}
}
