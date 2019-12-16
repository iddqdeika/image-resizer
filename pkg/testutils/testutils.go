package testutils

import (
	"fmt"
	"reflect"
	"testing"
)

func MustBe(target interface{}, methodName string, expected []interface{}, params []interface{}, t *testing.T) {
	err := mustBe(target, methodName, expected, params)
	if err != nil {
		t.Fatal(err)
	}
}

//chech that method "methodName" from "config" will return "expected" on given "key".
func mustBe(target interface{}, methodName string, expected []interface{}, params []interface{}) error {

	if target == nil {
		return fmt.Errorf("target is nil")
	}

	//проверяем, что таргет - структура
	switch reflect.TypeOf(target).Kind() {
	case reflect.Ptr:
		if reflect.ValueOf(target).Elem().Kind() != reflect.Struct {
			return fmt.Errorf("target must be interface or struct")
		}
	case reflect.Struct:
	default:
		return fmt.Errorf("target must be interface or struct")
	}

	//получаем метод
	someMeth := reflect.ValueOf(target).MethodByName(methodName)
	if !someMeth.IsValid() {
		return fmt.Errorf("methodName \"%v\" from  "+getTypeName(target)+" is invalid", methodName)
	}
	//делаем слайс параметров
	inputs := make([]reflect.Value, 0)
	for _, param := range params {
		inputs = append(inputs, reflect.ValueOf(param))
	}
	//вызываем метод с параметрами
	r := someMeth.Call(inputs)
	results := make([]interface{}, 0)
	for _, v := range r {
		results = append(results, v.Interface())
	}

	//получили результат метода (первый)
	if !isEqualSlices(results, expected) {
		return fmt.Errorf("incorrect value returned for method \"%v\" of \"%v\" with AsSlice \"%v\": \"%v\"", methodName, reflect.TypeOf(target).Name(), params, results)
	}

	return nil
}

func AsSlice(args ...interface{}) []interface{} {
	return args
}

//function to compare two slices
func isEqualSlices(first interface{}, second interface{}) bool {

	firstRef := reflect.ValueOf(first)
	secondRef := reflect.ValueOf(second)

	if (first == nil) && (second == nil) {
		return true
	}
	firstlen := firstRef.Len()
	secondlen := secondRef.Len()
	if ((first == nil) != (second == nil)) || firstlen != secondlen {
		return false
	}

	for i := 0; i < firstlen; i++ {
		firstVal := firstRef.Index(i).Interface()
		secondVal := secondRef.Index(i).Interface()
		if firstVal != secondVal {
			return false
		}
	}
	return true
}

//function to get name of struct
func getTypeName(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
