package utils

import (
	"math/big"
	"reflect"
)

func ComparePtrFieldsDesc(s1, s2 interface{}) bool {
	// Get the int values directly from the underlying pointer values
	val1 := reflect.ValueOf(s1).Elem().FieldByName("Time").Elem().Interface().(big.Int)
	val2 := reflect.ValueOf(s2).Elem().FieldByName("Time").Elem().Interface().(big.Int)

	valInInt := val1.Int64()
	valInInt2 := val2.Int64()
	return valInInt > valInInt2 // Reverse comparison for descending order
}
