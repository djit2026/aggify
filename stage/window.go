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

// --- Window-only expression helpers ---
// These are used exclusively as the Expr argument inside WindowFE().

// WinCovariancePop builds a $covariancePop window expression.
// Returns the population covariance of two numeric expressions over the window.
func WinCovariancePop(expr1, expr2 expr.Expr) expr.Expr {
	return expr.Raw(bson.D{{Key: "$covariancePop", Value: bson.A{expr1.Build(), expr2.Build()}}})
}

// WinCovarianceSamp builds a $covarianceSamp window expression.
// Returns the sample covariance of two numeric expressions over the window.
func WinCovarianceSamp(expr1, expr2 expr.Expr) expr.Expr {
	return expr.Raw(bson.D{{Key: "$covarianceSamp", Value: bson.A{expr1.Build(), expr2.Build()}}})
}

// WinDocumentNumber builds a $documentNumber window expression.
// Returns the position (1-indexed) of the current document within its partition.
func WinDocumentNumber() expr.Expr {
	return expr.Raw(bson.D{{Key: "$documentNumber", Value: bson.D{}}})
}

// WinRank builds a $rank window expression.
// Returns the rank of the current document within its partition (gaps for ties).
func WinRank() expr.Expr {
	return expr.Raw(bson.D{{Key: "$rank", Value: bson.D{}}})
}

// WinDenseRank builds a $denseRank window expression.
// Returns the rank without gaps for ties.
func WinDenseRank() expr.Expr {
	return expr.Raw(bson.D{{Key: "$denseRank", Value: bson.D{}}})
}

// WinShift builds a $shift window expression.
// Returns the value of output at a given offset from the current document.
// defaultValue is returned when the offset is out of range (optional, pass nil to omit).
func WinShift(output expr.Expr, by int, defaultValue expr.Expr) expr.Expr {
	doc := bson.D{
		{Key: "output", Value: output.Build()},
		{Key: "by", Value: by},
	}
	if defaultValue != nil {
		doc = append(doc, bson.E{Key: "default", Value: defaultValue.Build()})
	}
	return expr.Raw(bson.D{{Key: "$shift", Value: doc}})
}

// WinDerivative builds a $derivative window expression.
// Divides the change in output by the change in the sort field over the window.
// unit is optional (e.g. "second", "day") — only valid when the sort field is a date.
func WinDerivative(output expr.Expr, unit ...string) expr.Expr {
	doc := bson.D{{Key: "output", Value: output.Build()}}
	if len(unit) > 0 && unit[0] != "" {
		doc = append(doc, bson.E{Key: "unit", Value: unit[0]})
	}
	return expr.Raw(bson.D{{Key: "$derivative", Value: doc}})
}

// WinIntegral builds a $integral window expression.
// Integrates the output over the sorted window.
// unit is optional (e.g. "second", "day") — only valid when the sort field is a date.
func WinIntegral(output expr.Expr, unit ...string) expr.Expr {
	doc := bson.D{{Key: "output", Value: output.Build()}}
	if len(unit) > 0 && unit[0] != "" {
		doc = append(doc, bson.E{Key: "unit", Value: unit[0]})
	}
	return expr.Raw(bson.D{{Key: "$integral", Value: doc}})
}

// WinExpMovingAvg builds a $expMovingAvg window expression.
// Computes an exponential moving average of output.
// Supply either n (number of historical documents) OR alpha (decay factor 0–1), not both.
func WinExpMovingAvg(output expr.Expr, n int, alpha float64) expr.Expr {
	doc := bson.D{{Key: "input", Value: output.Build()}}
	if n > 0 {
		doc = append(doc, bson.E{Key: "N", Value: n})
	} else {
		doc = append(doc, bson.E{Key: "alpha", Value: alpha})
	}
	return expr.Raw(bson.D{{Key: "$expMovingAvg", Value: doc}})
}
