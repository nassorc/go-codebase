package gandalf

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

var (
	CountID = CreateComponentID[Count]()
	DummyID = CreateComponentID[Dummy]()
)

type Dummy struct {
	val int
}

type Count struct {
	val int
}

type Counter struct {
	entities []EntityHandle
}

func (c *Counter) Increment(entities []EntityHandle) {
	// save information
	fmt.Println(entities)
	c.entities = entities

	for _, entity := range entities {
		var count *Count
		entity.Unpack(&count)

		count.val += 1
	}
}

func DummySystem(_ []EntityHandle) {}

type MockSystem struct {
	count int
}

func (m *MockSystem) UpdateMock(_ []EntityHandle) {
	m.count += 1
}

func TestSystemFunctionCalled(t *testing.T) {
	world := NewWorld(10, 5)
	mockSys := MockSystem{0}

	world.RegisterSystem(mockSys.UpdateMock)

	world.Tick()
	world.Tick()
	world.Tick()

	if mockSys.count != 3 {
		t.Errorf("Expected=%v, Got=%v", 3, mockSys.count)
	}
}

func TestSystemEntities(t *testing.T) {
	var components = []ComponentID{CountID, DummyID}
	var counter = Counter{}

	world := NewWorld(10, 5)

	// Register Components
	world.RegisterComponents(components...)

	// register systems
	world.RegisterSystem(counter.Increment, CountID)
	world.RegisterSystem(DummySystem, DummyID)

	// create entities
	world.Create(&Dummy{1})
	c1 := world.Create(&Count{10})
	c2 := world.Create(&Count{109})
	world.Create(&Dummy{2})

	type T struct {
		id            int
		expectedCount int
	}

	tt := []T{
		{c1.Id(), 11},
		{c2.Id(), 110},
	}

	world.Tick()

	if len(counter.entities) != 2 {
		t.Errorf("Expected=%v, Got=%v", 2, len(counter.entities))
	}

	// validate the entities that the Increment function received
	for _, entity := range counter.entities {
		var id = entity.Id()
		var count *Count
		entity.Unpack(&count)

		got := T{id, count.val}

		if !slices.Contains(tt, got) {
			t.Errorf("test does not contain value=%v", got)
		}
	}

}

func TestSystemEntities_AddAndRemmove(t *testing.T) {
	var world = NewWorld(10, 8)
	var counter = Counter{}

	// Register Components
	world.RegisterComponents(CountID)

	// register systems
	world.RegisterSystem(counter.Increment, CountID)

	// create entities
	e1 := world.Create(&Count{7})
	e2 := world.Create(&Count{61})

	world.Tick()

	if len(counter.entities) != 2 {
		t.Errorf("Expected=%v, Got=%v", 2, len(counter.entities))
	}

	e1.Remove(CountID)
	world.Tick()

	if len(counter.entities) != 1 {
		t.Errorf("Expected=%v, Got=%v", 1, len(counter.entities))
	}

	e2.Remove(CountID)
	world.Tick()

	if len(counter.entities) != 0 {
		t.Errorf("Expected=%v, Got=%v", 0, len(counter.entities))
	}

	e1.Add(&Count{0})
	world.Tick()

	if len(counter.entities) != 1 {
		t.Errorf("Expected=%v, Got=%v", 1, len(counter.entities))
	}
}

// fix bug where the entity's signature doesn't update after adding new component
func TestEntitySignature_AddComponent(t *testing.T) {
	var world = NewWorld(10, 5)

	world.RegisterComponents(DummyID, CountID)

	entity := world.Create()

	entity.Add(&Count{0})

	sig := entity.Signature()
	ssig := sig.String()

	if !strings.HasSuffix(ssig, "0010") {
		t.Errorf("Expected=%v, Got=%v", "0010", ssig[len(ssig)-4:])
	}

	entity.Add(&Dummy{})

	sig = entity.Signature()
	ssig = sig.String()

	if !strings.HasSuffix(ssig, "0011") {
		t.Errorf("Expected=%v, Got=%v", "0011", ssig[len(ssig)-4:])
	}
}

func TestEntitySignature_RemoveComponent(t *testing.T) {
	var world = NewWorld(10, 5)

	world.RegisterComponents(DummyID, CountID)

	entity := world.Create(&Count{0}, &Dummy{})

	sig := entity.Signature()
	ssig := sig.String()

	if !strings.HasSuffix(ssig, "0011") {
		t.Errorf("Expected=%v, Got=%v", "0011", ssig[len(ssig)-4:])
	}

	entity.Remove(DummyID)

	sig = entity.Signature()
	ssig = sig.String()

	if !strings.HasSuffix(ssig, "0010") {
		t.Errorf("Expected=%v, Got=%v", "0010", ssig[len(ssig)-4:])
	}

	entity.Remove(CountID)

	sig = entity.Signature()
	ssig = sig.String()

	if !strings.HasSuffix(ssig, "0000") {
		t.Errorf("Expected=%v, Got=%v", "0000", ssig[len(ssig)-4:])
	}
}

func BenchmarkCreateEntity(b *testing.B) {
	type Position Vec2
	type Velocity Vec2

	var (
		pid = CreateComponentID[Position]()
		vid = CreateComponentID[Velocity]()
	)

	world := NewWorld(b.N, 2)
	world.RegisterComponents(pid, vid)

	for idx := 0; idx < b.N; idx++ {
		world.Create(&Position{}, &Velocity{})
	}
}

func BenchmarkRemoveEntity(b *testing.B) {
	type Position Vec2
	type Velocity Vec2

	var (
		pid = CreateComponentID[Position]()
		vid = CreateComponentID[Velocity]()
	)

	world := NewWorld(b.N, 2)
	world.RegisterComponents(pid, vid)

	for idx := 0; idx < b.N; idx++ {
		world.Create(&Position{}, &Velocity{})
	}

	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		world.RemoveEntity(idx)
	}
}

func BenchmarkRemove10Entity(b *testing.B) {
	type Position Vec2
	type Velocity Vec2

	var (
		pid = CreateComponentID[Position]()
		vid = CreateComponentID[Velocity]()
	)

	world := NewWorld(b.N, 2)
	world.RegisterComponents(pid, vid)

	for idx := 0; idx < b.N; idx++ {
		world.Create(&Position{}, &Velocity{})
	}

	b.ResetTimer()

	var count = 50

	for idx := 0; idx < (b.N / count); idx++ {

		for n := 0; n < count; n++ {
			world.RemoveEntity(idx + n)
		}

	}
}

type Position Vec2
type Velocity Vec2

func MovementSystem(entities []EntityHandle) {
	for _, entity := range entities {
		var pos *Position
		var vel *Velocity

		entity.Unpack(&pos, &vel)

		pos.X += vel.X
		pos.Y += vel.Y
	}
}

func BenchmarkItrEntity(b *testing.B) {
	var (
		pid = CreateComponentID[Position]()
		vid = CreateComponentID[Velocity]()
		K   = 9000
	)

	world := NewWorld(K, 2)
	world.RegisterComponents(pid, vid)
	world.RegisterSystem(MovementSystem, pid, vid)

	for idx := 0; idx < K; idx++ {
		world.Create(&Position{}, &Velocity{})
	}

	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		world.Tick()
	}
}
