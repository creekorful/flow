package runner

import "testing"

func TestInterpolate(t *testing.T) {
	cmd := "apt install {flags} {name}"

	_, err := interpolateCmd(cmd, map[string]string{"name": "nginx"})
	if err == nil {
		t.FailNow()
	}

	res, err := interpolateCmd(cmd, map[string]string{"name": "nginx", "flags": "-y"})
	if err != nil {
		t.FailNow()
	}
	if res != "apt install -y nginx" {
		t.FailNow()
	}
}
