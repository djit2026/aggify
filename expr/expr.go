// Package expr provides a composable expression engine for MongoDB aggregation
// pipelines. Every function in this package returns an Expr that compiles
// directly to a BSON value — no reflection, no runtime parsing.
//
// Typical usage:
//
//	expr.Cond(
//	    expr.Eq(expr.Field("status"), expr.Value("active")),
//	    expr.Field("price"),
//	    expr.Value(0),
//	)
package expr

import "go.mongodb.org/mongo-driver/bson"

// Expr is any MongoDB aggregation expression.
// It compiles to a BSON-compatible value (bson.D, bson.A, string, int, …).
type Expr interface {
	Build() any
}

// Field references a document field in an aggregation expression.
// The "$" prefix is added automatically.
//
//	expr.Field("price")        → "$price"
//	expr.Field("items.status") → "$items.status"
func Field(name string) Expr { return fieldExpr{name} }

// Var references a pipeline variable using the "$$" prefix.
//
//	expr.Var("item") → "$$item"
func Var(name string) Expr { return varExpr{name} }

// Value wraps a literal Go value as an aggregation expression.
//
//	expr.Value(42)       → 42
//	expr.Value("active") → "active"
func Value(v any) Expr { return valueExpr{v} }

// Raw wraps an already-built BSON value as an Expr.
// Use this as an escape hatch when you need to embed a hand-crafted bson.D.
func Raw(v any) Expr { return rawExpr{v} }

// BuildAll converts a slice of Exprs to a bson.A (BSON array).
func BuildAll(exprs []Expr) bson.A {
	out := make(bson.A, len(exprs))
	for i, e := range exprs {
		out[i] = e.Build()
	}
	return out
}

// --- internal primitives ----------------------------------------------------

type fieldExpr struct{ name string }

func (f fieldExpr) Build() any { return "$" + f.name }

type varExpr struct{ name string }

func (v varExpr) Build() any { return "$$" + v.name }

type valueExpr struct{ v any }

func (v valueExpr) Build() any { return v.v }

type rawExpr struct{ v any }

func (r rawExpr) Build() any { return r.v }
