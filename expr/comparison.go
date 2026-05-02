package expr

import "go.mongodb.org/mongo-driver/bson"

// binaryExpr handles two-operand operators that take an array argument.
type binaryExpr struct {
	op  string
	lhs Expr
	rhs Expr
}

func (b binaryExpr) Build() any {
	return bson.D{{Key: b.op, Value: bson.A{b.lhs.Build(), b.rhs.Build()}}}
}

// Eq builds a $eq aggregation expression: { $eq: [lhs, rhs] }
func Eq(lhs, rhs Expr) Expr { return binaryExpr{"$eq", lhs, rhs} }

// Ne builds a $ne aggregation expression: { $ne: [lhs, rhs] }
func Ne(lhs, rhs Expr) Expr { return binaryExpr{"$ne", lhs, rhs} }

// Gt builds a $gt aggregation expression: { $gt: [lhs, rhs] }
func Gt(lhs, rhs Expr) Expr { return binaryExpr{"$gt", lhs, rhs} }

// Gte builds a $gte aggregation expression: { $gte: [lhs, rhs] }
func Gte(lhs, rhs Expr) Expr { return binaryExpr{"$gte", lhs, rhs} }

// Lt builds a $lt aggregation expression: { $lt: [lhs, rhs] }
func Lt(lhs, rhs Expr) Expr { return binaryExpr{"$lt", lhs, rhs} }

// Lte builds a $lte aggregation expression: { $lte: [lhs, rhs] }
func Lte(lhs, rhs Expr) Expr { return binaryExpr{"$lte", lhs, rhs} }

// Cmp builds a $cmp aggregation expression: { $cmp: [lhs, rhs] }
// Returns -1, 0, or 1.
func Cmp(lhs, rhs Expr) Expr { return binaryExpr{"$cmp", lhs, rhs} }
