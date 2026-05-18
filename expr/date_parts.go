package expr

import "go.mongodb.org/mongo-driver/bson"

// --- ISO Date Operators ---

// Week builds a $week expression (week of year, 0–53).
func Week(date Expr) Expr { return rawExpr{bson.D{{Key: "$week", Value: date.Build()}}} }

// ISOWeek builds a $isoWeek expression (ISO 8601 week, 1–53).
func ISOWeek(date Expr) Expr { return rawExpr{bson.D{{Key: "$isoWeek", Value: date.Build()}}} }

// ISOWeekYear builds a $isoWeekYear expression (ISO 8601 year, matches ISOWeek).
func ISOWeekYear(date Expr) Expr {
	return rawExpr{bson.D{{Key: "$isoWeekYear", Value: date.Build()}}}
}

// ISODayOfWeek builds a $isoDayOfWeek expression (ISO 8601, 1=Monday … 7=Sunday).
func ISODayOfWeek(date Expr) Expr {
	return rawExpr{bson.D{{Key: "$isoDayOfWeek", Value: date.Build()}}}
}

// --- Date Construction / Deconstruction ---

// DateFromPartsOptions holds optional fields for $dateFromParts.
type DateFromPartsOptions struct {
	Year        Expr // required (or ISOWeekYear)
	ISOWeekYear Expr // ISO alternative to Year
	Month       Expr // 1–12 (gregorian), or ISOWeek
	ISOWeek     Expr
	Day         Expr // 1–31 (gregorian), or ISODayOfWeek
	ISODayOfWeek Expr
	Hour        Expr
	Minute      Expr
	Second      Expr
	Millisecond Expr
	Timezone    Expr
}

// DateFromParts builds a $dateFromParts expression.
// Use either the gregorian (Year/Month/Day) or ISO week-date (ISOWeekYear/ISOWeek/ISODayOfWeek) fields.
func DateFromParts(opts DateFromPartsOptions) Expr {
	doc := bson.D{}
	appendExpr := func(key string, e Expr) {
		if e != nil {
			doc = append(doc, bson.E{Key: key, Value: e.Build()})
		}
	}
	appendExpr("year", opts.Year)
	appendExpr("isoWeekYear", opts.ISOWeekYear)
	appendExpr("month", opts.Month)
	appendExpr("isoWeek", opts.ISOWeek)
	appendExpr("day", opts.Day)
	appendExpr("isoDayOfWeek", opts.ISODayOfWeek)
	appendExpr("hour", opts.Hour)
	appendExpr("minute", opts.Minute)
	appendExpr("second", opts.Second)
	appendExpr("millisecond", opts.Millisecond)
	appendExpr("timezone", opts.Timezone)
	return rawExpr{bson.D{{Key: "$dateFromParts", Value: doc}}}
}

// DateToPartsOptions holds optional fields for $dateToParts.
type DateToPartsOptions struct {
	// ISO if true returns isoWeekYear/isoWeek/isoDayOfWeek instead of year/month/day.
	ISO      bool
	Timezone Expr
}

// DateToParts builds a $dateToParts expression.
func DateToParts(date Expr, opts ...DateToPartsOptions) Expr {
	doc := bson.D{{Key: "date", Value: date.Build()}}
	if len(opts) > 0 {
		if opts[0].Timezone != nil {
			doc = append(doc, bson.E{Key: "timezone", Value: opts[0].Timezone.Build()})
		}
		if opts[0].ISO {
			doc = append(doc, bson.E{Key: "iso8601", Value: true})
		}
	}
	return rawExpr{bson.D{{Key: "$dateToParts", Value: doc}}}
}
