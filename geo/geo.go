// Package geo provides GeoJSON geometry constructors and legacy MongoDB shape
// helpers for use with geo query operators (q.Near, q.GeoWithin, etc.)
// and the $geoNear aggregation stage.
//
// MongoDB uses [longitude, latitude] coordinate order for GeoJSON, which is
// the opposite of the common [lat, lng] convention. All constructors in this
// package follow the GeoJSON / MongoDB [longitude, latitude] standard.
//
// Quick start:
//
//	// Find restaurants within 2 km of Times Square
//	agg.New().
//	    GeoNear(stage.GeoNearOptions{
//	        Near:          geo.Point(-73.9857, 40.7580),
//	        DistanceField: "dist.meters",
//	        MaxDistance:   ptr(2000.0),
//	        Spherical:     true,
//	    }).
//	    Limit(20).
//	    Build()
package geo

import "go.mongodb.org/mongo-driver/bson"

// Coord represents a [longitude, latitude] coordinate pair.
// Always longitude first, latitude second — this is the GeoJSON standard.
type Coord [2]float64

// LngLat is a convenience constructor for Coord that makes axis order explicit.
//
//	geo.LngLat(-73.98, 40.75)  // Times Square, NYC
func LngLat(lng, lat float64) Coord { return Coord{lng, lat} }

// Geometry is a GeoJSON geometry document ready for embedding in BSON queries.
type Geometry bson.D

// Build returns the underlying bson.D representation.
func (g Geometry) Build() bson.D { return bson.D(g) }

// Shape is a legacy MongoDB shape for $geoWithin queries (non-GeoJSON).
// Use geo.Box, geo.Circle, geo.CenterSphere, or geo.LegacyPolygon.
type Shape bson.D

// Build returns the underlying bson.D representation.
func (s Shape) Build() bson.D { return bson.D(s) }

// ─── GeoJSON Geometries ────────────────────────────────────────────────────

// Point builds a GeoJSON Point geometry.
//
//	geo.Point(-73.9857, 40.7580)   // Times Square, NYC
func Point(lng, lat float64) Geometry {
	return Geometry(bson.D{
		{Key: "type", Value: "Point"},
		{Key: "coordinates", Value: bson.A{lng, lat}},
	})
}

// LineString builds a GeoJSON LineString geometry.
//
//	geo.LineString(geo.LngLat(-73.98, 40.75), geo.LngLat(-74.00, 40.71))
func LineString(coords ...Coord) Geometry {
	return Geometry(bson.D{
		{Key: "type", Value: "LineString"},
		{Key: "coordinates", Value: coordsToArray(coords)},
	})
}

// Polygon builds a GeoJSON Polygon geometry.
// The first ring is the outer boundary; subsequent rings are holes.
// Each ring must be closed (first and last coordinate must match).
//
//	geo.Polygon([]geo.Coord{
//	    geo.LngLat(-73.99, 40.76),
//	    geo.LngLat(-73.98, 40.76),
//	    geo.LngLat(-73.98, 40.75),
//	    geo.LngLat(-73.99, 40.75),
//	    geo.LngLat(-73.99, 40.76), // closed
//	})
func Polygon(rings ...[]Coord) Geometry {
	outer := make(bson.A, len(rings))
	for i, ring := range rings {
		outer[i] = coordsToArray(ring)
	}
	return Geometry(bson.D{
		{Key: "type", Value: "Polygon"},
		{Key: "coordinates", Value: outer},
	})
}

// MultiPoint builds a GeoJSON MultiPoint geometry.
func MultiPoint(coords ...Coord) Geometry {
	return Geometry(bson.D{
		{Key: "type", Value: "MultiPoint"},
		{Key: "coordinates", Value: coordsToArray(coords)},
	})
}

// MultiLineString builds a GeoJSON MultiLineString geometry.
func MultiLineString(lines ...[]Coord) Geometry {
	arr := make(bson.A, len(lines))
	for i, line := range lines {
		arr[i] = coordsToArray(line)
	}
	return Geometry(bson.D{
		{Key: "type", Value: "MultiLineString"},
		{Key: "coordinates", Value: arr},
	})
}

// MultiPolygon builds a GeoJSON MultiPolygon geometry.
func MultiPolygon(polygons ...[][]Coord) Geometry {
	arr := make(bson.A, len(polygons))
	for i, poly := range polygons {
		rings := make(bson.A, len(poly))
		for j, ring := range poly {
			rings[j] = coordsToArray(ring)
		}
		arr[i] = rings
	}
	return Geometry(bson.D{
		{Key: "type", Value: "MultiPolygon"},
		{Key: "coordinates", Value: arr},
	})
}

// GeometryCollection builds a GeoJSON GeometryCollection.
func GeometryCollection(geoms ...Geometry) Geometry {
	arr := make(bson.A, len(geoms))
	for i, g := range geoms {
		arr[i] = g.Build()
	}
	return Geometry(bson.D{
		{Key: "type", Value: "GeometryCollection"},
		{Key: "geometries", Value: arr},
	})
}

// ─── Legacy Shapes (for $geoWithin) ─────────────────────────────────────────

// Box builds a legacy $box shape using planar (flat-earth) coordinates.
// The box is defined by its bottom-left and top-right corners.
//
//	geo.Box(geo.LngLat(-74.00, 40.71), geo.LngLat(-73.96, 40.75))
func Box(bottomLeft, topRight Coord) Shape {
	return Shape(bson.D{{Key: "$box", Value: bson.A{
		bson.A{bottomLeft[0], bottomLeft[1]},
		bson.A{topRight[0], topRight[1]},
	}}})
}

// Circle builds a legacy $center shape for flat-earth distance queries.
// radius is in the same units as the coordinate system (not meters/km).
// For spherical distance use CenterSphere.
func Circle(center Coord, radius float64) Shape {
	return Shape(bson.D{{Key: "$center", Value: bson.A{
		bson.A{center[0], center[1]}, radius,
	}}})
}

// CenterSphere builds a $centerSphere shape using spherical geometry.
// radius must be in radians. To convert from meters: meters / 6378100.0
//
//	// 5 km radius
//	geo.CenterSphere(geo.LngLat(-73.98, 40.75), 5000.0/6378100.0)
func CenterSphere(center Coord, radiusRadians float64) Shape {
	return Shape(bson.D{{Key: "$centerSphere", Value: bson.A{
		bson.A{center[0], center[1]}, radiusRadians,
	}}})
}

// LegacyPolygon builds a $polygon shape using planar coordinates.
// Points do not need to be closed (MongoDB closes it automatically).
func LegacyPolygon(coords ...Coord) Shape {
	return Shape(bson.D{{Key: "$polygon", Value: coordsToArray(coords)}})
}

// ─── internal helpers ─────────────────────────────────────────────────────

func coordsToArray(coords []Coord) bson.A {
	arr := make(bson.A, len(coords))
	for i, c := range coords {
		arr[i] = bson.A{c[0], c[1]}
	}
	return arr
}
