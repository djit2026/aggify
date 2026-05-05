package expr

import "go.mongodb.org/mongo-driver/bson"

// Concat builds a $concat expression.
func Concat(exprs ...Expr) Expr {
	return rawExpr{bson.D{{Key: "$concat", Value: BuildAll(exprs)}}}
}

// ToLower builds a $toLower expression.
func ToLower(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$toLower", Value: e.Build()}}}
}

// ToUpper builds a $toUpper expression.
func ToUpper(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$toUpper", Value: e.Build()}}}
}

// Trim builds a $trim expression (trims both ends).
func Trim(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$trim", Value: bson.D{{Key: "input", Value: e.Build()}}}}}
}

// TrimChars builds a $trim expression with explicit chars to trim.
func TrimChars(e Expr, chars string) Expr {
	return rawExpr{bson.D{{Key: "$trim", Value: bson.D{
		{Key: "input", Value: e.Build()},
		{Key: "chars", Value: chars},
	}}}}
}

// LTrim builds a $ltrim expression.
func LTrim(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$ltrim", Value: bson.D{{Key: "input", Value: e.Build()}}}}}
}

// RTrim builds a $rtrim expression.
func RTrim(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$rtrim", Value: bson.D{{Key: "input", Value: e.Build()}}}}}
}

// SubstrCP builds a $substrCP expression.
func SubstrCP(e Expr, start, length int) Expr {
	return rawExpr{bson.D{{Key: "$substrCP", Value: bson.A{e.Build(), start, length}}}}
}

// StrLenCP builds a $strLenCP expression (character count, UTF-8 aware).
func StrLenCP(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$strLenCP", Value: e.Build()}}}
}

// StrLenBytes builds a $strLenBytes expression (byte count).
func StrLenBytes(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$strLenBytes", Value: e.Build()}}}
}

// Split builds a $split expression: { $split: [string, delimiter] }
func Split(e, delimiter Expr) Expr {
	return rawExpr{bson.D{{Key: "$split", Value: bson.A{e.Build(), delimiter.Build()}}}}
}

// IndexOfCP builds a $indexOfCP expression (finds position of substring).
func IndexOfCP(str, sub Expr) Expr {
	return rawExpr{bson.D{{Key: "$indexOfCP", Value: bson.A{str.Build(), sub.Build()}}}}
}

// RegexMatch builds a $regexMatch expression.
// options is an optional flag string (e.g. "i" for case-insensitive).
func RegexMatch(input Expr, regex string, options ...string) Expr {
	doc := bson.D{
		{Key: "input", Value: input.Build()},
		{Key: "regex", Value: regex},
	}
	if len(options) > 0 && options[0] != "" {
		doc = append(doc, bson.E{Key: "options", Value: options[0]})
	}
	return rawExpr{bson.D{{Key: "$regexMatch", Value: doc}}}
}

// RegexFind builds a $regexFind expression.
func RegexFind(input Expr, regex string, options ...string) Expr {
	doc := bson.D{
		{Key: "input", Value: input.Build()},
		{Key: "regex", Value: regex},
	}
	if len(options) > 0 && options[0] != "" {
		doc = append(doc, bson.E{Key: "options", Value: options[0]})
	}
	return rawExpr{bson.D{{Key: "$regexFind", Value: doc}}}
}

// ReplaceOne builds a $replaceOne expression.
func ReplaceOne(input, find, replacement Expr) Expr {
	return rawExpr{bson.D{{Key: "$replaceOne", Value: bson.D{
		{Key: "input", Value: input.Build()},
		{Key: "find", Value: find.Build()},
		{Key: "replacement", Value: replacement.Build()},
	}}}}
}

// ReplaceAll builds a $replaceAll expression.
func ReplaceAll(input, find, replacement Expr) Expr {
	return rawExpr{bson.D{{Key: "$replaceAll", Value: bson.D{
		{Key: "input", Value: input.Build()},
		{Key: "find", Value: find.Build()},
		{Key: "replacement", Value: replacement.Build()},
	}}}}
}

// Substr builds a $substr expression. Note: $substr is deprecated in MongoDB in favor of $substrCP, but provided for completeness.
func Substr(str Expr, start, length int) Expr {
	return rawExpr{bson.D{{Key: "$substr", Value: bson.A{str.Build(), start, length}}}}
}

// IndexOfBytes builds a $indexOfBytes expression (finds position of substring in bytes).
func IndexOfBytes(str, sub Expr) Expr {
	return rawExpr{bson.D{{Key: "$indexOfBytes", Value: bson.A{str.Build(), sub.Build()}}}}
}
