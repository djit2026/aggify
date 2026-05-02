package stage

import "go.mongodb.org/mongo-driver/bson"

// matchStage implements Stage for $match.
type matchStage struct{ filter bson.D }

// Match creates a $match stage from a filter document.
// The filter is typically produced by the q package.
//
//	stage.Match(q.Eq("status", "active"))
//	→ { $match: { status: "active" } }
func Match(filter bson.D) Stage {
	return matchStage{filter}
}

func (m matchStage) Build() bson.D {
	return bson.D{{Key: "$match", Value: m.filter}}
}
