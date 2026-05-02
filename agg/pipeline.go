// Package agg provides the top-level fluent API for building MongoDB
// aggregation pipelines.
//
// Typical usage:
//
//	pipeline := agg.New().
//	    Match(q.Eq("userId", userId)).
//	    Stage(filterActiveItems()).   // reusable business-logic stage
//	    Stage(groupByStore()).
//	    Sort(stage.SortField{"total", stage.Desc}).
//	    Limit(10).
//	    Build()
//
//	cursor, err := coll.Aggregate(ctx, pipeline)
package agg

import (
	"encoding/json"

	"github.com/djit2026/aggify/expr"
	"github.com/djit2026/aggify/stage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Pipeline is a fluent aggregation pipeline builder.
// Methods return the same *Pipeline for chaining.
// Call Build() to produce the final mongo.Pipeline.
type Pipeline struct {
	stages []stage.Stage
}

// New creates an empty Pipeline.
func New() *Pipeline {
	return &Pipeline{}
}

// Stage appends any Stage to the pipeline.
// This is the primary composition point: reusable stage functions return
// stage.Stage and slot in cleanly here.
//
//	agg.New().Stage(myCustomStage())
func (p *Pipeline) Stage(s stage.Stage) *Pipeline {
	p.stages = append(p.stages, s)
	return p
}

// Match appends a $match stage.
func (p *Pipeline) Match(filter bson.D) *Pipeline {
	return p.Stage(stage.Match(filter))
}

// Group appends a $group stage.
func (p *Pipeline) Group(id expr.Expr, accs ...stage.Accumulator) *Pipeline {
	return p.Stage(stage.Group(id, accs...))
}

// Project appends a $project stage.
// Use the returned projectStage builder, then call its methods and pass it
// back via .Stage() — or use the convenience shortcut:
//
//	agg.New().Stage(stage.Project().Include("name").Exclude("_id"))
func (p *Pipeline) Project(s *stage.ProjectStage) *Pipeline {
	return p.Stage(s)
}

// Lookup appends a simple $lookup stage.
func (p *Pipeline) Lookup(from, localField, foreignField, as string) *Pipeline {
	return p.Stage(stage.Lookup(from, localField, foreignField, as))
}

// Unwind appends an $unwind stage (simple string form).
// For options (preserveNull, includeArrayIndex), use .Stage(stage.Unwind(...).PreserveNullAndEmpty(true)).
func (p *Pipeline) Unwind(path string) *Pipeline {
	return p.Stage(stage.Unwind(path))
}

// Sort appends a $sort stage.
func (p *Pipeline) Sort(fields ...stage.SortField) *Pipeline {
	return p.Stage(stage.Sort(fields...))
}

// SortAsc appends a single ascending $sort stage.
func (p *Pipeline) SortAsc(field string) *Pipeline {
	return p.Stage(stage.SortAsc(field))
}

// SortDesc appends a single descending $sort stage.
func (p *Pipeline) SortDesc(field string) *Pipeline {
	return p.Stage(stage.SortDesc(field))
}

// Limit appends a $limit stage.
func (p *Pipeline) Limit(n int64) *Pipeline {
	return p.Stage(stage.Limit(n))
}

// Skip appends a $skip stage.
func (p *Pipeline) Skip(n int64) *Pipeline {
	return p.Stage(stage.Skip(n))
}

// AddFields appends an $addFields stage.
func (p *Pipeline) AddFields(fields ...stage.FieldExpr) *Pipeline {
	return p.Stage(stage.AddFields(fields...))
}

// Set appends a $set stage (alias for $addFields).
func (p *Pipeline) Set(fields ...stage.FieldExpr) *Pipeline {
	return p.Stage(stage.Set(fields...))
}

// Unset appends an $unset stage.
func (p *Pipeline) Unset(fields ...string) *Pipeline {
	return p.Stage(stage.Unset(fields...))
}

// ReplaceRoot appends a $replaceRoot stage.
func (p *Pipeline) ReplaceRoot(newRoot expr.Expr) *Pipeline {
	return p.Stage(stage.ReplaceRoot(newRoot))
}

// ReplaceWith appends a $replaceWith stage.
func (p *Pipeline) ReplaceWith(newRoot expr.Expr) *Pipeline {
	return p.Stage(stage.ReplaceWith(newRoot))
}

// Count appends a $count stage.
func (p *Pipeline) Count(field string) *Pipeline {
	return p.Stage(stage.Count(field))
}

// Raw appends a raw bson.D as a stage. Use as an escape hatch.
func (p *Pipeline) Raw(d bson.D) *Pipeline {
	return p.Stage(stage.Raw(d))
}

// Len returns the current number of stages.
func (p *Pipeline) Len() int {
	return len(p.stages)
}

// Build compiles the pipeline into a mongo.Pipeline ready to pass to
// collection.Aggregate().
func (p *Pipeline) Build() mongo.Pipeline {
	out := make(mongo.Pipeline, len(p.stages))
	for i, s := range p.stages {
		out[i] = s.Build()
	}
	return out
}

// MustJSON returns the pipeline serialised as indented JSON.
// Useful for debugging and logging — do not use in production hot paths.
func (p *Pipeline) MustJSON() string {
	b, err := json.MarshalIndent(p.Build(), "", "  ")
	if err != nil {
		panic("mono-query/agg: MustJSON: " + err.Error())
	}
	return string(b)
}
