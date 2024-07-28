package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// UserWord holds the schema definition for the UserWord entity.
type UserWord struct {
	ent.Schema
}

func (UserWord) Mixin() []ent.Mixin {
	return []ent.Mixin{
		CommonMixin{},
	}
}

// Fields of the UserWord.
func (UserWord) Fields() []ent.Field {
	return []ent.Field{
		field.String("term").NotEmpty(),
		field.JSON("definitions", []Definition{}).Default([]Definition{}),
		field.Float("learningProgress").Default(0.0).Max(100).Min(0),
	}
}

// Edges of the UserWord.
func (UserWord) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("lang", Language.Type).
			Ref("userWords").
			Unique(),
		edge.From("user", User.Type).
			Ref("userWords").
			Unique(),
	}
}
