package stage

import (
	"github.com/djit2026/aggify/expr"
	"go.mongodb.org/mongo-driver/bson"
)

type replaceRootStage struct{ newRoot expr.Expr }

// ReplaceRoot creates a $replaceRoot stage.
//
//	stage.ReplaceRoot(expr.MergeObjects(expr.Field("details"), expr.Root))
func ReplaceRoot(newRoot expr.Expr) Stage {
	return replaceRootStage{newRoot}
}

// ReplaceWith creates a $replaceWith stage (alias for $replaceRoot in MongoDB 4.2+).
func ReplaceWith(newRoot expr.Expr) Stage {
	return replaceWithStage{newRoot}
}

func (r replaceRootStage) Build() bson.D {
	return bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: r.newRoot.Build()}}}}
}

type replaceWithStage struct{ newRoot expr.Expr }

func (r replaceWithStage) Build() bson.D {
	return bson.D{{Key: "$replaceWith", Value: r.newRoot.Build()}}
}
