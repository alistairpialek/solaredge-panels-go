package solaredge

import "testing"

func TestAdd(t *testing.T) {
	got := 11
	want := 11

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}
