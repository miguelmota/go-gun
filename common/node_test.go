package common

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestNode(t *testing.T) {
	node := NewNode("abc")

	spew.Dump(node)

	kv := map[string]interface{}{
		"foo": "bar",
	}
	kv2 := map[string]interface{}{
		"qux": "baz",
	}

	node2 := NewNode("abc", kv, kv2)

	spew.Dump(node2)
}
