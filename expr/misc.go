package expr

import "go.mongodb.org/mongo-driver/bson"

// --- Miscellaneous Operators ---

// Literal builds a $literal expression.
// Forces the value to be treated as a literal, not an expression.
// Useful when a value starts with "$" or would otherwise be interpreted.
//
//	expr.Literal("$notAField") → { $literal: "$notAField" }
func Literal(v any) Expr {
	return rawExpr{bson.D{{Key: "$literal", Value: v}}}
}

// Rand builds a $rand expression that returns a random float between 0 and 1.
func Rand() Expr {
	return rawExpr{bson.D{{Key: "$rand", Value: bson.D{}}}}
}

// --- Timestamp Operators (MongoDB 5.1+) ---

// TsSecond builds a $tsSecond expression.
// Returns the seconds component of a BSON Timestamp.
func TsSecond(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$tsSecond", Value: e.Build()}}}
}

// TsIncrement builds a $tsIncrement expression.
// Returns the incrementing ordinal of a BSON Timestamp.
func TsIncrement(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$tsIncrement", Value: e.Build()}}}
}
