package ops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func enableItem(storeDir, activeDir, name string) error {
	if err := validateName(name); err != nil {
		return err
	}

	storePath := filepath.Join(storeDir, name)
	activePath := filepath.Join(activeDir, name)

	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		return fmt.Errorf("not found in store: %s", storePath)
	}

	if _, err := os.Lstat(activePath); err == nil {
		return nil
	}

	return os.Symlink(storePath, activePath)
}

func disableItem(activeDir, name string) error {
	if err := validateName(name); err != nil {
		return err
	}

	activePath := filepath.Join(activeDir, name)

	info, err := os.Lstat(activePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("stat %s: %w", activePath, err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("not a symlink: %s", activePath)
	}

	return os.Remove(activePath)
}

func validateName(name string) error {
	if strings.Contains(name, string(filepath.Separator)) || strings.Contains(name, "..") {
		return fmt.Errorf("invalid name: %s", name)
	}
	return nil
}
