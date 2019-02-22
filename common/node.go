package common

import types "github.com/miguelmota/go-gun/types"

// NewNode returns a new node
func NewNode(name string, args ...types.Kv) types.Kv {
	return newNode(name, args...)
}

// GetStateOfProp ...
func GetStateOfProp(node types.Kv, key string) float64 {
	return getStateOfProp(node, key)
}

// IToKv ...
func IToKv(i interface{}) types.Kv {
	return iToKv(i)
}

func newNode(name string, args ...types.Kv) types.Kv {
	node := make(types.Kv)

	node["_"] = make(types.Kv)
	node["_"].(types.Kv)["#"] = name
	node["_"].(types.Kv)[">"] = make(types.Kv)

	for _, arg := range args {
		for k, v := range arg {
			node["_"].(types.Kv)[">"].(types.Kv)[k] = 0
			node[k] = v
		}
	}

	return node
}

func getState(node types.Kv) types.Kv {
	if meta, ok := node["_"]; ok {
		return meta.(types.Kv)[">"].(types.Kv)
	}

	return make(types.Kv)
}

func getStateOfProp(node types.Kv, key string) float64 {
	v, _ := iToKv(iToKv(node["_"])[">"])[key].(float64)
	return v
}

func get(node interface{}, key string) interface{} {
	return node.(types.Kv)[key]
}

func iToKv(i interface{}) types.Kv {
	switch v := i.(type) {
	case nil:
		return types.Kv{}
	case types.Kv:
		return v
	case map[string]interface{}:
		return types.Kv(i.(map[string]interface{}))
	default:
		return types.Kv{}
	}
}
