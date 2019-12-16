package testutils

import (
	"fmt"
	"reflect"
	"testing"
)

func MustBe(target interface{}, methodName string, expected []interface{}, params []interface{}, t *testing.T) {
	err := mustBe(target, methodName, expected, params)
	if err != nil{
		t.Fatal(err)
	}
}

//chech that method "methodName" from "config" will return "expected" on given "key".
func mustBe(target interface{}, methodName string, expected []interface{}, params []interface{}) error {

	if target == nil{
		return fmt.Errorf("target is nil")
	}

	//проверяем, что таргет - структура
	switch reflect.TypeOf(target).Kind() {
	case reflect.Ptr:
		if reflect.ValueOf(target).Elem().Kind() != reflect.Struct{
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
	for _, param := range params{
		inputs = append(inputs, reflect.ValueOf(param))
	}
	//вызываем метод с параметрами
	r := someMeth.Call(inputs)
	results := make([]interface{}, 0)
	for _, v := range r{
		results = append(results, v.Interface())
	}

	//получили результат метода (первый)
	if !isEqualSlices(results, expected){
		return fmt.Errorf("incorrect value returned for method \"%v\" of \"%v\" with AsSlice \"%v\": \"%v\"", methodName, reflect.TypeOf(target).Name(), params, results)
	}

	//получаем значение результата как интерфейс
	//value := result.Interface()
	//
	//kind := result.Kind()
	//switch kind {
	//case reflect.Map:
	//	expval := reflect.ValueOf(expected)
	//	if (!expval.IsValid() || expval.Pointer() == 0) && (!result.IsValid() || result.Pointer() == 0) {
	//		return nil
	//	}
	//	if !(reflect.ValueOf(expected).Kind() == reflect.Map) {
	//		return fmt.Errorf("given expected value for params \"%v\" for method \"%v\" from %v is \"%v\", but not slice", params, methodName, getTypeName(target), expected)
	//	}
	//	if len(result.MapKeys()) != len(reflect.ValueOf(expected).MapKeys()) {
	//		return fmt.Errorf("values count from params \"%v\" on method \"%v\" from %v is \"%v\", but is not equal to given expected values count \"%v\"", params, methodName, getTypeName(target), value, expected)
	//	}
	//	for _, v := range result.MapKeys() {
	//		res := result.MapIndex(v).Interface()
	//		ex := expval.MapIndex(v).Interface()
	//		//floattype := reflect.TypeOf(float64(0))
	//		if res != ex {
	//			return fmt.Errorf("value from params \"%v\" on method \"%v\" from %v is \"%v\", but is not equal to given expected value \"%v\"", params, methodName, getTypeName(target), value, expected)
	//
	//		}
	//	}
	//
	//case reflect.Slice:
	//
	//	expval := reflect.ValueOf(expected)
	//	if (!expval.IsValid() || expval.Pointer() == 0) && (!result.IsValid() || result.Pointer() == 0) {
	//		return nil
	//	}
	//	if !(reflect.ValueOf(expected).Kind() == reflect.Slice) {
	//		return fmt.Errorf("given expected value for params \"%v\" for method \"%v\" from %v is \"%v\", but not slice", params, methodName, getTypeName(target), expected)
	//	}
	//	if !isEqualSlices(value, expected) {
	//		return fmt.Errorf("value from params \"%v\" on method \"%v\" from %v is \"%v\", but is not equal to given expected value \"%v\"", params, methodName, getTypeName(target), value, expected)
	//	}
	//case reflect.Ptr:
	//	if value == nil {
	//		return fmt.Errorf("value for params \"%v\" from method \"%v\" from %v must not be nil", params, methodName, getTypeName(target))
	//	}
	//	if !result.Elem().IsValid() {
	//		if expected == nil {
	//			return nil
	//		}
	//		return fmt.Errorf("value for params \"%v\" from method \"%v\" from %v must not be nil", params, methodName, getTypeName(target))
	//	}
	//
	//	//проверяем типы значений в результате и ожидаемом значении
	//	//(чтобы в случае разных типов выделить это в отдельную ошибку)
	//	resultType := reflect.ValueOf(value).Elem().Type()
	//	expType := reflect.TypeOf(expected)
	//	if resultType != expType {
	//		return fmt.Errorf("for params \"%v\" from method %v from %v got type \"%v\", expected \"%v\"", params, methodName, getTypeName(target), reflect.TypeOf(value).Name(), reflect.TypeOf(expected).Name())
	//	}
	//	tName := result.Elem().Type().Name()
	//	//сравнение значений в зависимости от типа
	//	switch tName {
	//	case "Time":
	//		//сравниваем значения в результате и ожидаемом значении
	//		vint := reflect.ValueOf(value).Elem().Interface()
	//		if !vint.(time.Time).Equal(expected.(time.Time)) {
	//			return fmt.Errorf("incorrect value \"%v\" for params \"%v\" from method \"%v\" from %v, expected \"%v\"", value, params, methodName, getTypeName(target), expected)
	//		}
	//	default:
	//		//сравниваем значения в результате и ожидаемом значении
	//		vin := reflect.ValueOf(value).Elem().Interface()
	//		if vin != expected {
	//			return fmt.Errorf("incorrect value \"%v\" for params \"%v\" from method \"%v\" from %v, expected \"%v\"", value, params, methodName, getTypeName(target), expected)
	//		}
	//
	//	}
	//default:
	//	if value == nil {
	//		return fmt.Errorf("value for params \"%v\" from method \"%v\" from %v must not be nil", params, methodName, getTypeName(target))
	//	}
	//	if !result.Elem().IsValid() {
	//		if expected == nil {
	//			return nil
	//		}
	//		return fmt.Errorf("value for params \"%v\" from method \"%v\" from %v must not be nil", params, methodName, getTypeName(target))
	//	}
	//	//проверяем типы значений в результате и ожидаемом значении
	//	//(чтобы в случае разных типов выделить это в отдельную ошибку)
	//	resultType := reflect.ValueOf(value).Elem().Type()
	//	expType := reflect.TypeOf(expected)
	//	if resultType != expType {
	//		return fmt.Errorf("for params \"%v\" from method %v from %v got type \"%v\", expected \"%v\"", params, methodName, getTypeName(target), reflect.TypeOf(value).Name(), reflect.TypeOf(expected).Name())
	//	}
	//	//проверяем значения в результате и ожидаемом значении
	//	if reflect.ValueOf(value).Elem().Interface() != expected {
	//		return fmt.Errorf("incorrect value \"%v\" for params \"%v\" from method \"%v\" from %v, expected \"%v\"", value, params, methodName, getTypeName(target), expected)
	//	}
	//}
	//

	return nil
}

func AsSlice(args ...interface{}) []interface{}{
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
