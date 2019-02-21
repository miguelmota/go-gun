package storage

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// IStorage ... {
type IStorage interface {
	Put(soul string, key string, value interface{}, state interface{})
	Get(soul string, key *string) interface{}
	List() map[string]interface{}
	SetList(list map[string]interface{})
}

// DummyKV ...
type DummyKV struct {
	db map[string]interface{}
}

// NewDummyKV ...
func NewDummyKV() *DummyKV {
	return &DummyKV{
		db: make(map[string]interface{}),
	}
}

// Put ...
func (d *DummyKV) Put(soul string, key string, value interface{}, state interface{}) {
	k := fmt.Sprintf("%s:%s:%s", soul, key, state)
	d.db[k] = value
	//self.db["{soul}:{key}:{state}".format(**locals())] = value
}

// Get ...
func (d *DummyKV) Get(soul string, key *string) interface{} {
	ret := make(map[string]interface{})
	ret["#"] = soul
	ret["_"] = make(map[string]interface{})
	ret["_"].(map[string]interface{})["#"] = soul
	ret["_"].(map[string]interface{})[">"] = make(map[string]interface{})

	fmt.Println("DB")
	spew.Dump(d.db)

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

	fmt.Println("keys", keys)

	for _, k := range keys {
		parts := strings.Split(k, ":")
		//sol := parts[0]
		key := parts[1]
		state := parts[2]

		fmt.Println("KEY", key)
		fmt.Println("STATE", state)

		ret["_"].(map[string]interface{})[">"].(map[string]interface{})[key] = state
		ret[key] = d.db[k]
	}

	_ = keys

	/*
	   print("\n\n{} {} {}\n\n".format(soul, key, type(key)))
	   ret = {'#': soul, '_':{'#':soul, '>':{}}}
	   if isinstance(key, str):
	       keys = [k for k in self.db.keys() if k and k.startswith(soul+":"+key)]
	   else:
	       keys = [k for k in self.db.keys() if k.startswith(soul+":")]

	   for k in keys:
	       sol, key, state = k.split(":")
	       ret['_']['>'][key] = state
	       ret[key] = self.db[k]

	   return ret
	*/

	return ret
}

// putSoul ...
func (d *DummyKV) putSoul(soul string, souldict map[string]interface{}) {
	for k, v := range souldict {
		if k == "#" || k == "_" || k == ">" {
			continue
		}

		s1, ok := souldict["_"]
		if ok {
			s2, ok := s1.(map[string]interface{})[">"]
			if ok {
				kstate := s2.(map[string]interface{})[k]
				d.Put(soul, k, v, kstate)
				return
			}
		}

		kstate := make(map[string]interface{})
		kstate[k] = 0
		d.Put(soul, k, v, kstate)
	}
	/*
	   for k,v in souldict.items():
	       if k in "#_>":
	           continue
	       kstate = souldict.get("_", {}).get(">", {k:0})[k]
	       self.put(soul, k, v, kstate)
	*/
}

// List ...
func (d *DummyKV) List() map[string]interface{} {
	return d.db
	//return self.db.items()
}

// SetList ...
func (d *DummyKV) SetList(db map[string]interface{}) {
	d.db = db
}

// getItem ...
func (d *DummyKV) getItem(soul string) interface{} {
	return d.Get(soul, nil)
	//return self.get(soul, None)
}

// setItem ...
func (d *DummyKV) setItem(soul string, souldict map[string]interface{}) {
	d.putSoul(soul, souldict)
	//self.putsoul(soul, souldict)
}

/*

   def get(self, soul, key=None):
       print("\n\n{} {} {}\n\n".format(soul, key, type(key)))
       ret = {'#': soul, '_':{'#':soul, '>':{}}}
       if isinstance(key, str):
           keys = [k for k in self.db.keys() if k and k.startswith(soul+":"+key)]
       else:
           keys = [k for k in self.db.keys() if k.startswith(soul+":")]

       for k in keys:
           sol, key, state = k.split(":")
           ret['_']['>'][key] = state
           ret[key] = self.db[k]

       return ret
*/
