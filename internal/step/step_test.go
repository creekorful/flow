package step

import "testing"

func TestStep_Runnable(t *testing.T) {
	s := Step{}

	if s.Runnable() {
		t.FailNow()
	}

	s.Exec = "test"
	if !s.Runnable() {
		t.FailNow()
	}
}
