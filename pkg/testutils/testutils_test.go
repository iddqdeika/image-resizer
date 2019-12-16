package testutils

import "testing"

func TestIsEqualsSlices(t *testing.T) {
	val1 := []interface{}{1,2,3}
	val2 := []interface{}{1,2,3}
	val3 := []interface{}{3,2,1}
	if !isEqualSlices(val1, val2){
		t.Errorf("must be equal 1,2,3 and 1,2,3")
	}
	if isEqualSlices(val2, val3){
		t.Errorf("must not be equal 1,2,3 and 3,2,1")
	}
}

func TestAsSlice(t *testing.T) {
	if !isEqualSlices(AsSlice(1,2,"3"), []interface{}{1,2,"3"}){
		t.Errorf("AsSlice returns not equal slice of interfaces for equal input")
	}
	if isEqualSlices(AsSlice(""), []interface{}{}){
		t.Errorf("AsSlice returns equal slices for non equal input")
	}
}

func TestMustBe(t *testing.T) {
	err := mustBe(nil, "meth", AsSlice("someval"), nil)
	if err ==nil{
		t.Errorf("must be error if got nil target")
	}

	var ti testInterface = &testStruct{}

	err = mustBe(ti, "SomeMethod", AsSlice(2), AsSlice(1))
	if err != nil{
		t.Errorf("must pass with correct interface realization, but got err: %v", err)
	}

	err = mustBe(ti, "SomeAnotherMeth", nil, nil)
	if err != nil{
		t.Errorf("must pass with void method without params with expected and params set like nil")
	}

	err = mustBe(ti, "SomeAnotherMeth", []interface{}{1}, nil)
	if err == nil{
		t.Errorf("must not pass with void method without params and expected set not to nil")
	}

	defer func() {
		err := recover()
		if err == nil{
			t.Errorf("must panic")
		}
	}()
	mustBe(ti, "SomeAnotherMeth", nil, []interface{}{1})
}

type testInterface interface {
	SomeMethod(i int)int
}

type testStruct struct {

}

func (s *testStruct) SomeMethod(i int) int {
	return i + i
}

func (s *testStruct) SomeAnotherMeth(){

}