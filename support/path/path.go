package path

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
)

func callerStepFile(skip int) string {
	_, file, _, ok := runtime.Caller(skip)
	if !ok {
		panic(`Can not get current file info`)
	}

	return file
}

// Get current running file path
func RunFile() string {
	// 0 is current file, so except
	// 1 is current file, so except
	return callerStepFile(2)
}

// Get current running directory path
func RunDir() string {
	return path.Dir(callerStepFile(2))
}

// Get current directory relative path
func RunRelative(rpath string) string {
	rpath, _ = filepath.Abs(filepath.Join(path.Dir(callerStepFile(2)), rpath))
	return rpath
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
