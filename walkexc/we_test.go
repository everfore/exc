package walkexc

import (
	"path/filepath"
	"testing"
	// "toukii"
)

func TestSample(t *testing.T) {
	Setting(nil, "ls", "-a")
	filepath.Walk("../", WalkExc)
	Setting(nil, "go", "version")
	filepath.Walk("../", WalkExc)
}
