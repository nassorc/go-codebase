package gandalf

import "reflect"

type ComponentID = reflect.Type

func CreateComponentID[C any]() ComponentID {
	var id = reflect.TypeOf((*C)(nil))
	return id
}
