package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenTerminalInDirectory opens a new terminal window in the specified directory.
// It's platform-aware and tries multiple terminal options on each platform.
func OpenTerminalInDirectory(dir string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Try Windows Terminal first, then PowerShell, then CMD
		cmd = exec.Command("wt.exe", "-w", "0", "nt", "-d", dir, "pwsh.exe")
		if err := cmd.Start(); err == nil {
			return nil
		}

		cmd = exec.Command("pwsh.exe", "-NoExit", "-Command", fmt.Sprintf("Set-Location '%s'", dir))
		if err := cmd.Start(); err == nil {
			return nil
		}

		cmd = exec.Command("cmd.exe", "/k", fmt.Sprintf("cd /d \"%s\"", dir))
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to open terminal: %w", err)
		}

	case "darwin":
		// macOS - use Terminal.app with AppleScript
		script := fmt.Sprintf(`tell application "Terminal" to do script "cd '%s'"`, dir)
		cmd = exec.Command("osascript", "-e", script)
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to open terminal: %w", err)
		}

	case "linux":
		// Linux - try common terminal emulators
		terminals := [][]string{
			{"gnome-terminal", "--working-directory=" + dir},
			{"konsole", "--workdir", dir},
			{"xfce4-terminal", "--working-directory=" + dir},
			{"xterm", "-e", fmt.Sprintf("cd '%s' && bash", dir)},
		}

		var lastErr error
		for _, termCmd := range terminals {
			cmd = exec.Command(termCmd[0], termCmd[1:]...)
			if err := cmd.Start(); err == nil {
				return nil
			} else {
				lastErr = err
			}
		}

		if lastErr != nil {
			return fmt.Errorf("failed to open terminal (tried multiple): %w", lastErr)
		}

	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return nil
}
