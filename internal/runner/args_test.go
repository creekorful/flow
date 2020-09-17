package runner

import (
	"runtime"
	"testing"
)

func TestNewArgs(t *testing.T) {
	args := NewArgs()
	values := args.Values()

	if values["os.name"] != runtime.GOOS {
		t.FailNow()
	}
	if values["os.arch"] != runtime.GOARCH {
		t.FailNow()
	}
}

func TestArguments_UpdateIN(t *testing.T) {
	args := NewArgs()
	args = args.Update("in", "result", "test")
	values := args.Values()

	if values["in.result"] != "test" {
		t.FailNow()
	}
}
