package expr

import "go.mongodb.org/mongo-driver/bson"

// --- Trigonometric ---

// Sin builds a $sin expression (input in radians).
func Sin(e Expr) Expr { return rawExpr{bson.D{{Key: "$sin", Value: e.Build()}}} }

// Cos builds a $cos expression (input in radians).
func Cos(e Expr) Expr { return rawExpr{bson.D{{Key: "$cos", Value: e.Build()}}} }

// Tan builds a $tan expression (input in radians).
func Tan(e Expr) Expr { return rawExpr{bson.D{{Key: "$tan", Value: e.Build()}}} }

// Asin builds a $asin expression. Returns angle in radians.
func Asin(e Expr) Expr { return rawExpr{bson.D{{Key: "$asin", Value: e.Build()}}} }

// Acos builds a $acos expression. Returns angle in radians.
func Acos(e Expr) Expr { return rawExpr{bson.D{{Key: "$acos", Value: e.Build()}}} }

// Atan builds a $atan expression. Returns angle in radians.
func Atan(e Expr) Expr { return rawExpr{bson.D{{Key: "$atan", Value: e.Build()}}} }

// Atan2 builds a $atan2 expression: atan2(y, x). Returns angle in radians.
func Atan2(y, x Expr) Expr {
	return rawExpr{bson.D{{Key: "$atan2", Value: bson.A{y.Build(), x.Build()}}}}
}

// Sinh builds a $sinh expression (hyperbolic sine).
func Sinh(e Expr) Expr { return rawExpr{bson.D{{Key: "$sinh", Value: e.Build()}}} }

// Cosh builds a $cosh expression (hyperbolic cosine).
func Cosh(e Expr) Expr { return rawExpr{bson.D{{Key: "$cosh", Value: e.Build()}}} }

// Tanh builds a $tanh expression (hyperbolic tangent).
func Tanh(e Expr) Expr { return rawExpr{bson.D{{Key: "$tanh", Value: e.Build()}}} }

// Asinh builds a $asinh expression (inverse hyperbolic sine).
func Asinh(e Expr) Expr { return rawExpr{bson.D{{Key: "$asinh", Value: e.Build()}}} }

// Acosh builds a $acosh expression (inverse hyperbolic cosine).
func Acosh(e Expr) Expr { return rawExpr{bson.D{{Key: "$acosh", Value: e.Build()}}} }

// Atanh builds a $atanh expression (inverse hyperbolic tangent).
func Atanh(e Expr) Expr { return rawExpr{bson.D{{Key: "$atanh", Value: e.Build()}}} }

// --- Angle Unit Conversion ---

// DegreesToRadians builds a $degreesToRadians expression.
func DegreesToRadians(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$degreesToRadians", Value: e.Build()}}}
}

// RadiansToDegrees builds a $radiansToDegrees expression.
func RadiansToDegrees(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$radiansToDegrees", Value: e.Build()}}}
}
