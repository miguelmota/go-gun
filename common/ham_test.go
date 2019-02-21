package common

import "testing"

func TestHam(t *testing.T) {
	var machineState, incomingState, currentState, incomingValue, currentValue float64
	h, err := Ham(machineState, incomingState, currentState, incomingValue, currentValue)

	if err != nil {
		t.Error(err)
	}

	_ = h
}
