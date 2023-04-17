package tests

import (
	"reflect"
	"testing"
	"time"
)

type myConfig struct {
	Obj1         *subObj
	PropertyInt8 *int8
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
	ArrayInt       *[3]int
	ArrayObj       []subObjItem
	ArrayObjPtr    []*subObjItem
	MapStr         map[string]string
	MapInt         map[int]int
	MapObj         map[int]subObjItem
	MapObjPtr      map[bool]*subObjItem
}

type subObjItem struct {
	PropertyString string
	PropertyInt    int
	PropertyBool   bool
}

func validateObject(t *testing.T, expected myConfig, result myConfig) {

	if *result.PropertyInt8 != *expected.PropertyInt8 {
		t.Logf("validateObject::error should be '%v', but got '%v'", expected.PropertyInt8, result.PropertyInt8)
		t.Fail()
	}

	validateSubObject(t, *expected.Obj1, *result.Obj1)
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

	if !reflect.DeepEqual(result.ArrayInt, expected.ArrayInt) {
		t.Logf("validateSubObject::error should be '%v', but got '%v'", expected.ArrayInt, result.ArrayInt)
		t.Fail()
	}

	if !reflect.DeepEqual(result.ArrayStr, expected.ArrayStr) {
		t.Logf("validateSubObject::error should be '%v', but got '%v'", expected.ArrayStr, result.ArrayStr)
		t.Fail()
	}

	if !reflect.DeepEqual(result.ArrayObj, expected.ArrayObj) {
		t.Logf("validateSubObject::error should be '%v', but got '%v'", expected.ArrayObj, result.ArrayObj)
		t.Fail()
	}

	if !reflect.DeepEqual(result.ArrayObjPtr, expected.ArrayObjPtr) {
		t.Logf("validateSubObject::error should be '%v', but got '%v'", expected.ArrayObj, result.ArrayObj)
		t.Fail()
	}

	if !reflect.DeepEqual(result.MapObj, expected.MapObj) {
		t.Logf("validateSubObject::error should be '%v', but got '%v'", expected.ArrayObj, result.ArrayObj)
		t.Fail()
	}

	if !reflect.DeepEqual(result.MapObjPtr, expected.MapObjPtr) {
		t.Logf("validateSubObject::error should be '%v', but got '%v'", expected.ArrayObj, result.ArrayObj)
		t.Fail()
	}
}

func getJSONExpectedValue() myConfig {
	var pointerInt8 int8 = 45
	expected := myConfig{
		PropertyInt8: &pointerInt8,
		Obj1: &subObj{
			PropertyString: "TestObj1",
			PropertyInt:    1,
			PropertyInt8:   2,
			PropertyInt16:  3,
			PropertyInt64:  4,
			PropertyBool:   true,
			ArrayStr:       []string{"Test", "Test2"},
			ArrayInt:       &[3]int{1, 2},
			ArrayObj: []subObjItem{
				{
					PropertyString: "TestArrObj1",
					PropertyInt:    1,
					PropertyBool:   true,
				},
				{
					PropertyString: "TestArrObj2",
					PropertyInt:    2,
					PropertyBool:   false,
				},
			},
			ArrayObjPtr: []*subObjItem{
				{
					PropertyString: "TestArrObj1",
					PropertyInt:    1,
					PropertyBool:   true,
				},
				{
					PropertyString: "TestArrObj2",
					PropertyInt:    2,
					PropertyBool:   false,
				},
			},
			MapStr: map[string]string{"Key1": "Value1", "Key2": "Value2"},
			MapInt: map[int]int{1: 1, 2: 2},
			MapObj: map[int]subObjItem{
				1: {
					PropertyString: "TestArrObj1",
					PropertyInt:    1,
					PropertyBool:   true,
				},
				2: {
					PropertyString: "TestArrObj2",
					PropertyInt:    2,
					PropertyBool:   false,
				},
				99: {
					PropertyString: "MapObjStrYaml",
					PropertyInt:    87,
					PropertyBool:   true,
				},
			},
			MapObjPtr: map[bool]*subObjItem{
				false: {
					PropertyString: "TestArrObj1",
					PropertyInt:    1,
					PropertyBool:   true,
				},
				true: {
					PropertyString: "TestArrObj2",
					PropertyInt:    2,
					PropertyBool:   false,
				},
			},
		},
	}

	timeCfg, _ := time.Parse(time.RFC3339Nano, "2022-01-19T10:00:00Z")
	expected.Obj1.Time = timeCfg
	return expected
}
