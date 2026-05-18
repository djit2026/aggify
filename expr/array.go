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

// --- MongoDB 5.2+ Array & Accumulator Operators ---

// SortArrayOrder is the sort direction for $sortArray.
type SortArrayOrder = int

const (
	SortArrayAsc  SortArrayOrder = 1
	SortArrayDesc SortArrayOrder = -1
)

// SortArrayField specifies a sort field and direction for $sortArray.
type SortArrayField struct {
	Field string
	Order SortArrayOrder
}

// SortArray builds a $sortArray expression (MongoDB 5.2+).
// sortBy may be a single SortArrayField for document arrays,
// or pass nil to sort a scalar array using Value(1) / Value(-1).
//
//	// Sort array of scalars descending
//	expr.SortArray(expr.Field("scores"), nil, expr.Value(-1))
//	// Sort array of documents by field
//	expr.SortArray(expr.Field("items"), []SortArrayField{{"price", SortArrayDesc}}, nil)
func SortArray(input Expr, fields []SortArrayField, scalarDir Expr) Expr {
	var sortBy any
	if len(fields) > 0 {
		doc := bson.D{}
		for _, f := range fields {
			doc = append(doc, bson.E{Key: f.Field, Value: f.Order})
		}
		sortBy = doc
	} else if scalarDir != nil {
		sortBy = scalarDir.Build()
	}
	return rawExpr{bson.D{{Key: "$sortArray", Value: bson.D{
		{Key: "input", Value: input.Build()},
		{Key: "sortBy", Value: sortBy},
	}}}}
}

// Range builds a $range expression: { $range: [start, end, step] }
// step is optional (defaults to 1 when omitted).
func Range(start, end Expr, step ...Expr) Expr {
	arr := bson.A{start.Build(), end.Build()}
	if len(step) > 0 && step[0] != nil {
		arr = append(arr, step[0].Build())
	}
	return rawExpr{bson.D{{Key: "$range", Value: arr}}}
}

// FirstN builds a $firstN accumulator expression (MongoDB 5.2+).
// Returns the first n elements of an array.
func FirstN(n, input Expr) Expr {
	return rawExpr{bson.D{{Key: "$firstN", Value: bson.D{
		{Key: "n", Value: n.Build()},
		{Key: "input", Value: input.Build()},
	}}}}
}

// LastN builds a $lastN accumulator expression (MongoDB 5.2+).
// Returns the last n elements of an array.
func LastN(n, input Expr) Expr {
	return rawExpr{bson.D{{Key: "$lastN", Value: bson.D{
		{Key: "n", Value: n.Build()},
		{Key: "input", Value: input.Build()},
	}}}}
}

// MaxN builds a $maxN accumulator expression (MongoDB 5.2+).
// Returns the n largest values.
func MaxN(n, input Expr) Expr {
	return rawExpr{bson.D{{Key: "$maxN", Value: bson.D{
		{Key: "n", Value: n.Build()},
		{Key: "input", Value: input.Build()},
	}}}}
}

// MinN builds a $minN accumulator expression (MongoDB 5.2+).
// Returns the n smallest values.
func MinN(n, input Expr) Expr {
	return rawExpr{bson.D{{Key: "$minN", Value: bson.D{
		{Key: "n", Value: n.Build()},
		{Key: "input", Value: input.Build()},
	}}}}
}

// TopSortBy holds the sort spec for $top / $topN.
type TopSortBy = bson.D

// Top builds a $top accumulator expression (MongoDB 5.2+).
// Returns the document with the highest value according to sortBy.
func Top(sortBy TopSortBy, output Expr) Expr {
	return rawExpr{bson.D{{Key: "$top", Value: bson.D{
		{Key: "sortBy", Value: sortBy},
		{Key: "output", Value: output.Build()},
	}}}}
}

// Bottom builds a $bottom accumulator expression (MongoDB 5.2+).
// Returns the document with the lowest value according to sortBy.
func Bottom(sortBy TopSortBy, output Expr) Expr {
	return rawExpr{bson.D{{Key: "$bottom", Value: bson.D{
		{Key: "sortBy", Value: sortBy},
		{Key: "output", Value: output.Build()},
	}}}}
}

// TopN builds a $topN accumulator expression (MongoDB 5.2+).
// Returns the n documents with the highest values according to sortBy.
func TopN(n Expr, sortBy TopSortBy, output Expr) Expr {
	return rawExpr{bson.D{{Key: "$topN", Value: bson.D{
		{Key: "n", Value: n.Build()},
		{Key: "sortBy", Value: sortBy},
		{Key: "output", Value: output.Build()},
	}}}}
}

// BottomN builds a $bottomN accumulator expression (MongoDB 5.2+).
// Returns the n documents with the lowest values according to sortBy.
func BottomN(n Expr, sortBy TopSortBy, output Expr) Expr {
	return rawExpr{bson.D{{Key: "$bottomN", Value: bson.D{
		{Key: "n", Value: n.Build()},
		{Key: "sortBy", Value: sortBy},
		{Key: "output", Value: output.Build()},
	}}}}
}
