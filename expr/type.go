package expr

import "go.mongodb.org/mongo-driver/bson"

// Type builds a $type expression that returns the BSON type of an expression.
func Type(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$type", Value: e.Build()}}}
}

// ToString builds a $toString expression.
func ToString(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$toString", Value: e.Build()}}}
}

// ToDate builds a $toDate expression.
func ToDate(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$toDate", Value: e.Build()}}}
}

// ToInt builds a $toInt expression.
func ToInt(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$toInt", Value: e.Build()}}}
}

// ToLong builds a $toLong expression.
func ToLong(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$toLong", Value: e.Build()}}}
}

// ToDouble builds a $toDouble expression.
func ToDouble(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$toDouble", Value: e.Build()}}}
}

// ToBool builds a $toBool expression.
func ToBool(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$toBool", Value: e.Build()}}}
}

// ToDecimal builds a $toDecimal expression.
func ToDecimal(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$toDecimal", Value: e.Build()}}}
}

// ToObjectID builds a $toObjectId expression.
func ToObjectID(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$toObjectId", Value: e.Build()}}}
}

// ConvertOptions holds options for the $convert expression.
type ConvertOptions struct {
	// To is the target BSON type name (e.g. "string", "int", "date").
	To string
	// OnError is the expression returned when conversion fails (optional).
	OnError Expr
	// OnNull is the expression returned when the input is null (optional).
	OnNull Expr
}

// Convert builds a $convert expression with full options.
func Convert(input Expr, opts ConvertOptions) Expr {
	doc := bson.D{
		{Key: "input", Value: input.Build()},
		{Key: "to", Value: opts.To},
	}
	if opts.OnError != nil {
		doc = append(doc, bson.E{Key: "onError", Value: opts.OnError.Build()})
	}
	if opts.OnNull != nil {
		doc = append(doc, bson.E{Key: "onNull", Value: opts.OnNull.Build()})
	}
	return rawExpr{bson.D{{Key: "$convert", Value: doc}}}
}

// IsNumber builds a $isNumber expression.
func IsNumber(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$isNumber", Value: e.Build()}}}
}
