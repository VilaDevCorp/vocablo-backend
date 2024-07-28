package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Word holds the schema definition for the Word entity.
type Word struct {
	ent.Schema
}

type Definition struct {
	Definition string `json:"definition"`
	Example    string `json:"example"`
}

func (Word) Mixin() []ent.Mixin {
	return []ent.Mixin{
		CommonMixin{},
	}
}

// Fields of the Word.
func (Word) Fields() []ent.Field {
	return []ent.Field{
		field.String("term").NotEmpty(),
		field.JSON("definitions", []Definition{}).Default([]Definition{}),
	}
}

// Edges of the Word.
func (Word) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("lang", Language.Type).
			Ref("words").
			Unique()}
}
