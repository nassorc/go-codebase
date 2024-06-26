package gandalf

// test
// creating component array
// adding data
// getting data
// updating data
// removing

// check bookkeeping

// type Position struct {
// 	X int
// 	Y int
// }

// func TestComponentStoreRemove(t *testing.T) {
// 	componentType := reflect.TypeOf(&Position{})
// 	store := NewComponentStore(componentType)

// 	val := reflect.ValueOf(&Position{10, 20})
// 	store.Push(0, val)

// 	val = reflect.ValueOf(&Position{100, 110})
// 	store.Push(1, val)

// 	val = reflect.ValueOf(&Position{5000, 6000})
// 	store.Push(2, val)

// 	if store.Size() != 3 {
// 		t.Errorf("Expected=%v, Got=%v", 3, store.Size())
// 	}

// 	if _, ok := store.Get(0); ok != true {
// 		t.Errorf("Expected=%v, Got=%v", true, ok)
// 	}

// 	store.Remove(0)
// 	store.Remove(2)
// 	store.Remove(1)

// 	if _, ok := store.Get(0); ok != false {
// 		t.Errorf("Expected=%v, Got=%v", false, ok)
// 	}

// 	if store.Size() != 0 {
// 		t.Errorf("Expected=%v, Got=%v", 0, store.Size())
// 	}
// }

// func TestComponentHandle(t *testing.T) {
// 	componentType := reflect.TypeOf(&Position{})
// 	store := NewComponentStore(componentType)

// 	val := reflect.ValueOf(&Position{10, 20})
// 	store.Push(0, val)

// 	// var handle ComponentHandle[*Position]
// }
