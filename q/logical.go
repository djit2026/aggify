package q

import "go.mongodb.org/mongo-driver/bson"

// And combines multiple filters with $and.
//
//	q.And(q.Eq("active", true), q.Gt("age", 18))
//	→ { $and: [{"active": true}, {"age": {"$gt": 18}}] }
func And(filters ...bson.D) bson.D {
	docs := make(bson.A, len(filters))
	for i, f := range filters {
		docs[i] = f
	}
	return bson.D{{Key: "$and", Value: docs}}
}

// Or combines multiple filters with $or.
func Or(filters ...bson.D) bson.D {
	docs := make(bson.A, len(filters))
	for i, f := range filters {
		docs[i] = f
	}
	return bson.D{{Key: "$or", Value: docs}}
}

// Nor combines multiple filters with $nor.
func Nor(filters ...bson.D) bson.D {
	docs := make(bson.A, len(filters))
	for i, f := range filters {
		docs[i] = f
	}
	return bson.D{{Key: "$nor", Value: docs}}
}

// Not wraps a single field condition with $not.
//
//	q.Not("status", q.Eq("status", "banned"))
//	→ { "status": { $not: { "$eq": "banned" } } }
//
// Note: $not in query context applies to a field's condition, not a top-level filter.
// Pass the inner condition document (without the field name).
func Not(field string, condition bson.D) bson.D {
	// Extract the inner condition (the value of the first element).
	if len(condition) == 1 {
		return bson.D{{Key: field, Value: bson.D{{Key: "$not", Value: condition[0].Value}}}}
	}
	return bson.D{{Key: field, Value: bson.D{{Key: "$not", Value: condition}}}}
}
