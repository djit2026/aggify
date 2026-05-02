package stage

import (
	"github.com/djit2026/aggify/expr"
	"go.mongodb.org/mongo-driver/bson"
)

// BucketBoundary is the lower boundary value for a $bucket boundary.
type BucketBoundary = any

type bucketStage struct {
	groupBy    expr.Expr
	boundaries []BucketBoundary
	defaultID  any
	output     []Accumulator
}

// Bucket creates a $bucket stage.
//
//	stage.Bucket(
//	    expr.Field("price"),
//	    []any{0, 50, 100, 500},
//	    "Other",
//	    stage.Acc("count", expr.Sum(expr.Value(1))),
//	)
func Bucket(groupBy expr.Expr, boundaries []BucketBoundary, defaultID any, output ...Accumulator) Stage {
	return bucketStage{groupBy, boundaries, defaultID, output}
}

func (b bucketStage) Build() bson.D {
	bounds := make(bson.A, len(b.boundaries))
	for i, bnd := range b.boundaries {
		bounds[i] = bnd
	}
	doc := bson.D{
		{Key: "groupBy", Value: b.groupBy.Build()},
		{Key: "boundaries", Value: bounds},
	}
	if b.defaultID != nil {
		doc = append(doc, bson.E{Key: "default", Value: b.defaultID})
	}
	if len(b.output) > 0 {
		outDoc := make(bson.D, len(b.output))
		for i, a := range b.output {
			outDoc[i] = bson.E{Key: a.Field, Value: a.Expr.Build()}
		}
		doc = append(doc, bson.E{Key: "output", Value: outDoc})
	}
	return bson.D{{Key: "$bucket", Value: doc}}
}

// BucketAutoStage implements $bucketAuto.
type bucketAutoStage struct {
	groupBy     expr.Expr
	buckets     int
	granularity string
	output      []Accumulator
}

// BucketAuto creates a $bucketAuto stage.
func BucketAuto(groupBy expr.Expr, buckets int, output ...Accumulator) *bucketAutoStage {
	return &bucketAutoStage{groupBy: groupBy, buckets: buckets, output: output}
}

// Granularity sets the preferred number series for bucket boundaries.
func (b *bucketAutoStage) Granularity(g string) *bucketAutoStage {
	b.granularity = g
	return b
}

func (b *bucketAutoStage) Build() bson.D {
	doc := bson.D{
		{Key: "groupBy", Value: b.groupBy.Build()},
		{Key: "buckets", Value: b.buckets},
	}
	if b.granularity != "" {
		doc = append(doc, bson.E{Key: "granularity", Value: b.granularity})
	}
	if len(b.output) > 0 {
		outDoc := make(bson.D, len(b.output))
		for i, a := range b.output {
			outDoc[i] = bson.E{Key: a.Field, Value: a.Expr.Build()}
		}
		doc = append(doc, bson.E{Key: "output", Value: outDoc})
	}
	return bson.D{{Key: "$bucketAuto", Value: doc}}
}
