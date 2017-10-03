package tabconv

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

var (
	escaper = strings.NewReplacer("\\", "_", "/", "_", " ", "_")
	tilde   = byte('~')
)

// basename returns filename without extension
func basename(name string) string {
	name = path.Base(name)
	ext := filepath.Ext(name)
	return name[:len(name)-len(ext)]
}

// fix replaces any path separators with underscores
func fix(filename string) string {
	return escaper.Replace(filename)
}

// expand expands '~' into current user directory
func Expand(path string) (string, error) {
	if path[0] != tilde {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return strings.Replace(path, "~", usr.HomeDir, 1), nil
}

// fileExists checks if file or directory exists
func FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
