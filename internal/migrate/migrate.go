package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gauravism2017/cctl/internal/config"
)

type Result struct {
	Skills int
	Rules  int
	Agents int
}

func Run(paths config.Paths) (Result, error) {
	if err := paths.EnsureDirs(); err != nil {
		return Result{}, fmt.Errorf("create store directories: %w", err)
	}

	var result Result
	var err error

	result.Skills, err = migrateType(paths.Skills, paths.VaultPath("skills"), paths.StoreSkills)
	if err != nil {
		return result, fmt.Errorf("migrate skills: %w", err)
	}

	result.Rules, err = migrateType(paths.Rules, paths.VaultPath("rules"), paths.StoreRules)
	if err != nil {
		return result, fmt.Errorf("migrate rules: %w", err)
	}

	result.Agents, err = migrateType(paths.Agents, paths.VaultPath("agents"), paths.StoreAgents)
	if err != nil {
		return result, fmt.Errorf("migrate agents: %w", err)
	}

	return result, nil
}

func migrateType(activeDir, vaultDir, storeDir string) (int, error) {
	count := 0

	migrated, err := migrateDir(activeDir, storeDir, true)
	if err != nil {
		return count, err
	}
	count += migrated

	migrated, err = migrateDir(vaultDir, storeDir, false)
	if err != nil {
		return count, err
	}
	count += migrated

	removeEmptyDir(vaultDir)

	return count, nil
}

func migrateDir(srcDir, storeDir string, createSymlinks bool) (int, error) {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}

	count := 0
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		srcPath := filepath.Join(srcDir, name)

		info, err := os.Lstat(srcPath)
		if err != nil {
			continue
		}
		if info.Mode()&os.ModeSymlink != 0 {
			continue
		}

		storePath := filepath.Join(storeDir, name)

		if _, err := os.Stat(storePath); err == nil {
			continue
		}

		if err := os.Rename(srcPath, storePath); err != nil {
			return count, fmt.Errorf("move %s to store: %w", name, err)
		}

		if createSymlinks {
			if err := os.Symlink(storePath, srcPath); err != nil {
				return count, fmt.Errorf("create symlink for %s: %w", name, err)
			}
		}

		count++
	}

	return count, nil
}

func removeEmptyDir(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	hasReal := false
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".") {
			hasReal = true
			break
		}
	}

	if !hasReal {
		_ = os.RemoveAll(dir)
	}
}
