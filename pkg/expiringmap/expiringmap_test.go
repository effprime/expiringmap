package expiringmap

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	toSleep := 2 * time.Second
	emap := NewExpiringMap[int](Settings{Age: toSleep})
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
	time.Sleep(toSleep + 1*time.Second)
	val, ok = emap.Get("foo")
	if ok {
		t.Error("expected val to be gone after time frame")
		return
	}
	if val != 0 {
		t.Errorf("expected val to be 0 (empty value), was %v", val)
	}
}
