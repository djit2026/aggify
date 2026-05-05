package agg

import (
	"reflect"
	"strings"
	"sync"
)

var (
	typeCache sync.Map
)

// Key returns the BSON key path for a given struct field of type T.
// It reflects the `bson` tag of the field and caches the result for future calls
// to avoid reflection overhead in hot paths.
//
//	agg.Key[User]("Email") // returns "email" (assuming `bson:"email"` is set)
func Key[T any](fieldName string) string {
	var zero T
	t := reflect.TypeOf(zero)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	
	if t == nil {
		panic("mono-query/agg: Key[T] used with untyped nil")
	}
	if t.Kind() != reflect.Struct {
		panic("mono-query/agg: Key[T] used with non-struct type " + t.Name())
	}

	cacheKey := t.String()
	var fieldMap map[string]string

	if cached, ok := typeCache.Load(cacheKey); ok {
		fieldMap = cached.(map[string]string)
	} else {
		fieldMap = make(map[string]string)
		buildFieldMap(t, fieldMap, "")
		typeCache.Store(cacheKey, fieldMap)
	}

	if val, ok := fieldMap[fieldName]; ok {
		return val
	}
	panic("mono-query/agg: field " + fieldName + " not found on type " + t.Name())
}

func buildFieldMap(t reflect.Type, m map[string]string, prefix string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		
		if !field.IsExported() {
			continue
		}

		bsonTag := field.Tag.Get("bson")
		if bsonTag == "-" {
			continue
		}

		name := field.Name
		if bsonTag != "" {
			parts := strings.Split(bsonTag, ",")
			if parts[0] != "" {
				name = parts[0]
			}
		}

		fullPath := name
		if prefix != "" {
			fullPath = prefix + "." + name
		}

		m[field.Name] = fullPath
	}
}
