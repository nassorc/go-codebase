package gandalf

// import (
// 	"testing"
// )

// type Position struct {
// 	X int
// 	Y int
// }

// type Velocity struct {
// 	X int
// 	Y int
// }

// type Tag struct {
// 	tag string
// }

// func PositionSystem(world *World, entities []*Entity) {
// 	for _, entity := range entities {
// 		var (
// 			pos *Position
// 			vel *Velocity
// 		)
// 		world.GetData(entity, &pos, &vel)

// 		pos.X += vel.X
// 		pos.Y += vel.Y
// 	}
// }

// func TagSystem(world *World, entities []*Entity) {
// 	for _, entity := range entities {
// 		var tag *Tag
// 		world.GetData(entity, &tag)

// 	}
// }

// func TestRegisterSystem(t *testing.T) {
// 	var world *World = NewWorld()

// 	world.RegisterComponents(&Position{}, &Velocity{}, &Tag{})
// 	world.RegisterSystem(PositionSystem, &Position{}, &Velocity{})

// 	// test if system receives the entity that has the required system components
// 	var testEntity = world.CreateEntity(&Position{10, 15}, &Tag{"test entity"}, &Velocity{1, 1})

// 	world.Tick()

// 	var pos *Position
// 	world.GetData(testEntity, &pos)

// 	if pos.X != 11 {
// 		t.Errorf("Wrong x value, Expected=%d, Got=%d", 11, pos.X)
// 	}
// 	if pos.Y != 16 {
// 		t.Errorf("Wrong x value, Expected=%d, Got=%d", 16, pos.Y)
// 	}

// 	world.Tick()

// 	if pos.X != 12 {
// 		t.Errorf("Wrong x value, Expected=%d, Got=%d", 12, pos.X)
// 	}
// 	if pos.Y != 17 {
// 		t.Errorf("Wrong x value, Expected=%d, Got=%d", 17, pos.Y)
// 	}
// }

// func TestRegisterTwoSystems(t *testing.T) {
// 	var world *World = NewWorld()

// 	world.RegisterComponents(&Position{}, &Tag{}, &Velocity{})

// 	world.RegisterSystem(TagSystem, &Tag{})
// 	world.RegisterSystem(PositionSystem, &Position{}, &Velocity{})

// 	var testEntity = world.CreateEntity(&Position{10, 15}, &Velocity{2, 2}, &Tag{"test entity"})

// 	world.Tick()

// 	var pos *Position
// 	world.GetData(testEntity, &pos)

// 	if pos.X != 12 {
// 		t.Errorf("Wrong x value, Expected=%d, Got=%d", 12, pos.X)
// 	}
// 	if pos.Y != 17 {
// 		t.Errorf("Wrong x value, Expected=%d, Got=%d", 17, pos.Y)
// 	}
// }

// func TestEntity(t *testing.T) {
// 	var world *World = NewWorld()

// 	world.RegisterComponents(&Position{})
// 	world.RegisterSystem(PositionSystem, &Position{})

// 	var testEntity1 = world.CreateEntity(&Position{10, 15})
// 	var testEntity2 = world.CreateEntity(&Position{98, 99})
// }

// func TestRegisterComponent(t *testing.T) {
// 	var world *World = NewWorld()
// 	world.RegisterComponents(&TestComponentA{}, &TestComponentB{})

// 	var testEntity *Entity = world.CreateEntity(&TestComponentA{100}, &TestComponentB{"Hello, World!"})
// 	var actual *TestComponentA

// 	testEntity.GetData(&actual)

// 	if actual == nil {
// 		t.Errorf("Received wrong component data. Expected=%d, Got=%v", 100, actual)
// 	}
// 	if actual.value != 100 {
// 		t.Errorf("Received wrong component data. Expected=%d, Got=%v", 100, actual.value)
// 	}
// }

// func TestRegisterComponent_PanicsGivenNonPointerArg(t *testing.T) {
// 	defer func(t *testing.T) {
// 		if r := recover(); r == nil {
// 			t.Errorf("RegisterComponent did not panic.")
// 		}
// 	}(t)
// 	var world *World = NewWorld()
// 	world.RegisterComponents(TestComponentA{})
// }

// func TestRegisterSystem(t *testing.T) {
// 	var world *World = NewWorld()
// 	world.RegisterComponents(&TestComponentA{}, &TestComponentB{})

// 	var testSystem = &TestABSystem{}
// 	world.RegisterSystem(testSystem, &TestComponentB{})

// 	type TestValue = struct {
// 		id    int
// 		value string
// 	}

// 	tt := []TestValue{}

// 	// add entity that should exist in the system
// 	entityNotInSystem := world.CreateEntity(&TestComponentA{10})

// 	entity := world.CreateEntity(&TestComponentB{"A"})
// 	tt = append(tt, TestValue{entity.Id(), "A"})

// 	entity = world.CreateEntity(&TestComponentA{50}, &TestComponentB{"B"})
// 	tt = append(tt, TestValue{entity.Id(), "B"})

// 	entity = world.CreateEntity(&TestComponentB{"C"}, &TestComponentA{50})
// 	tt = append(tt, TestValue{entity.Id(), "C"})

// 	got := []TestValue{}

// 	for _, entity := range testSystem.entities {
// 		var componentB *TestComponentB
// 		entity.GetData(&componentB)

// 		got = append(got, TestValue{entity.Id(), componentB.value})
// 	}

// 	// checks if all expected data exists in the system's entity list
// 	for _, expected := range tt {
// 		if !slices.Contains(got, expected) {
// 			t.Errorf("Expected system to have entity with the value of %s.", expected.value)
// 		}
// 	}

// 	// check if the system doesn't have the wrong entity
// 	for _, entity := range testSystem.entities {
// 		if entity.Id() == entityNotInSystem.Id() {
// 			t.Errorf("System has wrong entity")
// 		}
// 	}

// }

// func TestRegisterSystem_DoesNotRegisterInvalidComponent(t *testing.T) {
// 	var world *World = NewWorld()
// 	world.RegisterComponents(&TestComponentA{}, &TestComponentB{})

// 	defer func(t *testing.T) {
// 		if r := recover(); r == nil {
// 			t.Errorf("RegisterSystem did not panic when given non pointer to a component.")
// 		}
// 	}(t)

// 	world.RegisterSystem(&TestABSystem{}, TestComponentB{})
// }

// func TestCreateEntity(t *testing.T) {
// 	var world *World = NewWorld()
// 	world.RegisterComponents(&TestComponentA{}, &TestComponentB{})

// 	var testEntity1 *Entity = world.CreateEntity(&TestComponentA{100}, &TestComponentB{"Hello, World!"})
// 	var testEntity2 *Entity = world.CreateEntity(&TestComponentA{40}, &TestComponentB{"player"})

// 	var compBEntity1 *TestComponentB
// 	var compBEntity2 *TestComponentB

// 	testEntity1.GetData(&compBEntity1)
// 	testEntity2.GetData(&compBEntity2)

// 	if compBEntity1.value != "Hello, World!" {
// 		t.Errorf("Expected entity to have component with value %s, Got=%s", "Hello, World!", compBEntity1.value)
// 	}
// 	if compBEntity2.value != "player" {
// 		t.Errorf("Expected entity to have component with value %s, Got=%s", "player", compBEntity2.value)
// 	}
// }

// func TestRemoveEntity(t *testing.T) {
// 	var world *World = NewWorld()
// 	world.RegisterComponents(&TestComponentA{}, &TestComponentB{})

// 	var testEntity1 *Entity = world.CreateEntity(&TestComponentA{100}, &TestComponentB{"Hello, World!"})

// 	world.RemoveEntity(testEntity1)
// }

// type TestComponentA struct {
// 	value int
// }

// type TestComponentB struct {
// 	value string
// }

// type TestABSystem struct {
// 	entities    []*Entity
// 	updateCount int
// }

// func (s *TestABSystem) Update(w *World) {
// 	s.updateCount += 1
// }
// func (s *TestABSystem) Render(w *World) {}
// func (s *TestABSystem) AddEntity(e *Entity) {
// 	s.entities = append(s.entities, e)
// }
