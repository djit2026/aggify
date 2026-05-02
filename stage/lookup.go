package stage

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// --- Simple $lookup ---

// lookupStage implements Stage for simple $lookup.
type lookupStage struct {
	from         string
	localField   string
	foreignField string
	as           string
}

// Lookup creates a simple $lookup stage (equality join).
//
//	stage.Lookup("users", "userId", "_id", "user")
//	→ { $lookup: { from: "users", localField: "userId", foreignField: "_id", as: "user" } }
func Lookup(from, localField, foreignField, as string) Stage {
	mustNotEmpty(from, "lookup from")
	mustNotEmpty(as, "lookup as")
	return lookupStage{from, localField, foreignField, as}
}

func (l lookupStage) Build() bson.D {
	return bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: l.from},
		{Key: "localField", Value: l.localField},
		{Key: "foreignField", Value: l.foreignField},
		{Key: "as", Value: l.as},
	}}}
}

// --- Pipeline $lookup (with $let) ---

// LookupPipelineBuilder builds a pipeline-form $lookup.
type LookupPipelineBuilder struct {
	from     string
	as       string
	letVars  bson.D
	pipeline mongo.Pipeline
}

// LookupPipeline starts a pipeline-form $lookup builder.
//
//	stage.LookupPipeline("assignments", "assignments").
//	    Let(bson.D{{"orderId", "$_id"}}).
//	    Pipeline(agg.New().Match(q.Expr(...)).Build())
func LookupPipeline(from, as string) *LookupPipelineBuilder {
	mustNotEmpty(from, "lookup from")
	mustNotEmpty(as, "lookup as")
	return &LookupPipelineBuilder{from: from, as: as}
}

// Let sets variable bindings for use inside the sub-pipeline.
func (b *LookupPipelineBuilder) Let(vars bson.D) *LookupPipelineBuilder {
	b.letVars = vars
	return b
}

// Pipeline sets the sub-pipeline.
func (b *LookupPipelineBuilder) Pipeline(p mongo.Pipeline) *LookupPipelineBuilder {
	b.pipeline = p
	return b
}

// Build implements Stage.
func (b *LookupPipelineBuilder) Build() bson.D {
	doc := bson.D{
		{Key: "from", Value: b.from},
	}
	if len(b.letVars) > 0 {
		doc = append(doc, bson.E{Key: "let", Value: b.letVars})
	}
	doc = append(doc, bson.E{Key: "pipeline", Value: b.pipeline})
	doc = append(doc, bson.E{Key: "as", Value: b.as})
	return bson.D{{Key: "$lookup", Value: doc}}
}
