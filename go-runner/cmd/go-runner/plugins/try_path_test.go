package plugins

import (
	"testing"
)

func TestJoin(t *testing.T) {
	res, err := Join("a", "/b")
	if err != nil {
		t.Error(err)
	}
	println(res)
}
