package client

import (
	"fmt"
	"testing"

	types "github.com/miguelmota/go-gun/types"
)

func TestClient(t *testing.T) {
	t.Skip("requires server")

	client := NewClient("ws://localhost:8080")

	defer client.Close()

	var v = make(types.Kv)
	v["width"] = "10"
	v["height"] = "20"

	client.Put("box", v)

	val := client.Get("box", nil)
	fmt.Println("get value", val)
	if val["width"] != "10" {
		t.Errorf("expected 10, got %v", val["width"])
	}
	if val["height"] != "20" {
		t.Errorf("expected 20, got %v", val["height"])
	}

	k := "width"
	val2 := client.Get("box", &k)
	if val2[k] != "10" {
		t.Errorf("expected 10, got %v", val2[k])
	}
}
