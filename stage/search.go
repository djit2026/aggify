package stage

import (
	"github.com/djit2026/aggify/search"
	"go.mongodb.org/mongo-driver/bson"
)

// Search creates a $search stage for Atlas Search.
func Search(op search.Operator) Stage {
	return rawStage{bson.D{{Key: "$search", Value: op.Build()}}}
}

// SearchMeta creates a $searchMeta stage for Atlas Search metadata.
func SearchMeta(op search.Operator) Stage {
	return rawStage{bson.D{{Key: "$searchMeta", Value: op.Build()}}}
}
