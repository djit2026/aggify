package stage

import (
	"github.com/djit2026/aggify/geo"
	"go.mongodb.org/mongo-driver/bson"
)

// GeoNearOptions holds all options for the $geoNear stage.
// Near and DistanceField are required; all others are optional.
type GeoNearOptions struct {
	// Near is the GeoJSON Point to measure distance from. Required.
	Near geo.Geometry

	// DistanceField is the output field that will contain the distance.
	// Use dot notation for nested fields, e.g. "dist.calculated". Required.
	DistanceField string

	// Spherical uses spherical geometry to calculate distances.
	// Must be true for GeoJSON Points and 2dsphere indexes.
	Spherical bool

	// MaxDistance is the maximum distance from the center point in meters.
	MaxDistance *float64

	// MinDistance is the minimum distance from the center point in meters.
	MinDistance *float64

	// Query is an optional filter applied before the geo sort.
	// Equivalent to a $match stage before $geoNear, but more efficient.
	Query bson.D

	// IncludeLocs specifies an output field to store the matched location coordinate.
	IncludeLocs string

	// DistanceMultiplier multiplies the distance before storing it in DistanceField.
	// Useful for unit conversion (e.g. 0.001 to convert meters to km).
	DistanceMultiplier *float64

	// Key specifies which geo index to use when a collection has multiple.
	Key string
}

type geoNearStage struct {
	opts GeoNearOptions
}

// GeoNear creates a $geoNear stage.
// $geoNear must be the FIRST stage in the pipeline and requires a geo index.
//
//	// Find restaurants within 2 km of Times Square, sorted by distance
//	agg.New().
//	    GeoNear(stage.GeoNearOptions{
//	        Near:          geo.Point(-73.9857, 40.7580),
//	        DistanceField: "dist.meters",
//	        MaxDistance:   ptr(2000.0),
//	        Spherical:     true,
//	        Query:         q.Eq("category", "restaurant"),
//	    }).
//	    Limit(10).
//	    Build()
func GeoNear(opts GeoNearOptions) Stage {
	mustNotEmpty(opts.DistanceField, "geoNear distanceField")
	return geoNearStage{opts: opts}
}

func (g geoNearStage) Build() bson.D {
	o := g.opts
	doc := bson.D{
		{Key: "near", Value: o.Near.Build()},
		{Key: "distanceField", Value: o.DistanceField},
		{Key: "spherical", Value: o.Spherical},
	}
	if o.MaxDistance != nil {
		doc = append(doc, bson.E{Key: "maxDistance", Value: *o.MaxDistance})
	}
	if o.MinDistance != nil {
		doc = append(doc, bson.E{Key: "minDistance", Value: *o.MinDistance})
	}
	if len(o.Query) > 0 {
		doc = append(doc, bson.E{Key: "query", Value: o.Query})
	}
	if o.IncludeLocs != "" {
		doc = append(doc, bson.E{Key: "includeLocs", Value: o.IncludeLocs})
	}
	if o.DistanceMultiplier != nil {
		doc = append(doc, bson.E{Key: "distanceMultiplier", Value: *o.DistanceMultiplier})
	}
	if o.Key != "" {
		doc = append(doc, bson.E{Key: "key", Value: o.Key})
	}
	return bson.D{{Key: "$geoNear", Value: doc}}
}
