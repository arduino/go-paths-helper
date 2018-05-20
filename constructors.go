package paths

import (
	"io/ioutil"
	"os"
	"runtime"
)

// NullPath return the path to the /dev/null equivalent for the current OS
func NullPath() *Path {
	if runtime.GOOS == "windows" {
		return New("nul")
	}
	return New("/dev/null")
}

// TempDir returns the default path to use for temporary files
func TempDir() *Path {
	return New(os.TempDir())
}

// MkTempDir creates a new temporary directory in the directory
// dir with a name beginning with prefix and returns the path of
// the new directory. If dir is the empty string, TempDir uses the
// default directory for temporary files
func MkTempDir(dir, prefix string) (*Path, error) {
	path, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		return nil, err
	}
	return New(path), nil
}
