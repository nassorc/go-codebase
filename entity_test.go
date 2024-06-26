package gandalf

// import "testing"

// func TestEntityManagerNewEntity(t *testing.T) {
// 	size := 5
// 	entityManager := newEntityManager(size)

// 	// entityManager.newEntity()
// 	for idx := 0; idx < size; idx++ {
// 		_, ok := entityManager.newEntity()

// 		if !ok {
// 			t.Errorf("Expected=%v, Got=%v", true, ok)
// 		}
// 	}

// 	if got := entityManager.empty(); got != true {
// 		t.Errorf("Expected=%v, Got=%v", true, got)
// 	}
// }

// func TestEntityManagerSignature(t *testing.T) {
// 	entityManager := newEntityManager(2)

// 	entity1, _ := entityManager.newEntity()
// 	signature1 := NewSignature(4)
// 	signature1.Set(0)
// 	signature1.Set(3)
// 	entityManager.setSignature(entity1, signature1)

// 	entity2, _ := entityManager.newEntity()
// 	signature2 := NewSignature(4)
// 	signature2.Set(1)
// 	signature2.Set(2)
// 	entityManager.setSignature(entity2, signature2)

// 	actualSignature1, ok := entityManager.getSignature(entity1)

// 	if !ok {
// 		t.Errorf("Expected=%v, Got=%v", true, ok)
// 	}

// 	if actualSignature1 == nil {
// 		t.Errorf("Expected=%v, Got=%v", signature1.String(), actualSignature1)
// 	}

// 	if actualSignature1.Int() != signature1.Int() {
// 		t.Errorf("Expected=%v, Got=%v", signature1.String(), actualSignature1.String())
// 	}

// 	actualSignature2, ok := entityManager.getSignature(entity2)

// 	if !ok {
// 		t.Errorf("Expected=%v, Got=%v", true, ok)
// 	}

// 	if actualSignature2 == nil {
// 		t.Errorf("Expected=%v, Got=%v", signature2.String(), actualSignature2)
// 	}

// 	if actualSignature2.Int() != signature2.Int() {
// 		t.Errorf("Expected=%v, Got=%v", signature2.String(), actualSignature2.String())
// 	}
// }

// func TestEntityManagerRemoveEntity(t *testing.T) {
// 	entityManager := newEntityManager(2)

// 	// test if entity1 id recycles after removing it
// 	entity1, _ := entityManager.newEntity()
// 	signature1 := NewSignature(4)
// 	signature1.Set(0)
// 	signature1.Set(3)
// 	entityManager.setSignature(entity1, signature1)

// 	entity2, _ := entityManager.newEntity()
// 	signature2 := NewSignature(4)
// 	signature2.Set(1)
// 	signature2.Set(2)
// 	entityManager.setSignature(entity2, signature2)

// 	oldId := entity1.Id()
// 	entityManager.removeEntity(entity1)

// 	entity3, ok := entityManager.newEntity()

// 	if !ok {
// 		t.Errorf("Expected=%v, Got=%v", true, ok)
// 	}

// 	// check if entity1 id (destroyed) has been requeued
// 	if entity3.Id() != oldId {
// 		t.Errorf("Expected=%v, Got=%v", oldId, entity3.Id())
// 	}
// }
