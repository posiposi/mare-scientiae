package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Book struct {
	ent.Schema
}

func (Book) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("google_books_id").
			MaxLen(50).
			Unique(),
		field.String("title").
			MaxLen(500),
		field.String("subtitle").
			MaxLen(500).
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (Book) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("authors", Author.Type),
	}
}
