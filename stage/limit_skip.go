package stage

import "go.mongodb.org/mongo-driver/bson"

type limitStage struct{ n int64 }

// Limit creates a $limit stage.
func Limit(n int64) Stage { return limitStage{n} }

func (l limitStage) Build() bson.D {
	return bson.D{{Key: "$limit", Value: l.n}}
}

type skipStage struct{ n int64 }

// Skip creates a $skip stage.
func Skip(n int64) Stage { return skipStage{n} }

func (s skipStage) Build() bson.D {
	return bson.D{{Key: "$skip", Value: s.n}}
}
