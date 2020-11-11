package otto

type _symbolObject struct {
	internalVal interface{} // TODO: not sure internalVal is a good name for this
	description interface{}
}

func (runtime *_runtime) newSymbolObject(description interface{}) *_object {
	self := runtime.newObject()
	self.class = "Symbol"

	symbol := _symbolObject{
		description: description,
	}
	// TODO: we can convert it directly to a string here so we don't have to worry
	// about it later
	symbol.internalVal = &symbol
	self.value = symbol

	self.defineProperty("description", toValue(description), 0000, false)

	if _, ok := runtime.symbols[description]; !ok {
		runtime.symbols[description] = toValue_object(self)
	}

	return self
}
