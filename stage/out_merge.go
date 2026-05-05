package stage

import (
	"go.mongodb.org/mongo-driver/bson"
)

// Out creates a $out stage to write to a collection in the same database.
func Out(collection string) Stage {
	mustNotEmpty(collection, "out collection")
	return rawStage{bson.D{{Key: "$out", Value: collection}}}
}

// OutToDB creates a $out stage to write to a collection in a specific database.
func OutToDB(db, collection string) Stage {
	mustNotEmpty(db, "out database")
	mustNotEmpty(collection, "out collection")
	return rawStage{bson.D{{Key: "$out", Value: bson.D{
		{Key: "db", Value: db},
		{Key: "coll", Value: collection},
	}}}}
}

// MergeOptions defines options for the $merge stage.
type MergeOptions struct {
	IntoDB         string
	IntoCollection string
	On             []string
	Let            bson.D
	WhenMatched    string // e.g. "replace", "keepExisting", "merge", "fail", "pipeline"
	WhenNotMatched string // e.g. "insert", "discard", "fail"
}

// Merge creates a $merge stage.
func Merge(opts MergeOptions) Stage {
	mustNotEmpty(opts.IntoCollection, "merge into collection")

	var into any = opts.IntoCollection
	if opts.IntoDB != "" {
		into = bson.D{
			{Key: "db", Value: opts.IntoDB},
			{Key: "coll", Value: opts.IntoCollection},
		}
	}

	doc := bson.D{{Key: "into", Value: into}}

	if len(opts.On) > 0 {
		if len(opts.On) == 1 {
			doc = append(doc, bson.E{Key: "on", Value: opts.On[0]})
		} else {
			doc = append(doc, bson.E{Key: "on", Value: opts.On})
		}
	}
	if len(opts.Let) > 0 {
		doc = append(doc, bson.E{Key: "let", Value: opts.Let})
	}
	if opts.WhenMatched != "" {
		doc = append(doc, bson.E{Key: "whenMatched", Value: opts.WhenMatched})
	}
	if opts.WhenNotMatched != "" {
		doc = append(doc, bson.E{Key: "whenNotMatched", Value: opts.WhenNotMatched})
	}

	return rawStage{bson.D{{Key: "$merge", Value: doc}}}
}
