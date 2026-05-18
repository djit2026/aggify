package q

import (
	"github.com/djit2026/aggify/geo"
	"go.mongodb.org/mongo-driver/bson"
)

// NearOptions holds optional distance bounds for $near and $nearSphere.
type NearOptions struct {
	// MaxDistance is the maximum distance in meters (GeoJSON) or radians (legacy).
	MaxDistance *float64
	// MinDistance is the minimum distance in meters (GeoJSON) or radians (legacy).
	MinDistance *float64
}

// GeoWithin builds a $geoWithin filter using a GeoJSON geometry.
// Use this for querying points, polygons, etc. stored as GeoJSON.
//
//	q.GeoWithin("location", geo.Polygon(...))
func GeoWithin(field string, geom geo.Geometry) bson.D {
	return bson.D{{Key: field, Value: bson.D{
		{Key: "$geoWithin", Value: bson.D{{Key: "$geometry", Value: geom.Build()}}},
	}}}
}

// GeoWithinShape builds a $geoWithin filter using a legacy shape (Box, Circle, etc.).
// Legacy shapes use flat (planar) geometry. For spherical use CenterSphere.
//
//	q.GeoWithinShape("location", geo.Box(geo.LngLat(-74, 40), geo.LngLat(-73, 41)))
func GeoWithinShape(field string, shape geo.Shape) bson.D {
	return bson.D{{Key: field, Value: bson.D{
		{Key: "$geoWithin", Value: shape.Build()},
	}}}
}

// GeoIntersects builds a $geoIntersects filter.
// Returns documents whose geo field intersects the given GeoJSON geometry.
//
//	q.GeoIntersects("area", geo.Polygon(...))
func GeoIntersects(field string, geom geo.Geometry) bson.D {
	return bson.D{{Key: field, Value: bson.D{
		{Key: "$geoIntersects", Value: bson.D{{Key: "$geometry", Value: geom.Build()}}},
	}}}
}

// Near builds a $near filter using a GeoJSON Point.
// Returns documents sorted by distance (closest first).
// Requires a 2dsphere index on the field.
//
//	q.Near("location", geo.Point(-73.98, 40.75), q.NearOptions{MaxDistance: ptr(5000.0)})
func Near(field string, point geo.Geometry, opts ...NearOptions) bson.D {
	inner := bson.D{{Key: "$geometry", Value: point.Build()}}
	if len(opts) > 0 {
		if opts[0].MaxDistance != nil {
			inner = append(inner, bson.E{Key: "$maxDistance", Value: *opts[0].MaxDistance})
		}
		if opts[0].MinDistance != nil {
			inner = append(inner, bson.E{Key: "$minDistance", Value: *opts[0].MinDistance})
		}
	}
	return bson.D{{Key: field, Value: bson.D{{Key: "$near", Value: inner}}}}
}

// NearSphere builds a $nearSphere filter using a GeoJSON Point.
// Like $near but always uses spherical geometry, even for legacy coordinate pairs.
//
//	q.NearSphere("location", geo.Point(-73.98, 40.75))
func NearSphere(field string, point geo.Geometry, opts ...NearOptions) bson.D {
	inner := bson.D{{Key: "$geometry", Value: point.Build()}}
	if len(opts) > 0 {
		if opts[0].MaxDistance != nil {
			inner = append(inner, bson.E{Key: "$maxDistance", Value: *opts[0].MaxDistance})
		}
		if opts[0].MinDistance != nil {
			inner = append(inner, bson.E{Key: "$minDistance", Value: *opts[0].MinDistance})
		}
	}
	return bson.D{{Key: field, Value: bson.D{{Key: "$nearSphere", Value: inner}}}}
}
