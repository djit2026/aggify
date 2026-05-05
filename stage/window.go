package stage

import (
	"github.com/djit2026/aggify/expr"
	"go.mongodb.org/mongo-driver/bson"
)

// SortWindow specifies sorting for the window.
// It maps directly to bson.D to maintain insertion order for multiple sort fields.
type SortWindow bson.D

// WindowOutput defines an output field for $setWindowFields.
type WindowOutput struct {
	Field       string
	Expr        expr.Expr
	WindowBound bson.D
}

// WindowFE creates a WindowOutput.
// bounds can be nil if not needed (e.g., for ranking functions).
func WindowFE(field string, e expr.Expr, bounds bson.D) WindowOutput {
	mustNotEmpty(field, "window output field")
	return WindowOutput{Field: field, Expr: e, WindowBound: bounds}
}

// WindowBounds creates a bounds document, e.g. "documents", ["unbounded", "current"]
func WindowBounds(unit string, lower, upper any) bson.D {
	return bson.D{{Key: unit, Value: bson.A{lower, upper}}}
}

type setWindowFieldsStage struct {
	partitionBy expr.Expr
	sortBy      SortWindow
	outputs     []WindowOutput
}

// SetWindowFields creates a $setWindowFields stage.
func SetWindowFields(partitionBy expr.Expr, sortBy SortWindow, outputs ...WindowOutput) Stage {
	return setWindowFieldsStage{
		partitionBy: partitionBy,
		sortBy:      sortBy,
		outputs:     outputs,
	}
}

func (s setWindowFieldsStage) Build() bson.D {
	doc := bson.D{}
	if s.partitionBy != nil {
		doc = append(doc, bson.E{Key: "partitionBy", Value: s.partitionBy.Build()})
	}
	if len(s.sortBy) > 0 {
		doc = append(doc, bson.E{Key: "sortBy", Value: s.sortBy})
	}

	outDoc := bson.D{}
	for _, o := range s.outputs {
		var fieldDoc bson.D
		if e := o.Expr; e != nil {
			val := o.Expr.Build()
			if vDoc, ok := val.(bson.D); ok && len(vDoc) == 1 {
				opDoc := vDoc
				if len(o.WindowBound) > 0 {
					opDoc = append(opDoc, bson.E{Key: "window", Value: o.WindowBound})
				}
				fieldDoc = opDoc
			} else {
				if vDoc, ok := val.(bson.D); ok {
					fieldDoc = vDoc
				} else {
					panic("mono-query/stage: window expression must build to bson.D")
				}
			}
		}
		outDoc = append(outDoc, bson.E{Key: o.Field, Value: fieldDoc})
	}
	
	doc = append(doc, bson.E{Key: "output", Value: outDoc})
	return bson.D{{Key: "$setWindowFields", Value: doc}}
}
