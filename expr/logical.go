package expr

import "go.mongodb.org/mongo-driver/bson"

type logicalExpr struct {
	op    string
	exprs []Expr
}

func (l logicalExpr) Build() any {
	return bson.D{{Key: l.op, Value: BuildAll(l.exprs)}}
}

// And builds a $and aggregation expression.
func And(exprs ...Expr) Expr { return logicalExpr{"$and", exprs} }

// Or builds a $or aggregation expression.
func Or(exprs ...Expr) Expr { return logicalExpr{"$or", exprs} }

// Not builds a $not aggregation expression.
// In aggregation context $not takes a single-element array.
func Not(e Expr) Expr {
	return rawExpr{bson.D{{Key: "$not", Value: bson.A{e.Build()}}}}
}
