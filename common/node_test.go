package common

import (
	"testing"

	types "github.com/miguelmota/go-gun/types"
)

func TestNode(t *testing.T) {
	node := NewNode("abc")

	soul := node["_"].(types.Kv)["#"].(string)
	if soul != "abc" {
		t.Error("expected soul to equal abc")
	}

	kv1 := types.Kv{
		"foo": "goo",
	}
	kv2 := types.Kv{
		"bar": "baz",
	}

	node2 := NewNode("abc", kv1, kv2)

	if node2["foo"].(string) != "goo" {
		t.Error("expected kv1 to equal goo")
	}
	if node2["bar"].(string) != "baz" {
		t.Error("expected kv2 to equal baz")
	}
}
