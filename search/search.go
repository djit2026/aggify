package search

import "go.mongodb.org/mongo-driver/bson"

// Operator represents an Atlas Search operator.
type Operator interface {
	Build() bson.D
}

type rawSearch struct { d bson.D }

func (r rawSearch) Build() bson.D { return r.d }

// Raw is an escape hatch for custom search operators.
func Raw(doc bson.D) Operator {
	return rawSearch{doc}
}

// Text builds a text search operator.
func Text(path, query string) Operator {
	return rawSearch{bson.D{{Key: "text", Value: bson.D{
		{Key: "path", Value: path},
		{Key: "query", Value: query},
	}}}}
}

// Autocomplete builds an autocomplete search operator.
func Autocomplete(path, query string) Operator {
	return rawSearch{bson.D{{Key: "autocomplete", Value: bson.D{
		{Key: "path", Value: path},
		{Key: "query", Value: query},
	}}}}
}

// CompoundBuilder builds a compound search operator.
type CompoundBuilder struct {
	must    []Operator
	should  []Operator
	mustNot []Operator
	filter  []Operator
}

// Compound returns a new CompoundBuilder.
func Compound() *CompoundBuilder {
	return &CompoundBuilder{}
}

// Must adds MUST clauses.
func (c *CompoundBuilder) Must(ops ...Operator) *CompoundBuilder {
	c.must = append(c.must, ops...)
	return c
}

// Should adds SHOULD clauses.
func (c *CompoundBuilder) Should(ops ...Operator) *CompoundBuilder {
	c.should = append(c.should, ops...)
	return c
}

// MustNot adds MUST_NOT clauses.
func (c *CompoundBuilder) MustNot(ops ...Operator) *CompoundBuilder {
	c.mustNot = append(c.mustNot, ops...)
	return c
}

// Filter adds FILTER clauses.
func (c *CompoundBuilder) Filter(ops ...Operator) *CompoundBuilder {
	c.filter = append(c.filter, ops...)
	return c
}

// Build compiles the compound operator to BSON.
func (c *CompoundBuilder) Build() bson.D {
	doc := bson.D{}

	buildArray := func(ops []Operator) bson.A {
		arr := make(bson.A, len(ops))
		for i, op := range ops {
			arr[i] = op.Build()
		}
		return arr
	}

	if len(c.must) > 0 {
		doc = append(doc, bson.E{Key: "must", Value: buildArray(c.must)})
	}
	if len(c.should) > 0 {
		doc = append(doc, bson.E{Key: "should", Value: buildArray(c.should)})
	}
	if len(c.mustNot) > 0 {
		doc = append(doc, bson.E{Key: "mustNot", Value: buildArray(c.mustNot)})
	}
	if len(c.filter) > 0 {
		doc = append(doc, bson.E{Key: "filter", Value: buildArray(c.filter)})
	}

	return bson.D{{Key: "compound", Value: doc}}
}
