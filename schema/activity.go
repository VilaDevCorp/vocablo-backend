package schema

// import (
// 	"entgo.io/ent"
// 	"entgo.io/ent/schema/edge"
// 	"entgo.io/ent/schema/field"
// )

// // Activity holds the schema definition for the Activity entity.
// type Activity struct {
// 	ent.Schema
// }

// func (Activity) Mixin() []ent.Mixin {
// 	return []ent.Mixin{
// 		CommonMixin{},
// 	}
// }

// // Fields of the Activity.
// func (Activity) Fields() []ent.Field {
// 	return []ent.Field{
// 		field.String("name").NotEmpty(),
// 		field.String("description").Optional(),
// 		field.String("icon").Optional(),
// 		field.Int("size").Default(1),
// 	}
// }

// // Edges of the Activity.
// func (Activity) Edges() []ent.Edge {
// 	return []ent.Edge{
// 		edge.From("user", User.Type).Ref("activities").Required().Unique(),
// 	}
// }
