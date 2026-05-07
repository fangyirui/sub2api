package schema

import (
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type SmsOrder struct {
	ent.Schema
}

func (SmsOrder) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sms_orders"},
	}
}

func (SmsOrder) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (SmsOrder) Fields() []ent.Field {
	return []ent.Field{
		field.String("order_no").
			MaxLen(100).
			NotEmpty().
			Unique(),
		field.String("phone_number").
			MaxLen(50).
			NotEmpty(),
		field.String("sms_content").
			SchemaType(map[string]string{
				dialect.Postgres: "text",
			}).
			Optional().
			Default(""),
		field.String("status").
			MaxLen(20).
			Default("pending"),
	}
}

func (SmsOrder) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
	}
}
