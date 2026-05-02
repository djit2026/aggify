package stage

import (
	"github.com/djit2026/aggify/expr"
	"go.mongodb.org/mongo-driver/bson"
)

// Accumulator pairs an output field name with an accumulator expression.
type Accumulator struct {
	Field string
	Expr  expr.Expr
}

// Acc is a convenience constructor for Accumulator.
//
//	stage.Acc("total", expr.Sum(expr.Field("price")))
func Acc(field string, e expr.Expr) Accumulator {
	mustNotEmpty(field, "accumulator field")
	return Accumulator{Field: field, Expr: e}
}

// groupStage implements Stage for $group.
type groupStage struct {
	id   expr.Expr
	accs []Accumulator
}

// Group creates a $group stage.
//
//	stage.Group(
//	    expr.Field("storeId"),                        // _id
//	    stage.Acc("total", expr.Sum(expr.Field("price"))),
//	    stage.Acc("items", expr.Push(expr.Root)),
//	)
//
// Pass expr.Value(nil) as id to group all documents into a single group.
func Group(id expr.Expr, accs ...Accumulator) Stage {
	return groupStage{id: id, accs: accs}
}

func (g groupStage) Build() bson.D {
	doc := bson.D{{Key: "_id", Value: g.id.Build()}}
	for _, a := range g.accs {
		doc = append(doc, bson.E{Key: a.Field, Value: a.Expr.Build()})
	}
	return bson.D{{Key: "$group", Value: doc}}
}
