package stage

import (
	"github.com/rick/mono-query/expr"
	"go.mongodb.org/mongo-driver/bson"
)

// FieldExpr pairs a field name with an aggregation expression.
type FieldExpr struct {
	Field string
	Expr  expr.Expr
}

// addFieldsStage implements Stage for $addFields / $set.
type addFieldsStage struct {
	op     string
	fields []FieldExpr
}

// AddFields creates an $addFields stage.
//
//	stage.AddFields(
//	    stage.FieldExpr{"total", expr.Sum(expr.Field("items"))},
//	    stage.FieldExpr{"active", expr.Filter(expr.Field("items"), "i", expr.Eq(expr.Var("i.status"), expr.Value("active")))},
//	)
func AddFields(fields ...FieldExpr) Stage {
	return addFieldsStage{op: "$addFields", fields: fields}
}

// Set creates a $set stage (alias for $addFields, available in MongoDB 4.2+).
func Set(fields ...FieldExpr) Stage {
	return addFieldsStage{op: "$set", fields: fields}
}

// FE is a short-hand constructor for FieldExpr.
//
//	stage.AddFields(stage.FE("discounted", expr.Multiply(expr.Field("price"), expr.Value(0.9))))
func FE(field string, e expr.Expr) FieldExpr {
	mustNotEmpty(field, "addFields field")
	return FieldExpr{Field: field, Expr: e}
}

func (a addFieldsStage) Build() bson.D {
	doc := make(bson.D, len(a.fields))
	for i, f := range a.fields {
		doc[i] = bson.E{Key: f.Field, Value: f.Expr.Build()}
	}
	return bson.D{{Key: a.op, Value: doc}}
}
