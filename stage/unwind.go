package stage

import "go.mongodb.org/mongo-driver/bson"

// unwindStage implements Stage for $unwind.
type unwindStage struct {
	path                       string
	includeArrayIndex          string
	preserveNullAndEmptyArrays bool
}

// Unwind starts an $unwind stage builder.
//
//	stage.Unwind("$items")
//	stage.Unwind("$items").PreserveNullAndEmpty(true)
func Unwind(path string) *unwindStage {
	mustNotEmpty(path, "unwind path")
	return &unwindStage{path: path}
}

// IncludeArrayIndex sets the field name for the array index.
func (u *unwindStage) IncludeArrayIndex(field string) *unwindStage {
	u.includeArrayIndex = field
	return u
}

// PreserveNullAndEmpty controls whether documents with missing/null/empty
// arrays are passed through the pipeline.
func (u *unwindStage) PreserveNullAndEmpty(v bool) *unwindStage {
	u.preserveNullAndEmptyArrays = v
	return u
}

// Build implements Stage.
func (u *unwindStage) Build() bson.D {
	if u.includeArrayIndex == "" && !u.preserveNullAndEmptyArrays {
		// Simple form — just a string.
		return bson.D{{Key: "$unwind", Value: u.path}}
	}
	doc := bson.D{{Key: "path", Value: u.path}}
	if u.includeArrayIndex != "" {
		doc = append(doc, bson.E{Key: "includeArrayIndex", Value: u.includeArrayIndex})
	}
	if u.preserveNullAndEmptyArrays {
		doc = append(doc, bson.E{Key: "preserveNullAndEmptyArrays", Value: true})
	}
	return bson.D{{Key: "$unwind", Value: doc}}
}
