package expr

import "go.mongodb.org/mongo-driver/bson"

// SetEquals builds a $setEquals expression.
func SetEquals(exprs ...Expr) Expr {
	return rawExpr{bson.D{{Key: "$setEquals", Value: BuildAll(exprs)}}}
}

// SetIntersection builds a $setIntersection expression.
func SetIntersection(exprs ...Expr) Expr {
	return rawExpr{bson.D{{Key: "$setIntersection", Value: BuildAll(exprs)}}}
}

// SetUnion builds a $setUnion expression.
func SetUnion(exprs ...Expr) Expr {
	return rawExpr{bson.D{{Key: "$setUnion", Value: BuildAll(exprs)}}}
}

// SetDifference builds a $setDifference expression: { $setDifference: [setA, setB] }
func SetDifference(a, b Expr) Expr {
	return rawExpr{bson.D{{Key: "$setDifference", Value: bson.A{a.Build(), b.Build()}}}}
}

// SetIsSubset builds a $setIsSubset expression: { $setIsSubset: [setA, setB] }
func SetIsSubset(a, b Expr) Expr {
	return rawExpr{bson.D{{Key: "$setIsSubset", Value: bson.A{a.Build(), b.Build()}}}}
}

// AnyElementTrue builds a $anyElementTrue expression.
func AnyElementTrue(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$anyElementTrue", Value: e.Build()}}}
}

// AllElementsTrue builds a $allElementsTrue expression.
func AllElementsTrue(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$allElementsTrue", Value: e.Build()}}}
}
