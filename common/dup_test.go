package common

import "testing"

func TestDup(t *testing.T) {
	d := NewDup()

	id := NewUID()
	exist := d.Check(id)

	if exist {
		t.Fail()
	}

	d.Track(id)

	exist = d.Check(id)

	if !exist {
		t.Fail()
	}
}
