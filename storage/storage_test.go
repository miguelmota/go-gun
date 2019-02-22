package storage

import (
	"testing"

	types "github.com/miguelmota/go-gun/types"
)

func TestDummyKV(t *testing.T) {
	backing := make(types.Kv)
	db := NewDummyKV(backing)

	db.Put("abc", "foo", "bar", make(types.Kv))

	v1 := db.Get("abc", nil)
	if v1["foo"].(string) != "bar" {
		t.Error("expected foo to be bar")
	}

	k := "foo"
	v2 := db.Get("abc", &k)
	if v2["foo"].(string) != "bar" {
		t.Error("expected foo to be bar")
	}
}
