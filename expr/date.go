package expr

import "go.mongodb.org/mongo-driver/bson"

// Year builds a $year expression.
func Year(date Expr) Expr {
	return rawExpr{bson.D{{Key: "$year", Value: date.Build()}}}
}

// Month builds a $month expression.
func Month(date Expr) Expr {
	return rawExpr{bson.D{{Key: "$month", Value: date.Build()}}}
}

// DayOfMonth builds a $dayOfMonth expression.
func DayOfMonth(date Expr) Expr {
	return rawExpr{bson.D{{Key: "$dayOfMonth", Value: date.Build()}}}
}

// DayOfWeek builds a $dayOfWeek expression.
func DayOfWeek(date Expr) Expr {
	return rawExpr{bson.D{{Key: "$dayOfWeek", Value: date.Build()}}}
}

// DayOfYear builds a $dayOfYear expression.
func DayOfYear(date Expr) Expr {
	return rawExpr{bson.D{{Key: "$dayOfYear", Value: date.Build()}}}
}

// Hour builds a $hour expression.
func Hour(date Expr) Expr {
	return rawExpr{bson.D{{Key: "$hour", Value: date.Build()}}}
}

// Minute builds a $minute expression.
func Minute(date Expr) Expr {
	return rawExpr{bson.D{{Key: "$minute", Value: date.Build()}}}
}

// Second builds a $second expression.
func Second(date Expr) Expr {
	return rawExpr{bson.D{{Key: "$second", Value: date.Build()}}}
}

// Millisecond builds a $millisecond expression.
func Millisecond(date Expr) Expr {
	return rawExpr{bson.D{{Key: "$millisecond", Value: date.Build()}}}
}

// DateAdd builds a $dateAdd expression.
// amount is a number expression, unit is a string (e.g., "day", "month").
func DateAdd(startDate Expr, amount Expr, unit string) Expr {
	return rawExpr{bson.D{{Key: "$dateAdd", Value: bson.D{
		{Key: "startDate", Value: startDate.Build()},
		{Key: "amount", Value: amount.Build()},
		{Key: "unit", Value: unit},
	}}}}
}

// DateSubtract builds a $dateSubtract expression.
func DateSubtract(startDate Expr, amount Expr, unit string) Expr {
	return rawExpr{bson.D{{Key: "$dateSubtract", Value: bson.D{
		{Key: "startDate", Value: startDate.Build()},
		{Key: "amount", Value: amount.Build()},
		{Key: "unit", Value: unit},
	}}}}
}

// DateDiff builds a $dateDiff expression.
func DateDiff(startDate, endDate Expr, unit string) Expr {
	return rawExpr{bson.D{{Key: "$dateDiff", Value: bson.D{
		{Key: "startDate", Value: startDate.Build()},
		{Key: "endDate", Value: endDate.Build()},
		{Key: "unit", Value: unit},
	}}}}
}

// DateTrunc builds a $dateTrunc expression.
// binSize is optional.
func DateTrunc(date Expr, unit string, binSize ...int) Expr {
	doc := bson.D{
		{Key: "date", Value: date.Build()},
		{Key: "unit", Value: unit},
	}
	if len(binSize) > 0 && binSize[0] > 1 {
		doc = append(doc, bson.E{Key: "binSize", Value: binSize[0]})
	}
	return rawExpr{bson.D{{Key: "$dateTrunc", Value: doc}}}
}

// DateToString builds a $dateToString expression.
func DateToString(date Expr, format string) Expr {
	return rawExpr{bson.D{{Key: "$dateToString", Value: bson.D{
		{Key: "date", Value: date.Build()},
		{Key: "format", Value: format},
	}}}}
}

// DateFromString builds a $dateFromString expression.
// format is optional.
func DateFromString(dateString Expr, format ...string) Expr {
	doc := bson.D{{Key: "dateString", Value: dateString.Build()}}
	if len(format) > 0 && format[0] != "" {
		doc = append(doc, bson.E{Key: "format", Value: format[0]})
	}
	return rawExpr{bson.D{{Key: "$dateFromString", Value: doc}}}
}
