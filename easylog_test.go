package easylog

import (
	"os"
	"testing"
)

func TestBasic(t *testing.T) {
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	Debug("Hello, World!")
}
