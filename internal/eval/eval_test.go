package eval

import "testing"

func TestEvaluate(t *testing.T) {
	if Evaluate("os.name == linux", map[string]string{}) {
		t.FailNow()
	}
	if Evaluate("os.name == linux", map[string]string{"os.name": "windows"}) {
		t.FailNow()
	}
	if !Evaluate("os.name == linux", map[string]string{"os.name": "linux"}) {
		t.FailNow()
	}
	if Evaluate("os.name != linux", map[string]string{"os.name": "linux"}) {
		t.FailNow()
	}
	if !Evaluate("os.name != linux", map[string]string{"os.name": "windows"}) {
		t.FailNow()
	}
}
