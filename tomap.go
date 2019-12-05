package tomap

import (
	"reflect"
)

type Error int

func (e Error) Error() string {
	switch e {
	case NotAMap:
		return "value is not a map"
	case NotASlice:
		return "value is not a slice"
	case NotAStruct:
		return "value is not a struct"
	case IdNotFound:
		return "a struct doesnt have an ID"
	case IdIsNotUint:
		return "ID type is not uint"
	}
	return ""
}

const (
	NotAMap Error = iota
	NotASlice
	NotAStruct
	IdNotFound
	IdIsNotUint
)

// mapToSlice - construct slice from map values
func mapToSlice(val reflect.Value) (*reflect.Value, error) {
	if val.Kind() != reflect.Map {
		return nil, NotAMap
	}
	mapValType := val.Type().Elem() // type of value of a map

	if err := validateType(mapValType); err != nil {
		return nil, err
	}

	slice := reflect.MakeSlice(reflect.SliceOf(mapValType), 0, 0)
	slicePtr := reflect.New(slice.Type())
	slicePtr.Elem().Set(slice)
	return &slicePtr, nil
}

// sliceToIdMap - fills map(where key is uint value of ID field of slice element and value is slice element type)
func sliceToIdMap(sliceVal, mappVal reflect.Value) error {
	if sliceVal.Kind() != reflect.Slice {
		return NotASlice
	}
	if mappVal.Kind() != reflect.Map {
		return NotAMap
	}
	for i := 0; i < sliceVal.Len(); i++ {
		nextVal := sliceVal.Index(i)
		if err := validateType(nextVal.Type()); err != nil {
			return err
		}
		id := reflect.Indirect(nextVal).FieldByName("ID")
		mappVal.SetMapIndex(id, nextVal)
	}
	return nil
}

// validateType - validate that type can be used for "mapping"
func validateType(ttype reflect.Type) error {
	if ttype.Kind() != reflect.Struct {
		if ttype.Kind() == reflect.Ptr && ttype.Elem().Kind() != reflect.Struct { // if value is not a pointer to a struct
			return NotAStruct
		}
	}
	if ttype.Kind() == reflect.Ptr { // "indirect" if value is a pointer to a struct
		ttype = ttype.Elem()
	}
	id, ok := ttype.FieldByName("ID")
	if !ok {
		return IdNotFound
	}
	if id.Type.Kind() != reflect.Uint {
		return IdIsNotUint
	}
	return nil
}