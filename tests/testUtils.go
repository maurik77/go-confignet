package tests

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
	ArrayStr       []string
	ArrayInt       []int
}

func validateObject(t *testing.T, expected myConfig, result myConfig) {

	if result.PropertyInt8 != expected.PropertyInt8 {
		t.Logf("validateObject::error should be '%v', but got '%v'", expected.PropertyInt8, result.PropertyInt8)
		t.Fail()
	}

	validateSubObject(t, expected.Obj1, result.Obj1)
}

func validateSubObject(t *testing.T, expected subObj, result subObj) {
	if result.PropertyString != expected.PropertyString {
		t.Logf("validateSubObject::error should be '%v', but got '%v'", expected.PropertyString, result.PropertyString)
		t.Fail()
	}

	if result.PropertyBool != expected.PropertyBool {
		t.Logf("validateSubObject::error should be '%v', but got '%v'", expected.PropertyBool, result.PropertyBool)
		t.Fail()
	}

	if result.PropertyInt != expected.PropertyInt {
		t.Logf("validateSubObject::error should be '%v', but got '%v'", expected.PropertyInt, result.PropertyInt)
		t.Fail()
	}

	if result.PropertyInt8 != expected.PropertyInt8 {
		t.Logf("validateSubObject::error should be '%v', but got '%v'", expected.PropertyInt8, result.PropertyInt8)
		t.Fail()
	}

	if result.PropertyInt16 != expected.PropertyInt16 {
		t.Logf("validateSubObject::error should be '%v', but got '%v'", expected.PropertyInt16, result.PropertyInt16)
		t.Fail()
	}

	if result.PropertyInt64 != expected.PropertyInt64 {
		t.Logf("validateSubObject::error should be '%v', but got '%v'", expected.PropertyInt64, result.PropertyInt64)
		t.Fail()
	}

	if result.Time != expected.Time {
		t.Logf("validateSubObject::error should be '%v', but got '%v'", expected.Time, result.Time)
		t.Fail()
	}
}
