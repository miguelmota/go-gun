package storage

import (
	"fmt"
	"testing"
)

func TestDummyKV(t *testing.T) {
	db := NewDummyKV()

	db.Put("abc", "foo", "bar", make(map[string]interface{}))

	val := db.Get("abc", nil)
	fmt.Println("RET1", val)

	k := "foo"
	value := db.Get("abc", &k)
	fmt.Println("RET2", value)
}
