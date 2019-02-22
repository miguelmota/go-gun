package storage

import (
	"fmt"
	"strings"

	types "github.com/miguelmota/go-gun/types"
)

// IStorage ... {
type IStorage interface {
	Put(soul string, key string, value interface{}, state interface{})
	Get(soul string, key *string) types.Kv
	List() types.Kv
	SetList(list types.Kv)
}

// DummyKV ...
type DummyKV struct {
	db types.Kv
}

// NewDummyKV ...
func NewDummyKV(kv types.Kv) *DummyKV {
	return &DummyKV{
		db: kv,
	}
}

// Put ...
func (d *DummyKV) Put(soul string, key string, value interface{}, state interface{}) {
	k := fmt.Sprintf("%s:%s:%s", soul, key, state)
	d.db[k] = value
}

// Get ...
func (d *DummyKV) Get(soul string, key *string) types.Kv {
	ret := make(types.Kv)
	ret["#"] = soul
	ret["_"] = make(types.Kv)
	ret["_"].(types.Kv)["#"] = soul
	ret["_"].(types.Kv)[">"] = make(types.Kv)

	var keys []string
	if key != nil {
		for k := range d.db {
			if strings.HasPrefix(k, soul+":"+*key) {
				keys = append(keys, k)
			}
		}
	} else {
		for k := range d.db {
			if strings.HasPrefix(k, soul+":") {
				keys = append(keys, k)
			}
		}
	}

	for _, k := range keys {
		parts := strings.Split(k, ":")
		//sol := parts[0]
		key := parts[1]
		state := parts[2]

		ret["_"].(types.Kv)[">"].(types.Kv)[key] = state
		ret[key] = d.db[k]
	}

	return ret
}

// putSoul ...
func (d *DummyKV) putSoul(soul string, souldict types.Kv) {
	for k, v := range souldict {
		if k == "#" || k == "_" || k == ">" {
			continue
		}

		s1, ok := souldict["_"]
		if ok {
			s2, ok := s1.(types.Kv)[">"]
			if ok {
				kstate := s2.(types.Kv)[k]
				d.Put(soul, k, v, kstate)
				return
			}
		}

		kstate := make(types.Kv)
		kstate[k] = 0
		d.Put(soul, k, v, kstate)
	}
}

// List ...
func (d *DummyKV) List() types.Kv {
	return d.db
}

// SetList ...
func (d *DummyKV) SetList(db types.Kv) {
	d.db = db
}

// GetItem ...
func (d *DummyKV) GetItem(soul string) types.Kv {
	return d.Get(soul, nil)
}

// SetItem ...
func (d *DummyKV) SetItem(soul string, souldict types.Kv) {
	d.putSoul(soul, souldict)
}
