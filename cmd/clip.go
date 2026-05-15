package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
)

// clipboardCmd describes how to write to the system clipboard on a given OS.
type clipboardCmd struct {
	name string   // executable to lookup in PATH
	args []string // arguments passed to the executable
}

// clipboardCmds returns the candidate commands for the current OS, in order
// of preference. The first one whose binary is found in PATH is used.
func clipboardCmds() []clipboardCmd {
	switch runtime.GOOS {
	case "darwin":
		return []clipboardCmd{{name: "pbcopy"}}
	case "linux":
		return []clipboardCmd{
			{name: "wl-copy"},
			{name: "xclip", args: []string{"-selection", "clipboard"}},
			{name: "xsel", args: []string{"--clipboard", "--input"}},
		}
	default:
		return nil
	}
}

// copyToClipboard writes s to the system clipboard. It returns the name of
// the underlying tool used (for logging) or an error if no supported tool
// is available.
func copyToClipboard(s string) (string, error) {
	cands := clipboardCmds()
	if len(cands) == 0 {
		return "", fmt.Errorf("clipboard not supported on %s", runtime.GOOS)
	}
	var tried []string
	for _, c := range cands {
		path, err := exec.LookPath(c.name)
		if err != nil {
			tried = append(tried, c.name)
			continue
		}
		cmd := exec.Command(path, c.args...)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return "", fmt.Errorf("%s stdin pipe: %w", c.name, err)
		}
		if err := cmd.Start(); err != nil {
			return "", fmt.Errorf("%s start: %w", c.name, err)
		}
		if _, err := stdin.Write([]byte(s)); err != nil {
			_ = stdin.Close()
			_ = cmd.Wait()
			return "", fmt.Errorf("%s write: %w", c.name, err)
		}
		if err := stdin.Close(); err != nil {
			_ = cmd.Wait()
			return "", fmt.Errorf("%s close stdin: %w", c.name, err)
		}
		if err := cmd.Wait(); err != nil {
			return "", fmt.Errorf("%s wait: %w", c.name, err)
		}
		return c.name, nil
	}
	return "", fmt.Errorf("no clipboard tool found in PATH (tried: %v)", tried)
}
