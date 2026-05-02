package expr

import "go.mongodb.org/mongo-driver/bson"

// Cond builds a $cond expression.
//
//	expr.Cond(
//	    expr.Eq(expr.Field("status"), expr.Value("active")),
//	    expr.Field("price"),
//	    expr.Value(0),
//	)
//	→ { $cond: { if: {$eq:[...]}, then: "$price", else: 0 } }
func Cond(ifExpr, thenExpr, elseExpr Expr) Expr {
	return rawExpr{bson.D{{Key: "$cond", Value: bson.D{
		{Key: "if", Value: ifExpr.Build()},
		{Key: "then", Value: thenExpr.Build()},
		{Key: "else", Value: elseExpr.Build()},
	}}}}
}

// IfNull builds a $ifNull expression.
// Returns expr if it is non-null, otherwise returns replacement.
func IfNull(e, replacement Expr) Expr {
	return rawExpr{bson.D{{Key: "$ifNull", Value: bson.A{e.Build(), replacement.Build()}}}}
}

// SwitchBranch is a single case/then pair for a $switch expression.
type SwitchBranch struct {
	Case Expr
	Then Expr
}

// Switch builds a $switch expression.
// defaultExpr may be nil if all cases are guaranteed to be covered.
//
//	expr.Switch(
//	    []expr.SwitchBranch{
//	        {expr.Eq(expr.Field("status"), expr.Value("vip")), expr.Value(0.1)},
//	        {expr.Eq(expr.Field("status"), expr.Value("premium")), expr.Value(0.05)},
//	    },
//	    expr.Value(0),
//	)
func Switch(branches []SwitchBranch, defaultExpr Expr) Expr {
	brs := make(bson.A, len(branches))
	for i, b := range branches {
		brs[i] = bson.D{
			{Key: "case", Value: b.Case.Build()},
			{Key: "then", Value: b.Then.Build()},
		}
	}
	doc := bson.D{{Key: "branches", Value: brs}}
	if defaultExpr != nil {
		doc = append(doc, bson.E{Key: "default", Value: defaultExpr.Build()})
	}
	return rawExpr{bson.D{{Key: "$switch", Value: doc}}}
}
