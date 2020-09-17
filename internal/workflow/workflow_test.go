package workflow

import (
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	r := strings.NewReader("id: hello-world\nname: The way\nauthor: Aloïs Micard\ndescription: desc")

	w, err := Read(r)
	if err != nil {
		t.FailNow()
	}

	if w.ID != "hello-world" {
		t.FailNow()
	}
	if w.Name != "The way" {
		t.FailNow()
	}
	if w.Author != "Aloïs Micard" {
		t.FailNow()
	}
	if w.Description != "desc" {
		t.FailNow()
	}
}