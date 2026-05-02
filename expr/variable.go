package expr

import "go.mongodb.org/mongo-driver/bson"

// Pre-defined MongoDB system variables.
var (
	// Root refers to $$ROOT — the root document of the pipeline stage.
	Root Expr = varExpr{"ROOT"}
	// Remove refers to $$REMOVE — used to conditionally exclude a field.
	Remove Expr = varExpr{"REMOVE"}
	// Now refers to $$NOW — the current datetime.
	Now Expr = varExpr{"NOW"}
	// Current refers to $$CURRENT — the current document being processed.
	Current Expr = varExpr{"CURRENT"}
)

// LetBinding represents a single variable binding in a $let expression.
type LetBinding struct {
	Name string
	Expr Expr
}

// Let constructs a $let expression that binds variables for use in an inner
// expression.
//
//	expr.Let(
//	    []expr.LetBinding{{"discount", expr.Field("$discount")}},
//	    expr.Multiply(expr.Field("$price"), expr.Var("discount")),
//	)
func Let(vars []LetBinding, in Expr) Expr {
	varsDoc := make(bson.D, len(vars))
	for i, v := range vars {
		varsDoc[i] = bson.E{Key: v.Name, Value: v.Expr.Build()}
	}
	return rawExpr{bson.D{
		{Key: "$let", Value: bson.D{
			{Key: "vars", Value: varsDoc},
			{Key: "in", Value: in.Build()},
		}},
	}}
}
