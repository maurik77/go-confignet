package confignet

import (
	"testing"
	"time"
)

type myConfig struct {
	Obj1         subObj
	PropertyInt8 int8
}

type subObj struct {
	PropertyString string
	PropertyInt    int
	PropertyInt8   int8
	PropertyInt16  int16
	PropertyInt64  int64
	PropertyBool   bool
	Time           time.Time
}

func validateSubObject(t *testing.T, expected subObj, result subObj) {
	if result.PropertyBool != expected.PropertyBool {
		t.Log("error should be", expected.PropertyBool, ", but got", result.PropertyBool)
		t.Fail()
	}

	if result.PropertyInt != expected.PropertyInt {
		t.Log("error should be", expected.PropertyInt, ", but got", result.PropertyInt)
		t.Fail()
	}

	if result.PropertyInt8 != expected.PropertyInt8 {
		t.Log("error should be", expected.PropertyInt8, ", but got", result.PropertyInt8)
		t.Fail()
	}

	if result.PropertyInt16 != expected.PropertyInt16 {
		t.Log("error should be", expected.PropertyInt16, ", but got", result.PropertyInt16)
		t.Fail()
	}

	if result.PropertyInt64 != expected.PropertyInt64 {
		t.Log("error should be", expected.PropertyInt64, ", but got", result.PropertyInt64)
		t.Fail()
	}

	if result.Time != expected.Time {
		t.Log("error should be", expected.Time, ", but got", result.Time)
		t.Fail()
	}
}
