package common

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
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
		return &HamResponse{Defer: true}, nil
	}
	if incomingState < currentState {
		return &HamResponse{Historical: true}, nil
	}
	if currentState < incomingState {
		return &HamResponse{Converge: true, Incoming: true}, nil
	}
	var incomingVal string
	var currentVal string
	if incomingState == currentState {
		incomingVal = Lexical(incomingValue)
		currentVal = Lexical(currentValue)
		if incomingVal == currentVal {
			return &HamResponse{State: true}, nil
		}
		if incomingVal < currentVal {
			return &HamResponse{Converge: true, Current: true}, nil
		}
		if currentVal < incomingVal {
			return &HamResponse{Converge: true, Incoming: true}, nil
		}
	}

	return nil, fmt.Errorf("Invalid CRDT Data: %s to %s at %v to %v", incomingVal, currentVal, incomingState, currentState)
}

// Mix applies updates 'change' to the graph
func Mix(change map[string]interface{}, graph map[string]interface{}) map[string]interface{} {
	machine := time.Now()
	var diff map[string]interface{}

	for soul := range change {
		fmt.Println("CHANGE", change)
		fmt.Println("SOUL", soul)
		fmt.Println("BOTH", change[soul])
		node, ok := change[soul].(map[string]interface{})
		if !ok {
			continue
		}
		fmt.Println("NODE", node)
		for key := range node {
			val := node[key]
			if key == "_" {
				continue
			}

			if node["_"] == nil {
				continue
			}

			// equiv: state = node._['>'][key]
			state := node["_"].(map[string]interface{})[">"].(map[string]interface{})[key]

			// equiv: was = (graph[soul]||{_:{'>':{}}})._['>'][key] || -Infinity
			var was interface{}
			soulv, ok := graph[soul]

			if ok {
				was = soulv.(map[string]interface{})["_"].(map[string]interface{})[">"].(map[string]interface{})[key]
			}

			if was == nil {
				// 'infinity'
				was = float64(math.MinInt64)
			}

			// equiv: known = (graph[soul]||{})[key]
			var known interface{}
			graphsoul, ok := graph[soul]
			if ok {
				known = graphsoul.(map[string]interface{})[key]
			}

			fmt.Println("SATE", state)
			if state == nil {
				state = 0
			}

			var stateF float64
			switch v := state.(type) {
			case nil:
				v = 0
			case string:
				var err error
				stateF, err = strconv.ParseFloat(v, 64)
				if err != nil {
					stateF = 0
					fmt.Println(err)
				}
			case float64:
				stateF = v
			}

			hm, err := Ham(
				float64(machine.Unix()),
				stateF,
				was.(float64),
				val,
				known,
			)

			if err != nil {
				log.Fatal(err)
			}
			if !hm.Incoming {
				if hm.Defer {
					// TODO: need to implement this
					// fmt.Println("defer", key, val)
				}

				continue
			}

			if diff == nil {
				diff = make(map[string]interface{})
			}

			_, ok = diff[soul]
			if !ok {
				// equiv: graph[soul] = {_:{'#':soul, '>':{}}}
				diff[soul] = make(map[string]interface{})
				diff[soul].(map[string]interface{})["_"] = make(map[string]interface{})
				diff[soul].(map[string]interface{})["_"].(map[string]interface{})["#"] = soul
				diff[soul].(map[string]interface{})["_"].(map[string]interface{})[">"] = make(map[string]interface{})
			}

			_, ok = graph[soul]
			if !ok {
				// equiv: graph[soul] = {_:{'#':soul, '>':{}}}
				graph[soul] = make(map[string]interface{})
				graph[soul].(map[string]interface{})["_"] = make(map[string]interface{})
				graph[soul].(map[string]interface{})["_"].(map[string]interface{})["#"] = soul
				graph[soul].(map[string]interface{})["_"].(map[string]interface{})[">"] = make(map[string]interface{})
			}

			graph[soul].(map[string]interface{})[key] = val
			diff[soul].(map[string]interface{})[key] = val

			diff[soul].(map[string]interface{})["_"].(map[string]interface{})[">"].(map[string]interface{})[key] = state
			graph[soul].(map[string]interface{})["_"].(map[string]interface{})[">"].(map[string]interface{})[key] = state
		}
	}

	return diff
}

// Lexical ...
func Lexical(value interface{}) string {
	js, err := json.Marshal(value)
	if err != nil {
		log.Fatal(err)
	}

	return string(js)
}
