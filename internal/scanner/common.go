package scanner

import (
	"os"
	"path/filepath"
)

func isSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}

func adoptItem(activeDir, storeDir, name string) error {
	activePath := filepath.Join(activeDir, name)
	storePath := filepath.Join(storeDir, name)

	if _, err := os.Stat(storePath); err == nil {
		return nil
	}

	if err := os.Rename(activePath, storePath); err != nil {
		return err
	}

	return os.Symlink(storePath, activePath)
}

func adoptNewItems(activeDir, storeDir string) {
	entries, err := os.ReadDir(activeDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		name := entry.Name()
		if name == "." || name == ".." || name[0] == '.' {
			continue
		}

		fullPath := filepath.Join(activeDir, name)
		if !isSymlink(fullPath) {
			_ = adoptItem(activeDir, storeDir, name)
		}
	}
}
