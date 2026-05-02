// Package q provides query filter helpers for MongoDB.
// These helpers produce bson.D filter documents suitable for use in
// collection.Find(), collection.FindOne(), and stage.Match().
//
// Note: q helpers are NOT aggregation expressions. For aggregation expressions
// (used inside $group, $project, $addFields, etc.) use the expr package.
//
// Example:
//
//	filter := q.And(
//	    q.Eq("status", "active"),
//	    q.Gte("age", 18),
//	    q.In("role", "admin", "editor"),
//	)
package q

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Eq builds an equality filter. MongoDB's implicit equality form is used
// (no $eq operator), which is the idiomatic and most performant form.
//
//	q.Eq("status", "active") → bson.D{{"status", "active"}}
func Eq(field string, value any) bson.D {
	return bson.D{{Key: field, Value: value}}
}

// Ne builds a $ne filter: { field: { $ne: value } }
func Ne(field string, value any) bson.D {
	return bson.D{{Key: field, Value: bson.D{{Key: "$ne", Value: value}}}}
}

// Gt builds a $gt filter: { field: { $gt: value } }
func Gt(field string, value any) bson.D {
	return bson.D{{Key: field, Value: bson.D{{Key: "$gt", Value: value}}}}
}

// Gte builds a $gte filter: { field: { $gte: value } }
func Gte(field string, value any) bson.D {
	return bson.D{{Key: field, Value: bson.D{{Key: "$gte", Value: value}}}}
}

// Lt builds a $lt filter: { field: { $lt: value } }
func Lt(field string, value any) bson.D {
	return bson.D{{Key: field, Value: bson.D{{Key: "$lt", Value: value}}}}
}

// Lte builds a $lte filter: { field: { $lte: value } }
func Lte(field string, value any) bson.D {
	return bson.D{{Key: field, Value: bson.D{{Key: "$lte", Value: value}}}}
}

// In builds a $in filter: { field: { $in: [values...] } }
func In(field string, values ...any) bson.D {
	arr := make(bson.A, len(values))
	copy(arr, values)
	return bson.D{{Key: field, Value: bson.D{{Key: "$in", Value: arr}}}}
}

// Nin builds a $nin filter: { field: { $nin: [values...] } }
func Nin(field string, values ...any) bson.D {
	arr := make(bson.A, len(values))
	copy(arr, values)
	return bson.D{{Key: field, Value: bson.D{{Key: "$nin", Value: arr}}}}
}

// Exists builds an $exists filter: { field: { $exists: exists } }
func Exists(field string, exists bool) bson.D {
	return bson.D{{Key: field, Value: bson.D{{Key: "$exists", Value: exists}}}}
}

// Regex builds a $regex filter using MongoDB's native regex type.
//
//	q.Regex("email", `^admin`, "i")
func Regex(field, pattern, options string) bson.D {
	return bson.D{{Key: field, Value: primitive.Regex{Pattern: pattern, Options: options}}}
}

// ElemMatch builds an $elemMatch filter.
//
//	q.ElemMatch("items", q.And(q.Eq("status", "active"), q.Gt("qty", 0)))
func ElemMatch(field string, filter bson.D) bson.D {
	return bson.D{{Key: field, Value: bson.D{{Key: "$elemMatch", Value: filter}}}}
}

// All builds an $all filter.
func All(field string, values ...any) bson.D {
	arr := make(bson.A, len(values))
	copy(arr, values)
	return bson.D{{Key: field, Value: bson.D{{Key: "$all", Value: arr}}}}
}

// ArrSize builds an array $size filter: { field: { $size: n } }
func ArrSize(field string, size int) bson.D {
	return bson.D{{Key: field, Value: bson.D{{Key: "$size", Value: size}}}}
}

// Type builds a $type filter.
func Type(field string, bsonType string) bson.D {
	return bson.D{{Key: field, Value: bson.D{{Key: "$type", Value: bsonType}}}}
}

// Expr wraps an aggregation expression for use in a query filter ($expr).
//
//	q.Expr(expr.Gt(expr.Field("spend"), expr.Field("budget")))
func Expr(e interface{ Build() any }) bson.D {
	return bson.D{{Key: "$expr", Value: e.Build()}}
}
