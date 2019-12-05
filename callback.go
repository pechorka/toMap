// tomap adds to gorm functionality of returning query data as map[uint]ModelType.
// key - is `ID` field of a model
package tomap

import (
	"github.com/jinzhu/gorm"
	"reflect"
)

const (
	toSlice     = "tomap:toslice"
	toMap       = "tomap:tomap"
	originalMap = "tomap:originalmap"
)

// beforeQuery replaces passed to scope map with slice of same structure as map value.
// map key should be uint and map value should be a struct with `ID` field
func beforeQuery(scope *gorm.Scope) {
	origVal := scope.Value
	val := reflect.ValueOf(origVal)
	sliceVal, err := mapToSlice(reflect.Indirect(val))
	if err != nil && err != NotAMap { // value of map is not a struct or doesnt have an `ID` field
		_ = scope.Err(err) // ignore because this method returns added error
		return
	}
	if err == nil { // if value is not a map, than err won't be nil, so gorm behavior won't change for standard supported types
		scope.Value = sliceVal.Interface()
		scope.Set(originalMap, origVal)
	}
}

// afterQuery - if successfully retries map from scope setting, than fills it
// with data from slice by the following rule: ID field of structure is key and structure is value
func afterQuery(scope *gorm.Scope) {
	if oMap, ok := scope.Get(originalMap); ok {
		sliceVal := reflect.ValueOf(scope.Value)
		valMap := reflect.ValueOf(oMap)
		err := sliceToIdMap(reflect.Indirect(sliceVal), reflect.Indirect(valMap))
		if err != nil {
			_ = scope.Err(err) // ignore because this method returns added error
			return
		}
		scope.Value = oMap
	}
}

// RegisterCallbacks - register callbacks for passed database instance
func RegisterCallbacks(db *gorm.DB) {
	callback := db.Callback()
	if callback.Query().Get(toSlice) == nil {
		callback.Query().Before("gorm:query").Register(toSlice, beforeQuery)
	}
	if callback.Query().Get(toMap) == nil {
		callback.Query().After("gorm:query").Register(toMap, afterQuery)
	}
}
