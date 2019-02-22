package common

import (
	"math/rand"
	"time"

	types "github.com/miguelmota/go-gun/types"
)

// Opt ...
type Opt struct {
	Max int64
	Age int64
}

// Dup ...
type Dup struct {
	value types.Kv
	opt   *Opt
}

// NewDup ...
func NewDup() *Dup {
	value := make(types.Kv)
	value["s"] = make(types.Kv)

	return &Dup{
		value: value,
		opt: &Opt{
			Max: 1000,
			Age: 1000 * 9,
		},
	}
}

// Check ...
func (d *Dup) Check(id string) bool {
	_, ok := d.value["s"].(types.Kv)[id]
	if ok {
		d.Track(id)
		return true
	}

	return false
}

// Track ...
func (d *Dup) Track(id string) string {
	d.value["s"].(types.Kv)[id] = time.Now()

	_, ok := d.value["to"]
	if !ok {
		d.value["to"] = time.AfterFunc(time.Duration(d.opt.Age), func() {
			for id, v := range d.value["s"].(types.Kv) {
				t := v.(time.Time)
				if d.opt.Age > time.Now().Unix()-t.Unix() {
					continue
				}

				delete(d.value["s"].(types.Kv), id)
			}

			delete(d.value, "to")
		})
	}

	return id
}

// NewUID generate a new unique identifier
func NewUID() string {
	return Random()
}

// Random ...
func Random() string {
	return randStringRunes(3)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
