package expr

import "go.mongodb.org/mongo-driver/bson"

// Filter builds a $filter expression that selects elements from an array.
//
//	expr.Filter(
//	    expr.Field("items"), "item",
//	    expr.Gt(expr.Var("item.qty"), expr.Value(5)),
//	)
func Filter(input Expr, as string, cond Expr) Expr {
	return rawExpr{bson.D{{Key: "$filter", Value: bson.D{
		{Key: "input", Value: input.Build()},
		{Key: "as", Value: as},
		{Key: "cond", Value: cond.Build()},
	}}}}
}

// Map builds a $map expression that applies an expression to each element.
func Map(input Expr, as string, in Expr) Expr {
	return rawExpr{bson.D{{Key: "$map", Value: bson.D{
		{Key: "input", Value: input.Build()},
		{Key: "as", Value: as},
		{Key: "in", Value: in.Build()},
	}}}}
}

// Reduce builds a $reduce expression.
func Reduce(input, initialValue, in Expr) Expr {
	return rawExpr{bson.D{{Key: "$reduce", Value: bson.D{
		{Key: "input", Value: input.Build()},
		{Key: "initialValue", Value: initialValue.Build()},
		{Key: "in", Value: in.Build()},
	}}}}
}

// Size builds a $size expression that returns the number of elements in an array.
func Size(arr Expr) Expr {
	return rawExpr{bson.D{{Key: "$size", Value: arr.Build()}}}
}

// First returns the first element of an array ($first as array operator).
func First(arr Expr) Expr {
	return rawExpr{bson.D{{Key: "$first", Value: arr.Build()}}}
}

// Last returns the last element of an array ($last as array operator).
func Last(arr Expr) Expr {
	return rawExpr{bson.D{{Key: "$last", Value: arr.Build()}}}
}

// ArrayElemAt returns the element at the given index of an array.
func ArrayElemAt(arr, idx Expr) Expr {
	return rawExpr{bson.D{{Key: "$arrayElemAt", Value: bson.A{arr.Build(), idx.Build()}}}}
}

// ArrIn builds a $in aggregation expression (array membership check).
// Note: this is the aggregation $in — use q.In for query filters.
func ArrIn(elem, arr Expr) Expr {
	return rawExpr{bson.D{{Key: "$in", Value: bson.A{elem.Build(), arr.Build()}}}}
}

// ConcatArrays builds a $concatArrays expression.
func ConcatArrays(arrs ...Expr) Expr {
	return rawExpr{bson.D{{Key: "$concatArrays", Value: BuildAll(arrs)}}}
}

// Zip builds a $zip expression.
func Zip(inputs []Expr, useLongest bool, defaults Expr) Expr {
	doc := bson.D{
		{Key: "inputs", Value: BuildAll(inputs)},
		{Key: "useLongestLength", Value: useLongest},
	}
	if defaults != nil {
		doc = append(doc, bson.E{Key: "defaults", Value: defaults.Build()})
	}
	return rawExpr{bson.D{{Key: "$zip", Value: doc}}}
}

// Slice builds a $slice expression that returns a subset of an array.
func Slice(arr, n Expr) Expr {
	return rawExpr{bson.D{{Key: "$slice", Value: bson.A{arr.Build(), n.Build()}}}}
}

// SliceFrom builds a $slice expression with position and n.
func SliceFrom(arr, position, n Expr) Expr {
	return rawExpr{bson.D{{Key: "$slice", Value: bson.A{arr.Build(), position.Build(), n.Build()}}}}
}

// --- Accumulator expressions (used inside $group) ---

// Push builds a $push accumulator expression.
func Push(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$push", Value: e.Build()}}}
}

// AddToSet builds a $addToSet accumulator expression.
func AddToSet(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$addToSet", Value: e.Build()}}}
}

// Sum builds a $sum expression (accumulator or arithmetic).
func Sum(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$sum", Value: e.Build()}}}
}

// Avg builds a $avg accumulator expression.
func Avg(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$avg", Value: e.Build()}}}
}

// Min builds a $min accumulator expression.
func Min(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$min", Value: e.Build()}}}
}

// Max builds a $max accumulator expression.
func Max(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$max", Value: e.Build()}}}
}

// Count builds a $count accumulator expression (MongoDB 5.0+).
func Count() Expr {
	return rawExpr{bson.D{{Key: "$count", Value: bson.D{}}}}
}
