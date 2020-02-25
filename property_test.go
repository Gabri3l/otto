package otto

import (
	"testing"
)

func TestPropertyWritableGet(t *testing.T) {
	tt(t, func() {
		test, _ := test()

		test(`
var obj = {}

Object.defineProperty(obj, 'getter', {
	get: function() { return 2 },
	configurable: true
});

Object.getOwnPropertyDescriptor(obj, 'getter');

Object.defineProperty(obj, 'getter', {
	configurable: true,
	writable: true
});

Object.getOwnPropertyDescriptor(obj, 'getter');
`)
	})
}
