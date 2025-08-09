package cmd

import (
	"os/exec"
	"runtime"
)

func OpenDir(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default: // linux or others
		cmd = exec.Command("xdg-open", path)
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
}
