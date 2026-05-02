package expr

import "go.mongodb.org/mongo-driver/bson"

// MergeObjects builds a $mergeObjects expression.
// When called with multiple expressions, they are merged left-to-right.
// When called with one expression (e.g. an array field), it merges elements of that array.
func MergeObjects(exprs ...Expr) Expr {
	if len(exprs) == 1 {
		return rawExpr{bson.D{{Key: "$mergeObjects", Value: exprs[0].Build()}}}
	}
	return rawExpr{bson.D{{Key: "$mergeObjects", Value: BuildAll(exprs)}}}
}

// ObjectToArray builds a $objectToArray expression.
// Converts a document to an array of {k, v} documents.
func ObjectToArray(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$objectToArray", Value: e.Build()}}}
}

// ArrayToObject builds a $arrayToObject expression.
// Converts an array of {k, v} documents (or [[key, value]] pairs) to an object.
func ArrayToObject(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$arrayToObject", Value: e.Build()}}}
}
