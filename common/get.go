package common

import (
	types "github.com/miguelmota/go-gun/types"
)

// Get ...
func Get(lex, graph types.Kv) types.Kv {
	if lex == nil {
		lex = make(types.Kv)
	}
	if graph == nil {
		graph = make(types.Kv)
	}
	soul := lex["#"].(string)
	node, ok := graph[soul]
	if !ok {
		return nil
	}
	key, ok := lex["."]
	if ok {
		tmp, ok := node.(types.Kv)[key.(string)]
		if !ok {
			return nil
		}
		// equiv: (node = {_: node._})[key] = tmp
		node = types.Kv{
			"_": node.(types.Kv)["_"],
		}
		node.(types.Kv)[key.(string)] = tmp

		// equiv: tmp = node._['>']
		tmp = node.(types.Kv)["_"].(types.Kv)[">"]

		// equiv: (node._['>'] = {})[key] = tmp[key]
		node.(types.Kv)["_"].(types.Kv)[">"] = make(types.Kv)

		node.(types.Kv)["_"].(types.Kv)[">"].(types.Kv)[key.(string)] = tmp.(types.Kv)[key.(string)]
	}

	ack := make(types.Kv)
	ack[soul] = node
	return ack
}
