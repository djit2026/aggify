package expr

import "go.mongodb.org/mongo-driver/bson"

// --- Dynamic Field Access (MongoDB 5.0+) ---

// GetField builds a $getField expression for safe dynamic field access.
// Useful when field names contain dots or dollar signs.
//
//	expr.GetField("price.usd", expr.Root)
func GetField(field string, input Expr) Expr {
	doc := bson.D{{Key: "field", Value: field}}
	if input != nil {
		doc = append(doc, bson.E{Key: "input", Value: input.Build()})
	}
	return rawExpr{bson.D{{Key: "$getField", Value: doc}}}
}

// GetFieldSimple builds a simple $getField expression using just the field name
// (MongoDB resolves it against $$CURRENT).
func GetFieldSimple(field string) Expr {
	return rawExpr{bson.D{{Key: "$getField", Value: field}}}
}

// SetField builds a $setField expression that adds or updates a field in a document.
//
//	expr.SetField("price", expr.Root, expr.Value(99))
func SetField(field string, input, value Expr) Expr {
	return rawExpr{bson.D{{Key: "$setField", Value: bson.D{
		{Key: "field", Value: field},
		{Key: "input", Value: input.Build()},
		{Key: "value", Value: value.Build()},
	}}}}
}

// UnsetField builds a $unsetField expression that removes a field from a document.
//
//	expr.UnsetField("password", expr.Root)
func UnsetField(field string, input Expr) Expr {
	return rawExpr{bson.D{{Key: "$unsetField", Value: bson.D{
		{Key: "field", Value: field},
		{Key: "input", Value: input.Build()},
	}}}}
}
