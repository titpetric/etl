package loader

import (
	"io/fs"
	"os"
)

var testConfig = struct {
	Storage fs.FS
	Path    string
}{
	Storage: os.DirFS("testdata"),
	Path:    "config.yml",
}
