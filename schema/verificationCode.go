package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// VerificationCode holds the schema definition for the VerificationCode entity.
type VerificationCode struct {
	ent.Schema
}

func (VerificationCode) Mixin() []ent.Mixin {
	return []ent.Mixin{
		CommonMixin{},
	}
}

func (VerificationCode) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "verification_codes"},
	}
}

// Fields of the VerificationCode.
func (VerificationCode) Fields() []ent.Field {
	return []ent.Field{
		field.String("type").NotEmpty(),
		field.String("code").NotEmpty(),
		field.Time("expireDate").StorageKey("expire_date"),
		field.Bool("used").Default(false),
	}
}

// Edges of the VerificationCode.
func (VerificationCode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("codes").Required().Unique(),
	}
}
