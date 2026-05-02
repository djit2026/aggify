package stage

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FacetPipeline associates a name with a sub-pipeline for $facet.
type FacetPipeline struct {
	Name     string
	Pipeline mongo.Pipeline
}

type facetStage struct{ facets []FacetPipeline }

// Facet creates a $facet stage.
// Each FacetPipeline runs independently on the input documents.
//
//	stage.Facet(
//	    stage.FacetPipeline{"byStatus", agg.New().Group(...).Build()},
//	    stage.FacetPipeline{"byStore",  agg.New().Group(...).Build()},
//	)
func Facet(facets ...FacetPipeline) Stage {
	return facetStage{facets}
}

func (f facetStage) Build() bson.D {
	doc := make(bson.D, len(f.facets))
	for i, fp := range f.facets {
		mustNotEmpty(fp.Name, "facet pipeline name")
		doc[i] = bson.E{Key: fp.Name, Value: fp.Pipeline}
	}
	return bson.D{{Key: "$facet", Value: doc}}
}
