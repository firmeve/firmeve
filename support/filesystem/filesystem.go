package filesystem

import "os"

func IsFile(path string) bool {
	file, err := os.Stat(path)
	if err != nil || file.IsDir() {
		return false
	}

	return true
}
