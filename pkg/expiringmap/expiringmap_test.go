package expiringmap

import (
	"testing"
	"time"
)

func TestExpiration(t *testing.T) {
	age := 2 * time.Second
	emap := NewExpiringMap[int](Settings{Age: age})
	emap.Set("foo", 1)
	val, ok := emap.Get("foo")
	if !ok {
		t.Error("expected to get val within time frame")
		return
	}
	if val != 1 {
		t.Errorf("expected val to be 1, was %v", val)
	}

	t.Log("Sleeping...")
	time.Sleep(age + 200*time.Millisecond)
	val, ok = emap.Get("foo")
	if ok {
		t.Error("expected val to be gone after time frame")
	}
	if val != 0 {
		t.Errorf("expected val to be 0 (empty value), was %v", val)
	}
}

func TestCapacityDeletion(t *testing.T) {
	emap := NewExpiringMap[int](Settings{MaxLength: 3})
	emap.Set("foo", 1)
	val, ok := emap.Get("foo")
	if !ok || val != 1 {
		t.Error("foo not set to 1")
	}

	emap.Set("bar", 2)
	val, ok = emap.Get("bar")
	if !ok || val != 2 {
		t.Error("bar not set to 2")
	}

	emap.Set("baz", 3)
	val, ok = emap.Get("baz")
	if !ok || val != 3 {
		t.Error("baz not set to 3")
	}

	emap.Set("toe", 4)
	val, ok = emap.Get("toe")
	if !ok || val != 4 {
		t.Error("toe not set to 4")
	}

	val, ok = emap.Get("foo")
	if ok {
		t.Error("expected foo field to be deleted")
	}
	val, ok = emap.Get("bar")
	if !ok {
		t.Error("expected bar field to not be deleted")
	}
}
