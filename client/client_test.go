package client

import (
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	client := NewClient("ws://localhost:8080")

	var v = make(map[string]interface{})
	v["width"] = "10"
	v["height"] = "20"

	client.Put("box", v)

	val := client.Get("box", nil)
	fmt.Println("get value", val)

	/*
		k := "width"
		val2 := client.Get("box", &k)
		fmt.Println("get value2", val2)
	*/

	/*
		result2 := client.Get("box", "width")
		fmt.Println(result2)
	*/
}
