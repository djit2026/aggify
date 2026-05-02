package expr

import "go.mongodb.org/mongo-driver/bson"

// Add builds a $add expression.
func Add(exprs ...Expr) Expr {
	return rawExpr{bson.D{{Key: "$add", Value: BuildAll(exprs)}}}
}

// Subtract builds a $subtract expression: { $subtract: [a, b] }
func Subtract(a, b Expr) Expr {
	return rawExpr{bson.D{{Key: "$subtract", Value: bson.A{a.Build(), b.Build()}}}}
}

// Multiply builds a $multiply expression.
func Multiply(exprs ...Expr) Expr {
	return rawExpr{bson.D{{Key: "$multiply", Value: BuildAll(exprs)}}}
}

// Divide builds a $divide expression: { $divide: [a, b] }
func Divide(a, b Expr) Expr {
	return rawExpr{bson.D{{Key: "$divide", Value: bson.A{a.Build(), b.Build()}}}}
}

// Mod builds a $mod expression: { $mod: [a, b] }
func Mod(a, b Expr) Expr {
	return rawExpr{bson.D{{Key: "$mod", Value: bson.A{a.Build(), b.Build()}}}}
}

// Abs builds a $abs expression.
func Abs(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$abs", Value: e.Build()}}}
}

// Ceil builds a $ceil expression.
func Ceil(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$ceil", Value: e.Build()}}}
}

// Floor builds a $floor expression.
func Floor(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$floor", Value: e.Build()}}}
}

// Round builds a $round expression.
// The optional places argument specifies decimal places (default 0).
func Round(e Expr, places ...int) Expr {
	p := 0
	if len(places) > 0 {
		p = places[0]
	}
	return rawExpr{bson.D{{Key: "$round", Value: bson.A{e.Build(), p}}}}
}

// Sqrt builds a $sqrt expression.
func Sqrt(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$sqrt", Value: e.Build()}}}
}

// Pow builds a $pow expression: { $pow: [base, exponent] }
func Pow(base, exp Expr) Expr {
	return rawExpr{bson.D{{Key: "$pow", Value: bson.A{base.Build(), exp.Build()}}}}
}

// Log builds a $log expression: { $log: [number, base] }
func Log(number, base Expr) Expr {
	return rawExpr{bson.D{{Key: "$log", Value: bson.A{number.Build(), base.Build()}}}}
}

// Log10 builds a $log10 expression.
func Log10(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$log10", Value: e.Build()}}}
}

// Exp builds a $exp expression (e^x).
func Exp(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$exp", Value: e.Build()}}}
}

// Ln builds a $ln expression (natural log).
func Ln(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$ln", Value: e.Build()}}}
}

// Trunc builds a $trunc expression.
func Trunc(e Expr, places ...int) Expr {
	if len(places) == 0 {
		return rawExpr{bson.D{{Key: "$trunc", Value: e.Build()}}}
	}
	return rawExpr{bson.D{{Key: "$trunc", Value: bson.A{e.Build(), places[0]}}}}
}
