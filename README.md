# aggify

> **A production-grade, composable MongoDB aggregation pipeline builder for Go.**

[![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

---

## Why?

Working with aggregation pipelines in the official MongoDB Go driver means writing deeply-nested `bson.M` and `bson.D` literals:

```go
// ❌ Before — raw BSON, hard to read and maintain
bson.D{{"$group", bson.D{
    {"_id", bson.D{{"storeId", "$items.storeId"}}},
    {"items", bson.D{{"$push", bson.D{{"$cond", bson.A{
        bson.D{{"$eq", bson.A{"$$item.status", "active"}}},
        "$$item", "$$REMOVE",
    }}}}}},
}}}
```

`aggify` replaces that with a composable, readable DSL:

```go
// ✅ After
stage.Group(
    expr.Raw(bson.D{{Key: "storeId", Value: "$items.storeId"}}),
    stage.Acc("items", expr.Push(
        expr.Cond(
            expr.Eq(expr.Var("item.status"), expr.Value("active")),
            expr.Var("item"),
            expr.Remove,
        ),
    )),
)
```

---

## Features

- **Zero reflection** — direct compilation to `bson.D` / `mongo.Pipeline`
- **Zero overhead** — output is identical to hand-written BSON
- **Composable stages** — any `func() stage.Stage` is a reusable pipeline unit
- **Full expression engine** — comparison, logical, conditional, array, object, arithmetic, string, type, date, set, trigonometry, and dynamic field operators
- **Complete stage coverage** — `$match`, `$group`, `$project`, `$lookup`, `$unwind`, `$sort`, `$limit`, `$skip`, `$addFields`, `$set`, `$unset`, `$replaceRoot`, `$replaceWith`, `$count`, `$facet`, `$bucket`, `$bucketAuto`, `$setWindowFields`, `$geoNear`, `$out`, `$merge`
- **Geo queries** — full GeoJSON geometry builders + `$geoNear` stage + `$near`, `$geoWithin`, `$geoIntersects` query filters
- **Atlas Search** — dedicated DSL for `$search` and `$searchMeta`
- **Type-safe schema generation** — optional `aggify-gen` CLI tool to generate typo-free nested BSON path constants from your structs
- **Escape hatches** — `stage.Raw(bson.D{...})` and `expr.Raw(...)` for anything not yet covered
- **Debuggable** — `pipeline.MustJSON()` prints the final pipeline as indented JSON

---

## Installation

```bash
go get github.com/djit2026/aggify
```

---

## Quick Start

```go
import (
    "github.com/djit2026/aggify/agg"
    "github.com/djit2026/aggify/expr"
    "github.com/djit2026/aggify/q"
    "github.com/djit2026/aggify/stage"
)

pipeline := agg.New().
    Match(q.And(
        q.Eq("userId", userID),
        q.Eq("status", "active"),
    )).
    Stage(filterActiveItems()).   // reusable business-logic stage
    Stage(groupByStore()).
    SortDesc("storeTotal").
    Limit(20).
    Build()

cursor, err := collection.Aggregate(ctx, pipeline)
```

---

## Packages

| Package  | Responsibility |
|----------|----------------|
| `agg`    | Fluent pipeline builder — assembles stages via `.Stage()` |
| `stage`  | Individual stage builders (`$match`, `$group`, `$geoNear`, …) |
| `expr`   | Aggregation expression engine (`$cond`, `$filter`, `$dateAdd`, `$sin`, `$getField`, …) |
| `q`      | Query filter helpers for `$match` / `Find` (including geo filters) |
| `geo`    | GeoJSON geometry constructors (`Point`, `Polygon`, `LineString`, …) and legacy shapes |
| `search` | MongoDB Atlas Search operators (`$search`) |

---

## Package: `expr` — Expression Engine

Every MongoDB aggregation operator is a typed Go function returning `expr.Expr`.

```go
// Primitives
expr.Field("price")        // "$price"
expr.Var("item")           // "$$item"
expr.Value(42)             // 42 (literal)
expr.Root                  // "$$ROOT"
expr.Remove                // "$$REMOVE"

// Comparison (aggregation form)
expr.Eq(expr.Field("status"), expr.Value("active"))
expr.Gt(expr.Field("age"), expr.Value(18))

// Logical
expr.And(expr.Eq(...), expr.Gt(...))
expr.Or(...)
expr.Not(...)

// Conditional
expr.Cond(ifExpr, thenExpr, elseExpr)
expr.IfNull(expr.Field("discount"), expr.Value(0))
expr.Switch([]expr.SwitchBranch{
    {Case: expr.Eq(expr.Field("tier"), expr.Value("gold")), Then: expr.Value(0.2)},
}, expr.Value(0))

// Array
expr.Filter(expr.Field("items"), "item", expr.Eq(expr.Var("item.status"), expr.Value("active")))
expr.Map(expr.Field("items"), "item", expr.Field("$$item.name"))
expr.Size(expr.Field("tags"))
expr.Push(expr.Root)
expr.Sum(expr.Field("price"))
expr.Avg(expr.Field("score"))

// Object
expr.MergeObjects(expr.Root, expr.Raw(bson.D{{Key: "extra", Value: 1}}))
expr.ObjectToArray(expr.Field("meta"))

// Arithmetic
expr.Add(expr.Field("a"), expr.Field("b"))
expr.Multiply(expr.Field("price"), expr.Value(1.1))
expr.Round(expr.Field("total"), 2)

// String
expr.Concat(expr.Field("first"), expr.Value(" "), expr.Field("last"))
expr.ToLower(expr.Field("email"))
expr.RegexMatch(expr.Field("email"), `^.+@.+\..+$`, "i")
expr.RegexFindAll(expr.Field("body"), `\d+`, "")
expr.ReplaceAll(expr.Field("name"), expr.Value(" "), expr.Value("-"))

// Date & Time
expr.DateTrunc(expr.Field("createdAt"), "month")
expr.DateAdd(expr.Field("lastLogin"), expr.Value(7), "day")
expr.Year(expr.Field("createdAt"))
expr.ISOWeek(expr.Field("createdAt"))
expr.DateFromParts(expr.DateFromPartsOptions{Year: expr.Value(2024), Month: expr.Value(6), Day: expr.Value(1)})
expr.DateToParts(expr.Field("createdAt"), expr.DateToPartsOptions{ISO: true})

// Set Math
expr.SetIntersection(expr.Field("roles"), expr.Value([]string{"admin", "editor"}))
expr.SetEquals(expr.Field("tags"), expr.Value([]string{"go", "mongodb"}))

// Trigonometry & Angle Conversion
expr.Sin(expr.Field("angleRad"))
expr.Atan2(expr.Field("y"), expr.Field("x"))
expr.DegreesToRadians(expr.Field("angleDeg"))
expr.RadiansToDegrees(expr.Field("angleRad"))

// Dynamic Field Access (MongoDB 5.0+)
expr.GetField("price.usd", expr.Root)   // safe dot/dollar access
expr.SetField("status", expr.Root, expr.Value("active"))
expr.UnsetField("password", expr.Root)

// Array (MongoDB 5.2+)
expr.SortArray(expr.Field("scores"), nil, expr.Value(-1))   // sort scalar array desc
expr.FirstN(expr.Value(3), expr.Field("scores"))            // top-3 elements
expr.MaxN(expr.Value(3), expr.Field("scores"))              // 3 largest values
expr.Top(bson.D{{"score", -1}}, expr.Field("playerName"))  // highest-score doc
expr.Range(expr.Value(0), expr.Value(10), expr.Value(2))   // [0, 2, 4, 6, 8]

// Misc
expr.Literal("$notAField")   // force literal interpretation
expr.Rand()                   // random float 0–1
expr.TsSecond(expr.Field("ts"))  // BSON Timestamp seconds component

// Type conversion
expr.ToString(expr.Field("_id"))
expr.ToDate(expr.Field("createdAtStr"))

// Variable binding
expr.Let(
    []expr.LetBinding{{"tax", expr.Multiply(expr.Field("price"), expr.Value(0.1))}},
    expr.Add(expr.Field("price"), expr.Var("tax")),
)
```

---

## Package: `q` — Query Filters

Produces `bson.D` filter documents for `$match` stages and `Find` calls.

```go
q.Eq("status", "active")
q.Ne("role", "banned")
q.Gt("age", 18)
q.In("role", "admin", "editor", "mod")
q.Nin("status", "deleted", "archived")
q.Exists("deletedAt", false)
q.Regex("email", `^admin`, "i")
q.ElemMatch("items", q.And(q.Eq("status", "active"), q.Gt("qty", 0)))

// Logical combinators
q.And(q.Eq("active", true), q.Gt("age", 18))
q.Or(q.Eq("role", "admin"), q.Eq("role", "superuser"))
q.Nor(q.Eq("status", "banned"))

// Bridge to aggregation expressions
q.Expr(expr.Gt(expr.Field("spend"), expr.Field("budget")))

// Geo query filters (require 2dsphere index)
q.Near("location", geo.Point(-73.98, 40.75), q.NearOptions{MaxDistance: ptr(5000.0)})
q.NearSphere("location", geo.Point(-73.98, 40.75))
q.GeoWithin("location", geo.Polygon(ring))                            // GeoJSON shape
q.GeoWithinShape("loc", geo.Box(geo.LngLat(-74, 40), geo.LngLat(-73, 41))) // legacy box
q.GeoIntersects("coverage", geo.Point(-73.98, 40.75))
```

---

## Package: `stage` — Stage Builders

```go
// $match
stage.Match(q.Eq("status", "active"))

// $group
stage.Group(
    expr.Field("storeId"),                        // _id
    stage.Acc("total", expr.Sum(expr.Field("price"))),
    stage.Acc("items", expr.Push(expr.Root)),
)

// $project
stage.Project().
    Include("name", "email").
    Exclude("_id").
    Computed("fullName", expr.Concat(expr.Field("first"), expr.Value(" "), expr.Field("last")))

// $lookup — simple form
stage.Lookup("users", "userId", "_id", "user")

// $lookup — pipeline form
stage.LookupPipeline("assignments", "assignments").
    Let(bson.D{{Key: "orderId", Value: "$_id"}}).
    Pipeline(agg.New().Match(q.Expr(...)).Build())

// $unwind
stage.Unwind("$items")
stage.Unwind("$items").PreserveNullAndEmpty(true).IncludeArrayIndex("idx")

// $sort
stage.Sort(stage.SortField{"createdAt", stage.Desc}, stage.SortField{"name", stage.Asc})
stage.SortDesc("score")  // convenience shorthand

// $addFields / $set
stage.AddFields(stage.FE("total", expr.Sum(expr.Field("items"))))
stage.Set(stage.FE("upper", expr.ToUpper(expr.Field("name"))))

// $unset
stage.Unset("password", "internalToken")

// $replaceRoot / $replaceWith
stage.ReplaceRoot(expr.MergeObjects(expr.Field("details"), expr.Root))
stage.ReplaceWith(expr.Field("nested"))

// $count
stage.Count("total")

// $facet
stage.Facet(
    stage.FacetPipeline{Name: "byStatus", Pipeline: agg.New().Group(...).Build()},
    stage.FacetPipeline{Name: "byTier",   Pipeline: agg.New().Group(...).Build()},
)

// $bucket / $bucketAuto
stage.Bucket(expr.Field("price"), []any{0, 50, 100, 500}, "Other",
    stage.Acc("count", expr.Sum(expr.Value(1))),
)
stage.BucketAuto(expr.Field("age"), 5)

// $setWindowFields
stage.SetWindowFields(
    expr.Field("state"),
    stage.SortWindow{{"orderDate", 1}},
    stage.WindowFE("runningTotal", expr.Sum(expr.Field("price")), stage.WindowBounds("documents", "unbounded", "current")),
)

// Window-specific operators (used inside WindowFE)
stage.WindowFE("rank",     stage.WinRank(), nil)
stage.WindowFE("denseRank",stage.WinDenseRank(), nil)
stage.WindowFE("docNum",   stage.WinDocumentNumber(), nil)
stage.WindowFE("shifted",  stage.WinShift(expr.Field("price"), -1, expr.Value(0)), nil)
stage.WindowFE("cov",      stage.WinCovariancePop(expr.Field("x"), expr.Field("y")), stage.WindowBounds("documents", "unbounded", "current"))
stage.WindowFE("ema",      stage.WinExpMovingAvg(expr.Field("price"), 10, 0), nil)

// $geoNear  (must be first stage; requires 2dsphere index)
stage.GeoNear(stage.GeoNearOptions{
    Near:          geo.Point(-73.9857, 40.7580), // [lng, lat]
    DistanceField: "dist.meters",
    Spherical:     true,
    MaxDistance:   ptr(2000.0),            // 2 km
    Query:         q.Eq("category", "restaurant"),
    IncludeLocs:   "dist.location",
})

// $search (Atlas Search)
stage.Search(search.Compound().Must(search.Text("title", "golang")))

// $out / $merge
stage.Out("targetCollection")
stage.Merge(stage.MergeOptions{IntoCollection: "dailySales", WhenMatched: "replace"})

// Escape hatch
stage.Raw(bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: 5}}}})
```

---

## Package: `agg` — Pipeline Builder

```go
pipeline := agg.New().
    GeoNear(stage.GeoNearOptions{         // must be first stage
        Near:          geo.Point(-73.98, 40.75),
        DistanceField: "dist.meters",
        Spherical:     true,
    }).
    Match(filter).
    Stage(myReusableStage()).
    Group(id, accs...).
    Lookup(from, local, foreign, as).
    Unwind("$field").
    Sort(stage.SortField{"dist.meters", stage.Asc}).
    Limit(20).
    Skip(0).
    AddFields(stage.FE("x", expr.Value(1))).
    Unset("secret").
    Count("total").
    SetWindowFields(expr.Field("category"), stage.SortWindow{{"date", 1}}, stage.WindowFE("rank", stage.WinRank(), nil)).
    Out("reports").
    Raw(bson.D{...}).      // escape hatch
    Build()                // → mongo.Pipeline

// Debug the pipeline
fmt.Println(agg.New().Match(...).MustJSON())
```

---

## Composable Stage Pattern

The real power of `aggify` is the **composable stage** pattern. Any `func() stage.Stage` becomes a reusable, testable pipeline unit:

```go
// In your domain package
func FilterActiveItems() stage.Stage {
    return stage.AddFields(
        stage.FE("items", expr.Filter(
            expr.Field("items"), "item",
            expr.Eq(expr.Var("item.status"), expr.Value("active")),
        )),
    )
}

func GroupByStore() stage.Stage {
    return stage.Group(
        expr.Raw(bson.D{{Key: "storeId", Value: "$items.storeId"}}),
        stage.Acc("items", expr.Push(expr.Field("items"))),
        stage.Acc("total", expr.Sum(expr.Field("items.price"))),
    )
}

// In your service
pipeline := agg.New().
    Match(q.Eq("userId", userID)).
    Stage(orders.FilterActiveItems()).
    Stage(orders.GroupByStore()).
    SortDesc("total").
    Build()
```

Each stage is **independently testable**:

```go
func TestFilterActiveItems(t *testing.T) {
    got := orders.FilterActiveItems().Build()
    // assert the bson.D structure
}
```

---

## Type-Safe Schemas (`aggify-gen`)

Eliminate typos and get perfect IDE autocomplete by generating type-safe BSON paths from your structs.

**1. Install the tool:**
```bash
go install github.com/djit2026/aggify/cmd/aggify-gen@latest
```

**2. Annotate your models:**
```go
//go:generate aggify-gen -type User -pkg schema -out schema/schema.go

type Address struct {
    City string `bson:"city"`
}

type User struct {
    Email   string  `bson:"email"`
    Address Address `bson:"address"`
}
```

**3. Use the generated schema in your queries:**
```go
import "myproject/schema"

pipeline := agg.New().
    Match(q.Eq(schema.User.Email, "test@test.com")).
    Project(stage.Project().
        Include(schema.User.Address.City),
    ).
    Build()
```

---

## Package: `geo` — Geospatial Support

The `geo` package provides type-safe GeoJSON geometry builders and legacy shape
helpers for use with geo query operators and the `$geoNear` stage.

> ⚠️ MongoDB uses **[longitude, latitude]** order for GeoJSON — the opposite of
> the common GPS [lat, lng] convention. All helpers in this package follow the
> GeoJSON standard.

```go
import "github.com/djit2026/aggify/geo"

// GeoJSON geometries
geo.Point(-73.9857, 40.7580)             // Times Square
geo.LineString(
    geo.LngLat(-73.98, 40.75),
    geo.LngLat(-74.00, 40.71),
)
geo.Polygon([]geo.Coord{
    geo.LngLat(-73.99, 40.76),
    geo.LngLat(-73.98, 40.76),
    geo.LngLat(-73.98, 40.75),
    geo.LngLat(-73.99, 40.76), // closed
})
geo.MultiPoint(geo.LngLat(-73.98, 40.75), geo.LngLat(-74.00, 40.71))
geo.GeometryCollection(geo.Point(-73.98, 40.75), geo.LineString(...))

// Legacy shapes (flat-earth / planar)
geo.Box(geo.LngLat(-74.00, 40.71), geo.LngLat(-73.96, 40.75))
geo.Circle(geo.LngLat(-73.98, 40.75), 100)             // 100 unit radius
geo.CenterSphere(geo.LngLat(-73.98, 40.75), 5000.0/6378100.0) // 5 km in radians
geo.LegacyPolygon(geo.LngLat(-74, 40), geo.LngLat(-73, 40), geo.LngLat(-73, 41))
```

### Geo Queries — Full Example

```go
// Find the 10 closest restaurants to Times Square within 2 km
pipeline := agg.New().
    GeoNear(stage.GeoNearOptions{
        Near:               geo.Point(-73.9857, 40.7580),
        DistanceField:      "dist.meters",
        MaxDistance:        ptr(2000.0),
        Spherical:          true,
        Query:              q.Eq("category", "restaurant"),
        IncludeLocs:        "dist.location",
        DistanceMultiplier: ptr(0.001), // convert to km in the output
    }).
    Limit(10).
    Build()

// Filter a Find query by polygon
filter := q.GeoWithin("location", geo.Polygon(ring))

// Filter by bounding box (very fast, legacy flat-earth)
filter = q.GeoWithinShape("location", geo.Box(
    geo.LngLat(-74.00, 40.71),
    geo.LngLat(-73.96, 40.75),
))

// Filter by spherical radius
filter = q.GeoWithinShape("location", geo.CenterSphere(
    geo.LngLat(-73.98, 40.75),
    5000.0/6378100.0, // 5 km in radians
))

// Near query (returns sorted by distance)
filter = q.Near("location", geo.Point(-73.98, 40.75), q.NearOptions{
    MaxDistance: ptr(500.0), // metres
})

// Intersects — useful for coverage areas / service zones
filter = q.GeoIntersects("serviceArea", geo.Point(-73.98, 40.75))
```

---

## Verification

```bash
go build ./...
go test ./...
go vet ./...
```

---

## Design Principles

1. **Zero reflection** — no `reflect` package in any hot path
2. **Deterministic output** — `bson.D` (ordered) everywhere, never `bson.M`
3. **Panic on misuse** — invalid inputs (empty field names, nil required exprs) panic at startup, not at query time
4. **No abstraction leakage** — MongoDB semantics are preserved 1:1
5. **Escape hatches first** — `stage.Raw` and `expr.Raw` ensure nothing is blocked

---

## License

MIT
