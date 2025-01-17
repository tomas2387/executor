package commands

import (
	"os"
	"path/filepath"
	"strings"

	"executor/internal/config"
	"executor/internal/terminal"
)

func isExecutable(cmd string) bool {
	if cmd == "" {
		return false
	}

	if strings.Contains(cmd, "/") {
		info, err := os.Stat(cmd)
		if err != nil {
			return false
		}
		if !info.IsDir() && info.Mode().IsRegular() && info.Mode().Perm()&0o100 != 0 {
			return true
		}
		return false
	}

	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	for _, dir := range paths {
		fullPath := filepath.Join(dir, cmd)
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}
		if !info.IsDir() && info.Mode().IsRegular() && info.Mode().Perm()&0o100 != 0 {
			return true
		}
	}

	return false
}

func Which(cfg *config.Config) error {
	ok := isExecutable(cfg.Command)

	if ok && cfg.Silent {
		return nil
	}

	terminal.SetNoColor(cfg.NoColor)
	desc := "Looking for " + cfg.Command
	terminal.ActionNoColon(terminal.InfoLevel, desc)
	terminal.DashedLine(len(desc) + 2)
	terminal.Result(ok)

	if !ok {
		terminal.Line(terminal.WarnLevel, cfg.NotFoundMsg)
		return ErrCommandNotFound
	}

	return nil
}
