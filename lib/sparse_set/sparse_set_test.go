package sparse_set

import "testing"

type Counter struct {
  val int
}

func TestSparseSet_InsertAndGet(t *testing.T) {
  store := NewSparseSet[int](5)
  store.Insert(0, 100)
  store.Insert(1, 200)

  store.Insert(2, 300)
  store.Insert(3, 400)

  if v, _ := store.Get(0); *v != 100 {
    t.Errorf("Expected %v, Got %v", 100, v)
  }
  if v, _ := store.Get(3); *v != 400 {
    t.Errorf("Expected %v, Got %v", 400, v)
  }
}

func TestSparseSet(t *testing.T) {
  store := NewSparseSet[Counter](3)

  if v := store.Has(0); v {
    t.Errorf("Expected %v, Got %v", false, true)
  }

  store.Insert(0, Counter{10})
  store.Insert(1, Counter{14})
  store.Insert(2, Counter{18})
  shouldPanic(t, func() { store.Insert(2, Counter{18}) })

  // Querying
  if c, _ := store.Get(2); c.val != 18 {
    t.Errorf("Expected %v, Got %v", 18, c.val)
  }

  if c, _ := store.Get(1); c.val != 14 {
    t.Errorf("Expected %v, Got %v", 14, c.val)
  }

  if c, _ := store.Get(0); c.val != 10 {
    t.Errorf("Expected %v, Got %v", 10, c.val)
  }

  // Removing
  store.Remove(0)
  if v := store.Has(0); v {
    t.Errorf("Expected %v, Got %v", false, true)
  }

  store.Remove(1)
  if v := store.Has(1); v {
    t.Errorf("Expected %v, Got %v", false, true)
  }

  store.Remove(2)
  if v := store.Has(2); v {
    t.Errorf("Expected %v, Got %v", false, true)
  }

  if v := store.Remove(2); v {
    t.Errorf("Expected %v, Got %v", false, true)
  }

  if store.Size() != 0 {
    t.Errorf("Expected %v, Got %v", 0, store.Size())
  }
}

func TestSparse_ShouldNotInvalidationWhileLooping(t *testing.T) {
  n := 3
  store := NewSparseSet[Counter](n)
  saw := []int{}

  store.Insert(0, Counter{10})
  store.Insert(1, Counter{14})
  store.Insert(2, Counter{18})

  for idx := store.Size()-1; idx >= 0; idx-- {
    if idx == 2 {
      store.Remove(2)
    }

    saw = append(saw, idx)
  }

  if len(saw) != 3 {
    t.Errorf("Expected %v, Got %v", 3, len(saw))
    t.Errorf("Saw %v", saw)
  }
}

func shouldPanic(t *testing.T, cb func()) {
  defer func() { recover() }()
  cb()

  t.Errorf("Should have panicked")
}
