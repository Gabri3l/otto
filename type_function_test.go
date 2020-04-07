package otto

import (
	"testing"
)

func TestCreateNativeFunction(t *testing.T) {
	tt(t, func() {
		vm := New()

		nativeFunc, err := vm.CreateNativeFunction("myFunc", "somerandomfile.js", 23, func(c FunctionCall) Value {
			return toValue_string(c.ArgumentList[0].String() + " world")
		})
		is(err, nil)

		nativeObj := nativeFunc._object().value.(_nativeFunctionObject)
		is(nativeObj.name, "myFunc")
		is(nativeObj.file, "somerandomfile.js")
		is(nativeObj.line, 23)

		is(vm.Set("callTheFunc", nativeFunc), nil)

		ret, err := vm.Run(`callTheFunc("hello")`)
		is(err, nil)
		is(ret, "hello world")

		_, err = vm.CreateNativeFunction("myFunc", "somerandomfile.js", 23, nil)
		is(err, "!=", nil)
	})
}

func TestFunctionSetNameProperty(t *testing.T) {
	vm := New()

	val, err := vm.Run(`(function() { 
		var testFunc = function() {}; 
		Object.defineProperty(testFunc, "name", {value: "hello"});
		return testFunc;
	})()`)
	is(err, nil)
	x, err := val.Object().Get("name")
	is(err, nil)
	is(x, "hello")

	val, err = vm.Run(`(function() { 
		var testFunc = function hi() {}; 
		Object.defineProperty(testFunc, "name", {value: "hello"});
		return testFunc;
	})()`)
	is(err, nil)
	x, err = val.Object().Get("name")
	is(err, nil)
	is(x, "hello")

	// verify function name isn't enumerable, taken from global_test
	call := func(object interface{}, src string, argumentList ...interface{}) Value {
		var tgt *Object
		switch object := object.(type) {
		case Value:
			tgt = object.Object()
		case *Object:
			tgt = object
		case *_object:
			tgt = toValue_object(object).Object()
		default:
			panic("Here be dragons.")
		}
		value, err := tgt.Call(src, argumentList...)
		is(err, nil)
		return value
	}

	is(call(val._object(), "propertyIsEnumerable", "name"), false)
}
