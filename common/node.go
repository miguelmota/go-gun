package common

// NewNode ...
func NewNode(name string, args ...map[string]interface{}) map[string]interface{} {
	node := make(map[string]interface{})

	node["_"] = make(map[string]interface{})
	node["_"].(map[string]interface{})["#"] = name
	node["_"].(map[string]interface{})[">"] = make(map[string]interface{})

	for _, arg := range args {
		for k, v := range arg {
			node["_"].(map[string]interface{})[">"].(map[string]interface{})[k] = 0
			node[k] = v
		}
	}

	return node
}

// GetState ...
func GetState(node map[string]interface{}) map[string]interface{} {
	if _, ok := node["_"]; ok {
		return node["_"].(map[string]interface{})[">"].(map[string]interface{})
	}

	return make(map[string]interface{})
}

/*
def new_node(name, **kwargs):
    # node with meta
    node = {'_': {'#':name, '>':{k:0 for k in kwargs}}, **kwargs}
    print("NODE IS :" , node)
    return node
*/
