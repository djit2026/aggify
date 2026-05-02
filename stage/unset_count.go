package stage

import "go.mongodb.org/mongo-driver/bson"

type unsetStage struct{ fields []string }

// Unset creates an $unset stage that removes the specified fields.
//
//	stage.Unset("password", "internalNotes")
func Unset(fields ...string) Stage {
	for _, f := range fields {
		mustNotEmpty(f, "unset field")
	}
	return unsetStage{fields}
}

func (u unsetStage) Build() bson.D {
	if len(u.fields) == 1 {
		return bson.D{{Key: "$unset", Value: u.fields[0]}}
	}
	arr := make(bson.A, len(u.fields))
	for i, f := range u.fields {
		arr[i] = f
	}
	return bson.D{{Key: "$unset", Value: arr}}
}

type countStage struct{ field string }

// Count creates a $count stage that counts the documents and stores the
// result in the given field.
//
//	stage.Count("total")
//	→ { $count: "total" }
func Count(field string) Stage {
	mustNotEmpty(field, "count field")
	return countStage{field}
}

func (c countStage) Build() bson.D {
	return bson.D{{Key: "$count", Value: c.field}}
}
