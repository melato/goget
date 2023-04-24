package util

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	g := NewGet[int, string](func(d int) string {
		return fmt.Sprintf("%d", d)
	})
	if g.Get(3) != "3" {
		t.Fail()
	}
}
