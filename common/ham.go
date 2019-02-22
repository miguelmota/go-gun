package common

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	log "github.com/sirupsen/logrus"

	types "github.com/miguelmota/go-gun/types"
)

// HamResponse ...
type HamResponse struct {
	Defer      bool
	Historical bool
	Converge   bool
	Incoming   bool
	Current    bool
	State      bool
}

// Ham is conflict resolution algorithm
func Ham(
	machineState float64,
	incomingState float64,
	currentState float64,
	incomingValue interface{},
	currentValue interface{},
) (*HamResponse, error) {
	if machineState < incomingState {
		// the incoming value is outside the boundary of the machine's state, it must be reprocessed in another state
		return &HamResponse{Defer: true}, nil
	}
	if incomingState < currentState {
		// the incoming value is within the boundary of the machine's state, but not within the range
		return &HamResponse{Historical: true}, nil
	}
	if currentState < incomingState {
		// the incoming value is within both the boundary and the range of the machine's state
		return &HamResponse{Converge: true, Incoming: true}, nil
	}
	var incomingVal string
	var currentVal string
	if incomingState == currentState {
		incomingVal = lexical(incomingValue)
		currentVal = lexical(currentValue)
		// NOTE: while these are practically the same, the deltas could technically be different
		if incomingVal == currentVal {
			return &HamResponse{State: true}, nil
		}
		// string only works on primitive values
		if incomingVal < currentVal {
			return &HamResponse{Converge: true, Current: true}, nil
		}
		// string only works on primitive values
		if currentVal < incomingVal {
			return &HamResponse{Converge: true, Incoming: true}, nil
		}
	}

	return nil, fmt.Errorf("Invalid CRDT Data: %s to %s at %v to %v", incomingVal, currentVal, incomingState, currentState)
}

// Mix applies updates 'change' to the graph
func Mix(change types.Kv, graph types.Kv) types.Kv {
	// NOTE: the timestamp are coming from js client that uses nanoseconds
	machine := float64(time.Now().Unix() * 1e3)
	diff := make(types.Kv)

	for soul, inode := range change {
		node := iToKv(inode)
		for key, val := range node {
			if key == "_" {
				continue
			}

			if node["_"] == nil {
				log.Fatal("should not be nil")
			}

			// equiv: state = node._['>'][key]
			state := getStateOfProp(node, key)
			oldNode := iToKv(graph[soul])
			was := getStateOfProp(oldNode, key)
			if was == 0 {
				// NOTE: "infinity"
				was = float64(math.MinInt64)
			}

			// equiv: known = (graph[soul]||{})[key]
			var known interface{} = get(oldNode, key)
			hm, err := Ham(
				machine,
				state,
				was,
				val,
				known,
			)
			if err != nil {
				log.Fatal(err)
			}

			if !hm.Incoming {
				if hm.Defer {
					// TODO: need to implement this
				}

				continue
			}

			_, ok := diff[soul]
			if !ok {
				// equiv: graph[soul] = {_:{'#':soul, '>':{}}}
				diff[soul] = newNode(soul)
			}

			_, ok = graph[soul]
			if !ok {
				// equiv: graph[soul] = {_:{'#':soul, '>':{}}}
				graph[soul] = newNode(soul)
			}

			graph[soul].(types.Kv)[key] = val
			diff[soul].(types.Kv)[key] = val

			diff[soul].(types.Kv)["_"].(types.Kv)[">"].(types.Kv)[key] = state
			graph[soul].(types.Kv)["_"].(types.Kv)[">"].(types.Kv)[key] = state
		}
	}

	return diff
}

// lexical ...
func lexical(value interface{}) string {
	js, err := json.Marshal(value)
	if err != nil {
		log.Fatal(err)
	}

	return string(js)
}
