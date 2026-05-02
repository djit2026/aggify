package stage

import (
	"github.com/djit2026/aggify/expr"
	"go.mongodb.org/mongo-driver/bson"
)

// ProjectField represents a single field specification in $project.
type ProjectField struct {
	key   string
	value any // 1, 0, or an expr.Expr.Build() result
}

// ProjectStage implements Stage for $project.
// Construct one with Project() and chain Include/Exclude/Computed calls.
type ProjectStage struct {
	fields []ProjectField
}

// Project starts a $project stage builder.
//
//	stage.Project().
//	    Include("name", "email").
//	    Exclude("_id").
//	    Computed("fullName", expr.Concat(expr.Field("first"), expr.Value(" "), expr.Field("last")))
func Project() *ProjectStage {
	return &ProjectStage{}
}

// Include adds fields to include (value = 1).
func (p *ProjectStage) Include(fields ...string) *ProjectStage {
	for _, f := range fields {
		mustNotEmpty(f, "project include field")
		p.fields = append(p.fields, ProjectField{key: f, value: 1})
	}
	return p
}

// Exclude adds fields to exclude (value = 0).
func (p *ProjectStage) Exclude(fields ...string) *ProjectStage {
	for _, f := range fields {
		mustNotEmpty(f, "project exclude field")
		p.fields = append(p.fields, ProjectField{key: f, value: 0})
	}
	return p
}

// Computed adds a computed field using an aggregation expression.
func (p *ProjectStage) Computed(field string, e expr.Expr) *ProjectStage {
	mustNotEmpty(field, "project computed field")
	p.fields = append(p.fields, ProjectField{key: field, value: e.Build()})
	return p
}

// Build implements Stage.
func (p *ProjectStage) Build() bson.D {
	doc := make(bson.D, len(p.fields))
	for i, f := range p.fields {
		doc[i] = bson.E{Key: f.key, Value: f.value}
	}
	return bson.D{{Key: "$project", Value: doc}}
}
