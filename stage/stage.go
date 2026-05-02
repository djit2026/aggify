// Package stage provides builders for every MongoDB aggregation stage.
// Each builder implements the Stage interface, which means stages can be
// composed freely and passed to agg.Pipeline via .Stage().
//
// All builders panic on structurally invalid input (e.g. empty field names).
// This surfaces API misuse at startup rather than at query time.
package stage

import "go.mongodb.org/mongo-driver/bson"

// Stage is implemented by every aggregation stage builder.
// Build() returns the stage as a bson.D ready for inclusion in a
// mongo.Pipeline.
type Stage interface {
	Build() bson.D
}

// Raw wraps an already-built bson.D as a Stage.
// Use this as an escape hatch when a stage is not yet supported by the library.
//
//	stage.Raw(bson.D{{"$sample", bson.D{{"size", 5}}}})
func Raw(d bson.D) Stage { return rawStage{d} }

type rawStage struct{ d bson.D }

func (r rawStage) Build() bson.D { return r.d }

// mustNotEmpty panics if s is an empty string, providing a clear error message.
func mustNotEmpty(s, label string) {
	if s == "" {
		panic("mono-query/stage: " + label + " must not be empty")
	}
}
